package whois

import (
	"os"
	"path/filepath"
	"testing"
)

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return data
}

func TestParseRDAP(t *testing.T) {
	tests := map[string]struct {
		fixture string
		// why documents the registry quirk each fixture stands for.
		why     string
		want    Info
		wantErr bool
	}{
		"ARIN: registrant vCard, no ISO country anywhere": {
			fixture: "arin_google.json",
			why:     "ARIN ships the address as a free-text label, not a code",
			want: Info{
				Org:     "Google LLC",
				Network: "GOOGLE",
				CIDR:    "209.85.128.0/17",
				Country: "",
			},
		},
		"RIPE: lowercase country and structured address": {
			fixture: "ripe_hetzner.json",
			why:     "country must be upper-cased; registrant outranks the admin role",
			want: Info{
				Org:     "Hetzner Online GmbH",
				Network: "DE-HETZNER-20110222",
				CIDR:    "5.9.0.0 - 5.9.255.255",
				Country: "DE",
			},
		},
		"RIPE: a registrant named after its own handle is skipped": {
			fixture: "ripe_placeholder_registrant.json",
			why:     "the readable company name sits on a nested organisation entity",
			want: Info{
				Org:     "Hetzner Online GmbH",
				Network: "HETZNER-fsn1-dc7",
				CIDR:    "5.9.1.0 - 5.9.1.31",
				Country: "DE",
			},
		},
		"APNIC: no registrant, admin outranks technical": {
			fixture: "apnic_admin_only.json",
			want: Info{
				Org:     "APNIC Research and Development",
				Network: "APNIC-LABS",
				CIDR:    "1.1.1.0 - 1.1.1.255",
				Country: "AU",
			},
		},
		"minimal: no entities at all falls back to the network name": {
			fixture: "lacnic_minimal.json",
			want: Info{
				Org:     "LACNIC-EXAMPLE",
				Network: "LACNIC-EXAMPLE",
				CIDR:    "200.3.12.0 - 200.3.15.255",
			},
		},
		"malformed vCards are skipped, never fatal": {
			fixture: "broken_vcard.json",
			want: Info{
				Org:     "BROKEN-NET",
				Network: "BROKEN-NET",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := parseRDAP(loadFixture(t, tt.fixture))
			if err != nil {
				t.Fatalf("parseRDAP: %v", err)
			}
			if got.Org != tt.want.Org {
				t.Errorf("Org = %q, want %q", got.Org, tt.want.Org)
			}
			if got.Network != tt.want.Network {
				t.Errorf("Network = %q, want %q", got.Network, tt.want.Network)
			}
			if got.CIDR != tt.want.CIDR {
				t.Errorf("CIDR = %q, want %q", got.CIDR, tt.want.CIDR)
			}
			if got.Country != tt.want.Country {
				t.Errorf("Country = %q, want %q", got.Country, tt.want.Country)
			}
		})
	}
}

func TestParseRDAPRejectsNonJSON(t *testing.T) {
	if _, err := parseRDAP([]byte("<html>rate limited</html>")); err == nil {
		t.Fatal("expected an error for a non-JSON body")
	}
}

func TestNormalizeCountry(t *testing.T) {
	tests := map[string]string{
		"de":            "DE",
		"  us  ":        "US",
		"DE":            "DE",
		"United States": "",
		"D":             "",
		"12":            "",
		"":              "",
	}
	for in, want := range tests {
		if got := normalizeCountry(in); got != want {
			t.Errorf("normalizeCountry(%q) = %q, want %q", in, got, want)
		}
	}
}
