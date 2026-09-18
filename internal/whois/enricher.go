package whois

import (
	"context"
	"net"
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
	// lookupPause spaces out registry requests so we stay a polite client. It
	// applies only when a registry was actually asked.
	lookupPause = 250 * time.Millisecond
	// workers decides how many lookups run at once. Registries rate-limit
	// bursts, so this stays small; the real speed-up comes from reusing an
	// answer across every address in the range it covered.
	workers = 3
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

	// ranges remembers which address range each answer covered.
	ranges rangeIndex

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
	e.seedRanges()
	e.log.Info().
		Dur("ttl", e.ttl).
		Int("workers", workers).
		Int("known_ranges", e.ranges.Len()).
		Msg("whois enrichment started")

	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.work(ctx)
		}()
	}
	wg.Wait()
	e.log.Debug().Msg("whois enrichment stopped")
}

// work drains the queue until the context is cancelled.
func (e *Enricher) work(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case ip := <-e.queue:
			queried := e.lookupAndStore(ctx, ip)
			e.done(ip)
			// Only an actual registry request earns a pause; a reserved range
			// or an answer taken from a range we already knew costs nothing.
			if queried {
				select {
				case <-ctx.Done():
					return
				case <-time.After(lookupPause):
				}
			}
		}
	}
}

// seedRanges loads the ranges of previous answers, so a restart does not
// re-query a registry for an address it can already place.
func (e *Enricher) seedRanges() {
	cached, err := e.store.GetFreshWhoisRanges(time.Now().Unix())
	if err != nil {
		e.log.Debug().Err(err).Msg("could not preload known whois ranges")
		return
	}
	for _, c := range cached {
		e.ranges.Add(Info{
			Org:     c.Org,
			Network: c.Network,
			CIDR:    c.CIDR,
			Country: c.Country,
			Source:  c.Source,
		})
	}
}

// EnqueueStale queues the sources whose cached entry is missing or expired.
func (e *Enricher) EnqueueStale(sources []storage.TopSourceIP) {
	now := time.Now().Unix()
	stale := make([]string, 0, len(sources))
	for i := range sources {
		if sources[i].WhoisStale(now) {
			stale = append(stale, sources[i].SourceIP)
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

// lookupAndStore resolves one address and caches the result. It reports
// whether a registry was actually queried.
func (e *Enricher) lookupAndStore(ctx context.Context, ip string) bool {
	start := time.Now()

	info, queried := e.resolve(ctx, ip)

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
		return queried
	}

	e.log.Debug().
		Str("ip", ip).
		Str("source", info.Source).
		Str("org", info.Org).
		Bool("queried_registry", queried).
		Dur("duration", time.Since(start)).
		Msg("whois lookup complete")

	if strings.Contains(info.Err, "429") {
		e.warnRateLimited(ip)
	}
	return queried
}

// resolve answers from a known range where possible, and asks a registry
// otherwise. A report typically lists many relays out of one assignment, so
// this turns a dozen registry requests into one plus a dozen PTR lookups.
func (e *Enricher) resolve(ctx context.Context, ip string) (Info, bool) {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed != nil && !isReservedIP(parsed) {
		if known, ok := e.ranges.Lookup(parsed); ok {
			known.IP = ip
			known.Hostname = e.client.Hostname(ctx, ip)
			return known, false
		}
	}

	info := e.client.Lookup(ctx, ip)
	if info.Source == SourceRDAP || info.Source == SourceWHOIS {
		e.ranges.Add(Info{
			Org:     info.Org,
			Network: info.Network,
			CIDR:    info.CIDR,
			Country: info.Country,
			Source:  info.Source,
		})
	}
	return info, info.Source != SourcePrivate
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
