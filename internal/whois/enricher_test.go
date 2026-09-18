package whois

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/meysam81/parse-dmarc/internal/storage"
)

// newTestEnricher builds an Enricher over a real on-disk database. A file is
// used rather than :memory: because the worker runs on its own goroutine and
// database/sql may hand it a second connection, which for :memory: would be a
// separate, empty database.
func newTestEnricher(t *testing.T, handler http.HandlerFunc) (*Enricher, *storage.Storage) {
	t.Helper()

	store, err := storage.NewStorage(filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatalf("create storage: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	client := New(
		WithBaseURL(srv.URL+"/ip/"),
		WithHTTPClient(srv.Client()),
		WithResolver(&fakeResolver{}),
		// Keep the enricher tests to the RDAP path and off the network.
		WithWhoisFallback(false),
		WithDialer(refuseDial),
	)
	log := zerolog.Nop()

	return NewEnricher(client, store, &log, 24*time.Hour), store
}

// waitForEntry polls until the worker has written the row, or fails.
func waitForEntry(t *testing.T, store *storage.Storage, ip string) storage.IPWhois {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		entries, err := store.GetIPWhois([]string{ip})
		if err != nil {
			t.Fatalf("get ip whois: %v", err)
		}
		if entry, ok := entries[ip]; ok {
			return entry
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("no cache entry for %s within the deadline", ip)
	return storage.IPWhois{}
}

func TestEnricherStoresLookupOnce(t *testing.T) {
	var hits atomic.Int32
	fixture := loadFixture(t, "ripe_hetzner.json")

	enricher, store := newTestEnricher(t, func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		_, _ = w.Write(fixture)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go enricher.Start(ctx)

	// The same IP queued repeatedly must still cost exactly one request.
	enricher.Enqueue("5.9.1.1", "5.9.1.1", "5.9.1.1")

	entry := waitForEntry(t, store, "5.9.1.1")

	if entry.Org != "Hetzner Online GmbH" {
		t.Errorf("Org = %q", entry.Org)
	}
	if entry.Source != SourceRDAP {
		t.Errorf("Source = %q, want %q", entry.Source, SourceRDAP)
	}
	wantExpiry := time.Now().Add(24 * time.Hour).Unix()
	if diff := entry.ExpiresAt - wantExpiry; diff > 5 || diff < -5 {
		t.Errorf("ExpiresAt is %d off the configured TTL", diff)
	}
	if got := hits.Load(); got != 1 {
		t.Errorf("registry requests = %d, want 1", got)
	}
}

func TestEnricherCachesFailuresBriefly(t *testing.T) {
	enricher, store := newTestEnricher(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go enricher.Start(ctx)

	enricher.Enqueue("203.0.114.5")
	entry := waitForEntry(t, store, "203.0.114.5")

	if entry.Source != SourceError {
		t.Errorf("Source = %q, want %q", entry.Source, SourceError)
	}
	if entry.LastError == "" {
		t.Error("LastError is empty, want the failure recorded")
	}
	// A failure must be retried in an hour, not on the next dashboard load.
	wantExpiry := time.Now().Add(errorTTL).Unix()
	if diff := entry.ExpiresAt - wantExpiry; diff > 5 || diff < -5 {
		t.Errorf("ExpiresAt is %d off the error TTL", diff)
	}
}

func TestEnricherCachesReservedRangesWithoutRequests(t *testing.T) {
	var hits atomic.Int32
	enricher, store := newTestEnricher(t, func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go enricher.Start(ctx)

	enricher.Enqueue("10.0.0.1")
	entry := waitForEntry(t, store, "10.0.0.1")

	if entry.Source != SourcePrivate {
		t.Errorf("Source = %q, want %q", entry.Source, SourcePrivate)
	}
	if hits.Load() != 0 {
		t.Errorf("registry requests = %d, want 0", hits.Load())
	}
}

func TestEnqueueStaleSkipsFreshEntries(t *testing.T) {
	var hits atomic.Int32
	fixture := loadFixture(t, "lacnic_minimal.json")

	enricher, store := newTestEnricher(t, func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		_, _ = w.Write(fixture)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go enricher.Start(ctx)

	now := time.Now().Unix()
	enricher.EnqueueStale([]storage.TopSourceIP{
		{SourceIP: "200.3.12.1", WhoisExpiresAt: now + 3600}, // fresh, skip
		{SourceIP: "200.3.12.2", WhoisExpiresAt: now - 1},    // expired, refresh
		{SourceIP: "200.3.12.3", WhoisExpiresAt: 0},          // never looked up
	})

	waitForEntry(t, store, "200.3.12.2")
	waitForEntry(t, store, "200.3.12.3")

	entries, err := store.GetIPWhois([]string{"200.3.12.1"})
	if err != nil {
		t.Fatalf("get ip whois: %v", err)
	}
	if _, found := entries["200.3.12.1"]; found {
		t.Error("a fresh entry was looked up again")
	}
	if got := hits.Load(); got != 2 {
		t.Errorf("registry requests = %d, want 2", got)
	}
}

func TestEnricherStopsOnContextCancel(t *testing.T) {
	enricher, _ := newTestEnricher(t, func(w http.ResponseWriter, _ *http.Request) {})

	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})
	go func() {
		enricher.Start(ctx)
		close(stopped)
	}()

	cancel()
	select {
	case <-stopped:
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not return after the context was cancelled")
	}
}

// The point of the range index: a report listing a dozen relays out of one
// assignment must cost one registry request, not a dozen.
func TestEnricherReusesAnswersAcrossARange(t *testing.T) {
	var hits atomic.Int32
	fixture := loadFixture(t, "ripe_individual_contacts.json")

	enricher, store := newTestEnricher(t, func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		_, _ = w.Write(fixture)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go enricher.Start(ctx)

	// Every address below sits in 80.67.31.0 - 80.67.31.255, the range the
	// fixture covers.
	ips := []string{"80.67.31.39", "80.67.31.100", "80.67.31.41", "80.67.31.35", "80.67.31.42"}
	enricher.Enqueue(ips...)

	for _, ip := range ips {
		entry := waitForEntry(t, store, ip)
		if entry.Org != "domainfactory" {
			t.Errorf("%s: Org = %q", ip, entry.Org)
		}
		if entry.CIDR != "80.67.31.0 - 80.67.31.255" {
			t.Errorf("%s: CIDR = %q", ip, entry.CIDR)
		}
	}

	// Workers run in parallel, so a second address may start before the first
	// answer is indexed; what must not happen is one request per address.
	if got := hits.Load(); got >= int32(len(ips)) {
		t.Errorf("registry requests = %d for %d addresses in one range", got, len(ips))
	}
}

// A restart must not re-query a registry for a range it already knows.
func TestEnricherSeedsRangesFromCache(t *testing.T) {
	var hits atomic.Int32
	fixture := loadFixture(t, "ripe_individual_contacts.json")

	enricher, store := newTestEnricher(t, func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		_, _ = w.Write(fixture)
	})

	if err := store.UpsertIPWhois(&storage.IPWhois{
		IP: "80.67.31.1", Org: "domainfactory", Network: "DOMAINFACTORY-20060601",
		CIDR: "80.67.31.0 - 80.67.31.255", Country: "DE", Source: SourceRDAP,
		LookedUpAt: time.Now().Unix(), ExpiresAt: time.Now().Add(time.Hour).Unix(),
	}); err != nil {
		t.Fatalf("seed cache: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go enricher.Start(ctx)

	enricher.Enqueue("80.67.31.77")
	entry := waitForEntry(t, store, "80.67.31.77")

	if entry.Org != "domainfactory" {
		t.Errorf("Org = %q", entry.Org)
	}
	if hits.Load() != 0 {
		t.Errorf("registry requests = %d, want 0 for a range already cached", hits.Load())
	}
}
