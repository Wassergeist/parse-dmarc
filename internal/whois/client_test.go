package whois

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// errResolver is a stand-in for any reverse-DNS failure.
var errResolver = errors.New("resolver unavailable")

// fakeResolver stands in for reverse DNS so no test touches the network.
type fakeResolver struct {
	names map[string][]string
	err   error
	calls atomic.Int32
}

func (f *fakeResolver) LookupAddr(_ context.Context, addr string) ([]string, error) {
	f.calls.Add(1)
	if f.err != nil {
		return nil, f.err
	}
	names, ok := f.names[addr]
	if !ok {
		return nil, &net.DNSError{Err: "no such host", Name: addr, IsNotFound: true}
	}
	return names, nil
}

// newTestClient wires a Client to a local server and a fake resolver.
func newTestClient(t *testing.T, handler http.HandlerFunc, resolver Resolver, opts ...Option) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	base := []Option{
		WithBaseURL(srv.URL + "/ip/"),
		WithHTTPClient(srv.Client()),
		WithResolver(resolver),
	}
	return New(append(base, opts...)...)
}

func TestLookupSuccess(t *testing.T) {
	fixture := loadFixture(t, "ripe_hetzner.json")
	resolver := &fakeResolver{names: map[string][]string{"5.9.1.1": {"static.1.1.9.5.clients.your-server.de."}}}

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(fixture)
	}, resolver)

	got := client.Lookup(context.Background(), "5.9.1.1")

	if got.Source != SourceRDAP {
		t.Errorf("Source = %q, want %q (err: %s)", got.Source, SourceRDAP, got.Err)
	}
	if got.Org != "Hetzner Online GmbH" {
		t.Errorf("Org = %q", got.Org)
	}
	// The trailing dot of a PTR record must not reach the dashboard.
	if got.Hostname != "static.1.1.9.5.clients.your-server.de" {
		t.Errorf("Hostname = %q", got.Hostname)
	}
}

func TestLookupFollowsRedirect(t *testing.T) {
	// rdap.org answers with a redirect to the responsible RIR.
	fixture := loadFixture(t, "arin_google.json")
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/ip/") {
			http.Redirect(w, r, "/rir/net", http.StatusFound)
			return
		}
		_, _ = w.Write(fixture)
	}, &fakeResolver{})

	got := client.Lookup(context.Background(), "209.85.128.1")

	if got.Org != "Google LLC" {
		t.Errorf("Org = %q, want Google LLC (err: %s)", got.Org, got.Err)
	}
}

func TestLookupSkipsReservedIPs(t *testing.T) {
	var hits atomic.Int32
	resolver := &fakeResolver{}
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		_, _ = w.Write([]byte("{}"))
	}, resolver)

	reserved := []string{
		"10.1.2.3",        // RFC 1918
		"192.168.1.1",     // RFC 1918
		"127.0.0.1",       // loopback
		"::1",             // IPv6 loopback
		"fe80::1",         // link-local
		"100.64.0.1",      // CGNAT
		"192.0.2.1",       // TEST-NET-1, used by this repo's fixtures
		"::ffff:10.0.0.1", // IPv4-mapped private address
		"224.0.0.1",       // multicast
		"0.0.0.0",         // unspecified
	}

	for _, ip := range reserved {
		got := client.Lookup(context.Background(), ip)
		if got.Source != SourcePrivate {
			t.Errorf("Lookup(%s).Source = %q, want %q", ip, got.Source, SourcePrivate)
		}
	}

	if hits.Load() != 0 {
		t.Errorf("reserved IPs caused %d registry requests, want 0", hits.Load())
	}
	if resolver.calls.Load() != 0 {
		t.Errorf("reserved IPs caused %d DNS lookups, want 0", resolver.calls.Load())
	}
}

func TestLookupInvalidIP(t *testing.T) {
	var hits atomic.Int32
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
	}, &fakeResolver{})

	// An RFC 7489 null report can leave source_ip empty.
	for _, ip := range []string{"", "not-an-ip"} {
		got := client.Lookup(context.Background(), ip)
		if got.Source != SourceError {
			t.Errorf("Lookup(%q).Source = %q, want %q", ip, got.Source, SourceError)
		}
	}
	if hits.Load() != 0 {
		t.Errorf("invalid IPs caused %d requests, want 0", hits.Load())
	}
}

func TestLookupRegistryFailures(t *testing.T) {
	tests := map[string]struct {
		status     int
		hostname   string
		wantSource string
	}{
		"404 with a PTR record still yields the hostname": {
			status: http.StatusNotFound, hostname: "mail.example.com", wantSource: SourceRDNS,
		},
		"404 with nothing else is an error": {
			status: http.StatusNotFound, wantSource: SourceError,
		},
		"429 falls back to reverse DNS": {
			status: http.StatusTooManyRequests, hostname: "mail.example.com", wantSource: SourceRDNS,
		},
		"500 with nothing else is an error": {
			status: http.StatusInternalServerError, wantSource: SourceError,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			resolver := &fakeResolver{names: map[string][]string{}}
			if tt.hostname != "" {
				resolver.names["203.0.114.5"] = []string{tt.hostname}
			}
			client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
			}, resolver)

			got := client.Lookup(context.Background(), "203.0.114.5")

			if got.Source != tt.wantSource {
				t.Errorf("Source = %q, want %q", got.Source, tt.wantSource)
			}
			if got.Err == "" {
				t.Error("Err is empty, want the failure recorded for caching")
			}
			if got.Hostname != tt.hostname {
				t.Errorf("Hostname = %q, want %q", got.Hostname, tt.hostname)
			}
		})
	}
}

func TestLookupRespectsConcurrencyLimit(t *testing.T) {
	var inFlight, maxInFlight atomic.Int32
	fixture := loadFixture(t, "lacnic_minimal.json")

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		current := inFlight.Add(1)
		for {
			observed := maxInFlight.Load()
			if current <= observed || maxInFlight.CompareAndSwap(observed, current) {
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
		inFlight.Add(-1)
		_, _ = w.Write(fixture)
	}, &fakeResolver{}, WithConcurrency(2))

	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client.Lookup(context.Background(), "200.3.12."+string(rune('0'+i)))
		}()
	}
	wg.Wait()

	if got := maxInFlight.Load(); got > 2 {
		t.Errorf("max concurrent requests = %d, want at most 2", got)
	}
}

func TestLookupTimesOut(t *testing.T) {
	release := make(chan struct{})
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		<-release
	}, &fakeResolver{}, WithHTTPClient(&http.Client{Timeout: 50 * time.Millisecond}))
	// Registered after the server so it runs before the server's own cleanup:
	// httptest.Server.Close waits for the blocked handler to return.
	t.Cleanup(func() { close(release) })

	done := make(chan Info, 1)
	go func() { done <- client.Lookup(context.Background(), "203.0.114.5") }()

	select {
	case got := <-done:
		if got.Err == "" {
			t.Error("expected a timeout to be recorded in Err")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Lookup did not return after the HTTP client timeout")
	}
}

func TestLookupCapsResponseSize(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		// A registry we do not control must not be able to exhaust memory.
		chunk := strings.Repeat("a", 64*1024)
		for range 40 {
			_, _ = w.Write([]byte(chunk))
		}
	}, &fakeResolver{})

	got := client.Lookup(context.Background(), "203.0.114.5")

	if got.Err == "" {
		t.Error("expected an oversized body to be rejected")
	}
}

func TestLookupSurvivesResolverFailure(t *testing.T) {
	fixture := loadFixture(t, "ripe_hetzner.json")
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(fixture)
	}, &fakeResolver{err: errResolver})

	got := client.Lookup(context.Background(), "5.9.1.1")

	if got.Source != SourceRDAP {
		t.Errorf("Source = %q, want %q", got.Source, SourceRDAP)
	}
	if got.Hostname != "" {
		t.Errorf("Hostname = %q, want empty", got.Hostname)
	}
	// A missing PTR record is normal and must not be reported as a failure.
	if got.Err != "" {
		t.Errorf("Err = %q, want empty", got.Err)
	}
}
