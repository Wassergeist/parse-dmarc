package whois

import (
	"net"
	"testing"
)

func TestParseRangeAcceptsBothRegistryForms(t *testing.T) {
	tests := map[string]struct {
		in           string
		member       string
		nonMember    string
		wantUnusable bool
	}{
		"CIDR prefix": {
			in: "209.85.128.0/17", member: "209.85.220.41", nonMember: "209.86.0.1",
		},
		"first - last, as RIPE writes an inetnum": {
			in: "80.67.31.0 - 80.67.31.255", member: "80.67.31.42", nonMember: "80.67.32.1",
		},
		"IPv6 prefix": {
			in: "2a01:4f8::/29", member: "2a01:4f8:1::1", nonMember: "2a02::1",
		},
		"nothing usable": {in: "", wantUnusable: true},
		"not a range":    {in: "DOMAINFACTORY", wantUnusable: true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			first, last, ok := parseRange(tt.in)
			if tt.wantUnusable {
				if ok {
					t.Fatalf("parseRange(%q) accepted a value that names no range", tt.in)
				}
				return
			}
			if !ok {
				t.Fatalf("parseRange(%q) failed", tt.in)
			}
			if !withinRange(net.ParseIP(tt.member), first, last) {
				t.Errorf("%s should be inside %s", tt.member, tt.in)
			}
			if withinRange(net.ParseIP(tt.nonMember), first, last) {
				t.Errorf("%s should be outside %s", tt.nonMember, tt.in)
			}
		})
	}
}

func TestRangeIndexPrefersTheMostRecentMatch(t *testing.T) {
	var idx rangeIndex

	// A wide allocation is usually seen before the assignment inside it; the
	// more specific answer must win.
	idx.Add(Info{Org: "RIPE NCC", CIDR: "80.64.0.0/11", Source: SourceRDAP})
	idx.Add(Info{Org: "domainfactory", CIDR: "80.67.31.0 - 80.67.31.255", Source: SourceRDAP})

	got, ok := idx.Lookup(net.ParseIP("80.67.31.42"))
	if !ok {
		t.Fatal("address in a known range was not found")
	}
	if got.Org != "domainfactory" {
		t.Errorf("Org = %q, want the more specific answer", got.Org)
	}

	if _, ok := idx.Lookup(net.ParseIP("9.9.9.9")); ok {
		t.Error("an address outside every known range was matched")
	}
}

func TestRangeIndexIgnoresAnswersWithoutARange(t *testing.T) {
	var idx rangeIndex
	idx.Add(Info{Org: "Example", CIDR: "", Source: SourceRDNS})

	if idx.Len() != 0 {
		t.Error("an answer covering no range must not be reusable")
	}
}

// An IPv4 address written as IPv4-in-IPv6 must compare equal to the plain one.
func TestWithinRangeNormalizesAddressForm(t *testing.T) {
	first, last, ok := parseRange("80.67.31.0/24")
	if !ok {
		t.Fatal("parseRange failed")
	}
	if !withinRange(net.ParseIP("::ffff:80.67.31.42"), first, last) {
		t.Error("IPv4-mapped address was not recognised as inside the range")
	}
}
