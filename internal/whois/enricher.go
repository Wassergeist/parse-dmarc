package whois

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/meysam81/parse-dmarc/internal/storage"
)

const (
	// errorTTL is how long a failed lookup is remembered. It keeps a dashboard
	// that polls every few minutes from hammering an unreachable registry,
	// which matters most for installations with no outbound access at all.
	errorTTL = time.Hour
	// privateTTL applies to reserved ranges: those never change hands.
	privateTTL = 30 * 24 * time.Hour
	// lookupPause spaces out requests so we stay a polite RDAP client.
	lookupPause = 250 * time.Millisecond
	// queueSize bounds the backlog; beyond it new IPs are dropped and picked
	// up on the next trigger rather than blocking a caller.
	queueSize = 256
	// busyRetries covers the short window where the fetch loop holds the
	// SQLite write lock.
	busyRetries = 3
)

// Enricher fills the ip_whois cache in the background. Callers only ever
// enqueue; nothing on an HTTP path waits for a lookup.
type Enricher struct {
	client *Client
	store  *storage.Storage
	log    *zerolog.Logger
	ttl    time.Duration

	queue chan string

	mu      sync.Mutex
	pending map[string]struct{}

	// rateLimitedAt throttles the "we are being rate limited" warning.
	rateLimitedAt time.Time
}

// NewEnricher builds an Enricher writing successful lookups with the given TTL.
func NewEnricher(client *Client, store *storage.Storage, log *zerolog.Logger, ttl time.Duration) *Enricher {
	return &Enricher{
		client:  client,
		store:   store,
		log:     log,
		ttl:     ttl,
		queue:   make(chan string, queueSize),
		pending: make(map[string]struct{}),
	}
}

// Start drains the queue until ctx is cancelled. Run it in its own goroutine.
func (e *Enricher) Start(ctx context.Context) {
	e.log.Info().Dur("ttl", e.ttl).Msg("whois enrichment started")
	for {
		select {
		case <-ctx.Done():
			e.log.Debug().Msg("whois enrichment stopped")
			return
		case ip := <-e.queue:
			info := e.lookupAndStore(ctx, ip)
			e.done(ip)
			// Reserved ranges cost no request, so they need no pause.
			if info.Source != SourcePrivate {
				select {
				case <-ctx.Done():
					return
				case <-time.After(lookupPause):
				}
			}
		}
	}
}

// EnqueueStale queues the sources whose cached entry is missing or expired.
func (e *Enricher) EnqueueStale(sources []storage.TopSourceIP) {
	now := time.Now().Unix()
	stale := make([]string, 0, len(sources))
	for _, s := range sources {
		if s.WhoisExpiresAt < now {
			stale = append(stale, s.SourceIP)
		}
	}
	e.Enqueue(stale...)
}

// Enqueue schedules IPs for lookup. It never blocks: an IP already queued or
// in flight is skipped, and a full queue drops the rest until next time.
func (e *Enricher) Enqueue(ips ...string) {
	for _, ip := range ips {
		if ip == "" || !e.claim(ip) {
			continue
		}
		select {
		case e.queue <- ip:
		default:
			e.done(ip)
			e.log.Debug().Str("ip", ip).Msg("whois queue full, skipping")
		}
	}
}

// claim reserves an IP, reporting false if it is already queued or in flight.
func (e *Enricher) claim(ip string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.pending[ip]; exists {
		return false
	}
	e.pending[ip] = struct{}{}
	return true
}

func (e *Enricher) done(ip string) {
	e.mu.Lock()
	delete(e.pending, ip)
	e.mu.Unlock()
}

// lookupAndStore performs one lookup and caches whatever came back.
func (e *Enricher) lookupAndStore(ctx context.Context, ip string) Info {
	start := time.Now()
	info := e.client.Lookup(ctx, ip)

	now := time.Now()
	ttl := e.ttl
	switch info.Source {
	case SourcePrivate:
		ttl = privateTTL
	case SourceError:
		ttl = errorTTL
	}

	entry := &storage.IPWhois{
		IP:         ip,
		Org:        info.Org,
		Network:    info.Network,
		CIDR:       info.CIDR,
		Country:    info.Country,
		Hostname:   info.Hostname,
		Source:     info.Source,
		LastError:  info.Err,
		LookedUpAt: now.Unix(),
		ExpiresAt:  now.Add(ttl).Unix(),
	}

	if err := e.upsertWithRetry(ctx, entry); err != nil {
		e.log.Warn().Err(err).Str("ip", ip).Msg("failed to cache whois lookup")
		return info
	}

	e.log.Debug().
		Str("ip", ip).
		Str("source", info.Source).
		Str("org", info.Org).
		Dur("duration", time.Since(start)).
		Msg("whois lookup complete")

	if strings.Contains(info.Err, "429") {
		e.warnRateLimited(ip)
	}
	return info
}

// upsertWithRetry works around the fetch loop briefly holding the write lock.
func (e *Enricher) upsertWithRetry(ctx context.Context, entry *storage.IPWhois) error {
	var err error
	delay := 100 * time.Millisecond
	for attempt := range busyRetries {
		if err = e.store.UpsertIPWhois(entry); err == nil {
			return nil
		}
		if !isBusyErr(err) || attempt == busyRetries-1 {
			return err
		}
		select {
		case <-ctx.Done():
			return err
		case <-time.After(delay):
		}
		delay *= 3
	}
	return err
}

func isBusyErr(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "database is locked") || strings.Contains(msg, "busy")
}

// warnRateLimited logs at most once a minute so a throttling registry does
// not fill the log with the same line.
func (e *Enricher) warnRateLimited(ip string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if time.Since(e.rateLimitedAt) < time.Minute {
		return
	}
	e.rateLimitedAt = time.Now()
	e.log.Warn().Str("ip", ip).Msg("rdap registry is rate limiting lookups; results will fill in gradually")
}
