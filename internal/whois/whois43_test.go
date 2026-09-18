package whois

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
)

// ripeResponse is the shape RIPE answers an IP query with: the most specific
// object first, then the contact objects we deliberately ignore.
const ripeResponse = `% This is the RIPE Database query service.

inetnum:        5.9.1.0 - 5.9.1.31
netname:        HETZNER-fsn1-dc7
descr:          Hetzner Online GmbH
descr:          Datacenter fsn1-dc7
country:        DE
admin-c:        HOAC1-RIPE
status:         ASSIGNED PA
mnt-by:         HOS-GUN

role:           Hetzner Online GmbH - Contact Role
address:        Industriestrasse 25
phone:          +49 9831 505-0
`

// arinResponse uses entirely different field names for the same facts.
const arinResponse = `#
# ARIN WHOIS data and services are subject to the Terms of Use
#

NetRange:       209.85.128.0 - 209.85.255.255
CIDR:           209.85.128.0/17
NetName:        GOOGLE
Organization:   Google LLC (GOGL)
Country:        US
`

const ianaReferral = `% IANA WHOIS server

refer:          whois.ripe.net

inetnum:      5.0.0.0 - 5.255.255.255
organisation: RIPE NCC
`

// whoisTestServer serves canned responses, keyed by the server address the
// client believes it is talking to.
type whoisTestServer struct {
	listener  net.Listener
	responses map[string]string
	queries   atomic.Int32
}

func newWhoisTestServer(t *testing.T, responses map[string]string) *whoisTestServer {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	srv := &whoisTestServer{listener: ln, responses: responses}
	t.Cleanup(func() { _ = ln.Close() })

	go srv.serve()
	return srv
}

func (s *whoisTestServer) serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go func() {
			defer func() { _ = conn.Close() }()
			buf := make([]byte, 256)
			n, err := conn.Read(buf)
			if err != nil {
				return
			}
			s.queries.Add(1)
			// The dialer encodes which registry was asked in the query line.
			query := strings.TrimSpace(string(buf[:n]))
			for name, response := range s.responses {
				if strings.HasPrefix(query, name+" ") {
					_, _ = io.WriteString(conn, response)
					return
				}
			}
			_, _ = io.WriteString(conn, s.responses["default"])
		}()
	}
}

// dialer routes every WHOIS server to the test listener, prefixing the query
// with the server name so the handler can tell them apart.
func (s *whoisTestServer) dialer() Dialer {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			host = addr
		}
		conn, err := (&net.Dialer{}).DialContext(ctx, network, s.listener.Addr().String())
		if err != nil {
			return nil, err
		}
		return &prefixConn{Conn: conn, prefix: host + " "}, nil
	}
}

// prefixConn tags the outgoing query with the server it was meant for.
type prefixConn struct {
	net.Conn
	prefix string
	sent   bool
}

func (p *prefixConn) Write(b []byte) (int, error) {
	if p.sent {
		return p.Conn.Write(b)
	}
	// Send the tag and the query as one write: the test server answers and
	// closes after a single read, so a split write would race it.
	p.sent = true
	if _, err := p.Conn.Write(append([]byte(p.prefix), b...)); err != nil {
		return 0, err
	}
	return len(b), nil
}

// newFallbackClient builds a client whose RDAP endpoint always fails, so the
// WHOIS path is what gets exercised.
func newFallbackClient(t *testing.T, srv *whoisTestServer) *Client {
	t.Helper()
	return newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}, &fakeResolver{},
		WithWhoisFallback(true),
		WithDialer(srv.dialer()),
		WithWhoisBootstrap("whois.iana.org"),
	)
}

func TestWhoisFallbackFollowsReferral(t *testing.T) {
	srv := newWhoisTestServer(t, map[string]string{
		"whois.iana.org": ianaReferral,
		"whois.ripe.net": ripeResponse,
		"default":        "",
	})
	client := newFallbackClient(t, srv)

	got := client.Lookup(context.Background(), "5.9.1.1")

	if got.Source != SourceWHOIS {
		t.Fatalf("Source = %q, want %q (err: %s)", got.Source, SourceWHOIS, got.Err)
	}
	if got.Org != "Hetzner Online GmbH" {
		t.Errorf("Org = %q", got.Org)
	}
	if got.Network != "HETZNER-fsn1-dc7" {
		t.Errorf("Network = %q", got.Network)
	}
	if got.Country != "DE" {
		t.Errorf("Country = %q", got.Country)
	}
	if got.CIDR != "5.9.1.0 - 5.9.1.31" {
		t.Errorf("CIDR = %q", got.CIDR)
	}
	// A successful fallback is a result, not a failure: it must be cached for
	// the full TTL rather than retried in an hour.
	if got.Err != "" {
		t.Errorf("Err = %q, want empty after a successful fallback", got.Err)
	}
	if srv.queries.Load() != 2 {
		t.Errorf("whois queries = %d, want 2 (bootstrap then registry)", srv.queries.Load())
	}
}

func TestWhoisFallbackParsesARINFieldNames(t *testing.T) {
	srv := newWhoisTestServer(t, map[string]string{
		"whois.iana.org": "refer: whois.arin.net\n",
		"whois.arin.net": arinResponse,
		"default":        "",
	})
	client := newFallbackClient(t, srv)

	got := client.Lookup(context.Background(), "209.85.128.1")

	if got.Org != "Google LLC" {
		t.Errorf("Org = %q, want the ARIN handle stripped", got.Org)
	}
	if got.Network != "GOOGLE" {
		t.Errorf("Network = %q", got.Network)
	}
	if got.CIDR != "209.85.128.0/17" {
		t.Errorf("CIDR = %q", got.CIDR)
	}
	if got.Country != "US" {
		t.Errorf("Country = %q", got.Country)
	}
}

func TestWhoisFallbackNotUsedWhenRDAPAnswers(t *testing.T) {
	srv := newWhoisTestServer(t, map[string]string{"default": ripeResponse})
	fixture := loadFixture(t, "ripe_hetzner.json")

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(fixture)
	}, &fakeResolver{}, WithWhoisFallback(true), WithDialer(srv.dialer()),
		WithWhoisBootstrap("whois.iana.org"))

	got := client.Lookup(context.Background(), "5.9.1.1")

	if got.Source != SourceRDAP {
		t.Errorf("Source = %q, want %q", got.Source, SourceRDAP)
	}
	if srv.queries.Load() != 0 {
		t.Errorf("whois was queried %d times despite RDAP answering", srv.queries.Load())
	}
}

func TestWhoisFallbackCanBeDisabled(t *testing.T) {
	srv := newWhoisTestServer(t, map[string]string{"default": ripeResponse})

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}, &fakeResolver{}, WithDialer(srv.dialer()), WithWhoisFallback(false),
		WithWhoisBootstrap("whois.iana.org"))

	got := client.Lookup(context.Background(), "5.9.1.1")

	if got.Source == SourceWHOIS {
		t.Error("the fallback ran even though it was disabled")
	}
	if srv.queries.Load() != 0 {
		t.Errorf("whois was queried %d times with the fallback off", srv.queries.Load())
	}
}

func TestWhoisFallbackReportsBothFailures(t *testing.T) {
	srv := newWhoisTestServer(t, map[string]string{
		"whois.iana.org": "refer: whois.ripe.net\n",
		"whois.ripe.net": "% no entries found\n",
		"default":        "",
	})
	client := newFallbackClient(t, srv)

	got := client.Lookup(context.Background(), "5.9.1.1")

	if got.Source != SourceError {
		t.Errorf("Source = %q, want %q", got.Source, SourceError)
	}
	// The cached error should say what both attempts did, for debugging.
	for _, want := range []string{"rdap status 503", "whois fallback"} {
		if !strings.Contains(got.Err, want) {
			t.Errorf("Err = %q, missing %q", got.Err, want)
		}
	}
}

func TestParseWhoisFieldsSkipsBannersAndComments(t *testing.T) {
	fields := parseWhoisFields(ripeResponse)

	if v := firstValue(fields, "netname"); v != "HETZNER-fsn1-dc7" {
		t.Errorf("netname = %q", v)
	}
	// The first descr wins; the second is the datacenter, not the company.
	if v := firstValue(fields, "descr"); v != "Hetzner Online GmbH" {
		t.Errorf("descr = %q", v)
	}
	for _, f := range fields {
		if strings.HasPrefix(f.key, "%") {
			t.Errorf("comment line leaked into the fields: %+v", f)
		}
	}
}

func TestWhoisReferralFormats(t *testing.T) {
	tests := map[string]string{
		"refer:          whois.ripe.net":                    "whois.ripe.net:43",
		"ReferralServer: whois://whois.ripe.net":            "whois.ripe.net:43",
		"ReferralServer: rwhois://rwhois.example.com:4321/": "rwhois.example.com:4321",
	}
	for line, want := range tests {
		if got := whoisReferral(parseWhoisFields(line)); got != want {
			t.Errorf("whoisReferral(%q) = %q, want %q", line, got, want)
		}
	}
}
