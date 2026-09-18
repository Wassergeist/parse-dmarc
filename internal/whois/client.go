// Package whois resolves the owner of a sending IP address. It queries RDAP
// (RFC 7482), which returns structured JSON, rather than classic WHOIS, whose
// free-text output differs per registry; reverse DNS is looked up alongside it
// because a PTR record is often the name a reader recognises first.
package whois

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Outcomes recorded in Info.Source. They are mirrored in the ip_whois table.
const (
	SourceRDAP    = "rdap"
	SourceRDNS    = "rdns"
	SourcePrivate = "private"
	SourceError   = "error"
)

const (
	defaultRDAPTimeout = 10 * time.Second
	defaultRDNSTimeout = 3 * time.Second
	// maxBodySize caps what we read from a registry we do not control.
	maxBodySize = 1 << 20
	// defaultConcurrency stays low: RIR RDAP endpoints rate-limit bursts.
	defaultConcurrency = 2
)

// Info is the result of one lookup. An empty field means "not published",
// which is common and never an error on its own.
type Info struct {
	IP       string
	Org      string
	Network  string
	CIDR     string
	Country  string
	Hostname string
	Source   string
	// Err carries the last failure so it can be cached, keeping a broken or
	// unreachable registry from being retried on every dashboard load.
	Err string
}

// Resolver is the reverse-DNS surface we need; *net.Resolver satisfies it.
type Resolver interface {
	LookupAddr(ctx context.Context, addr string) ([]string, error)
}

// Client performs RDAP and reverse-DNS lookups.
type Client struct {
	http      *http.Client
	baseURL   string
	resolver  Resolver
	sem       chan struct{}
	userAgent string
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL sets the RDAP endpoint an IP is appended to.
func WithBaseURL(u string) Option {
	return func(c *Client) {
		if u != "" {
			if !strings.HasSuffix(u, "/") {
				u += "/"
			}
			c.baseURL = u
		}
	}
}

// WithHTTPClient replaces the HTTP client, mainly for tests.
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		if h != nil {
			c.http = h
		}
	}
}

// WithResolver replaces the reverse-DNS resolver, mainly for tests.
func WithResolver(r Resolver) Option {
	return func(c *Client) {
		if r != nil {
			c.resolver = r
		}
	}
}

// WithConcurrency caps how many lookups may be in flight at once.
func WithConcurrency(n int) Option {
	return func(c *Client) {
		if n > 0 {
			c.sem = make(chan struct{}, n)
		}
	}
}

// WithUserAgent sets the User-Agent; registries ask clients to identify.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		if ua != "" {
			c.userAgent = ua
		}
	}
}

// New builds a Client with sensible defaults for talking to public registries.
func New(opts ...Option) *Client {
	c := &Client{
		http:      &http.Client{Timeout: defaultRDAPTimeout},
		baseURL:   "https://rdap.org/ip/",
		resolver:  net.DefaultResolver,
		sem:       make(chan struct{}, defaultConcurrency),
		userAgent: "parse-dmarc (+https://github.com/dmarcguardhq/parse-dmarc)",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Lookup resolves one IP. It never returns an error: a failure is reported in
// the returned Info so the caller can cache it and move on.
func (c *Client) Lookup(ctx context.Context, ip string) Info {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		// Reports with no rows (RFC 7489 null reports) can carry an empty IP.
		return Info{IP: ip, Source: SourceError, Err: "invalid ip address"}
	}
	if isReservedIP(parsed) {
		return Info{IP: ip, Source: SourcePrivate}
	}

	select {
	case c.sem <- struct{}{}:
		defer func() { <-c.sem }()
	case <-ctx.Done():
		return Info{IP: ip, Source: SourceError, Err: ctx.Err().Error()}
	}

	canonical := parsed.String()
	info := Info{IP: ip}

	var wg sync.WaitGroup
	var hostname string
	var rdapInfo Info
	var rdapErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		hostname = c.reverseDNS(ctx, canonical)
	}()
	go func() {
		defer wg.Done()
		rdapInfo, rdapErr = c.rdap(ctx, canonical)
	}()
	wg.Wait()

	info.Hostname = hostname
	if rdapErr != nil {
		info.Err = rdapErr.Error()
		// A PTR record alone is still worth showing and worth caching.
		if hostname != "" {
			info.Source = SourceRDNS
		} else {
			info.Source = SourceError
		}
		return info
	}

	info.Org = rdapInfo.Org
	info.Network = rdapInfo.Network
	info.CIDR = rdapInfo.CIDR
	info.Country = rdapInfo.Country
	info.Source = SourceRDAP
	return info
}

// reverseDNS returns the first PTR name, or "" if there is none. A missing
// PTR record is the norm for plenty of senders, so failures stay silent.
func (c *Client) reverseDNS(ctx context.Context, ip string) string {
	ctx, cancel := context.WithTimeout(ctx, defaultRDNSTimeout)
	defer cancel()

	names, err := c.resolver.LookupAddr(ctx, ip)
	if err != nil || len(names) == 0 {
		return ""
	}
	return strings.TrimSuffix(names[0], ".")
}

// rdap fetches and parses the registry entry for ip.
func (c *Client) rdap(ctx context.Context, ip string) (Info, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRDAPTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+ip, nil)
	if err != nil {
		return Info{}, fmt.Errorf("build rdap request: %w", err)
	}
	req.Header.Set("Accept", "application/rdap+json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return Info{}, fmt.Errorf("rdap request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	switch {
	case resp.StatusCode == http.StatusNotFound:
		// Unallocated space. Nothing to retry for, so this is cached as a
		// normal result rather than a transient failure.
		return Info{}, errNotFound
	case resp.StatusCode != http.StatusOK:
		return Info{}, fmt.Errorf("rdap status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
	if err != nil {
		return Info{}, fmt.Errorf("read rdap body: %w", err)
	}

	return parseRDAP(body)
}
