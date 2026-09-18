package storage

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/meysam81/parse-dmarc/internal/parser"
)

func TestGetStatistics_HasData(t *testing.T) {
	// Create an in-memory SQLite database for testing
	storage, err := NewStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() { _ = storage.Close() }()

	t.Run("empty database", func(t *testing.T) {
		stats, err := storage.GetStatistics()
		if err != nil {
			t.Fatalf("Failed to get statistics: %v", err)
		}

		if stats.HasData {
			t.Errorf("Expected HasData to be false for empty database, got true")
		}

		if stats.TotalReports != 0 {
			t.Errorf("Expected TotalReports to be 0, got %d", stats.TotalReports)
		}

		if stats.TotalMessages != 0 {
			t.Errorf("Expected TotalMessages to be 0, got %d", stats.TotalMessages)
		}

		if stats.CompliantMessages != 0 {
			t.Errorf("Expected CompliantMessages to be 0, got %d", stats.CompliantMessages)
		}

		if stats.ComplianceRate != 0 {
			t.Errorf("Expected ComplianceRate to be 0, got %f", stats.ComplianceRate)
		}
	})

	t.Run("database with report", func(t *testing.T) {
		xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<feedback>
  <version>1.0</version>
  <report_metadata>
    <org_name>google.com</org_name>
    <email>noreply-dmarc-support@google.com</email>
    <report_id>12345678901234567890</report_id>
    <date_range>
      <begin>1609459200</begin>
      <end>1609545600</end>
    </date_range>
  </report_metadata>
  <policy_published>
    <domain>example.com</domain>
    <adkim>r</adkim>
    <aspf>r</aspf>
    <p>none</p>
    <sp>none</sp>
    <pct>100</pct>
  </policy_published>
  <record>
    <row>
      <source_ip>192.0.2.1</source_ip>
      <count>100</count>
      <policy_evaluated>
        <disposition>none</disposition>
        <dkim>pass</dkim>
        <spf>pass</spf>
      </policy_evaluated>
    </row>
    <identifiers>
      <header_from>example.com</header_from>
    </identifiers>
    <auth_results>
      <spf>
        <domain>example.com</domain>
        <result>pass</result>
      </spf>
      <dkim>
        <domain>example.com</domain>
        <result>pass</result>
      </dkim>
    </auth_results>
  </record>
</feedback>`

		feedback, err := parser.ParseReport([]byte(xmlData))
		if err != nil {
			t.Fatalf("Failed to parse report: %v", err)
		}

		err = storage.SaveReport(feedback)
		if err != nil {
			t.Fatalf("Failed to save report: %v", err)
		}

		stats, err := storage.GetStatistics()
		if err != nil {
			t.Fatalf("Failed to get statistics after adding report: %v", err)
		}

		if !stats.HasData {
			t.Errorf("Expected HasData to be true after adding report, got false")
		}

		if stats.TotalReports != 1 {
			t.Errorf("Expected TotalReports to be 1, got %d", stats.TotalReports)
		}

		if stats.TotalMessages != 100 {
			t.Errorf("Expected TotalMessages to be 100, got %d", stats.TotalMessages)
		}

		if stats.CompliantMessages != 100 {
			t.Errorf("Expected CompliantMessages to be 100, got %d", stats.CompliantMessages)
		}

		if stats.ComplianceRate != 100.0 {
			t.Errorf("Expected ComplianceRate to be 100.0, got %f", stats.ComplianceRate)
		}
	})
}

// The backend owns the status decision; the dashboard only renders it.
func TestReportSummaryClassify(t *testing.T) {
	cases := map[string]struct {
		total, compliant int
		wantStatus       string
		wantRate         string // "nil" or the formatted rate
	}{
		"null report": {0, 0, "empty", "nil"},
		"all pass":    {100, 100, "pass", "100.0"},
		"warn":        {100, 80, "warn", "80.0"},
		"fail":        {100, 79, "fail", "79.0"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			r := ReportSummary{TotalMessages: tc.total, CompliantMessages: tc.compliant}
			r.classify()

			gotRate := "nil"
			if r.ComplianceRate != nil {
				gotRate = fmt.Sprintf("%.1f", *r.ComplianceRate)
			}
			if r.Status != tc.wantStatus || gotRate != tc.wantRate {
				t.Fatalf("got status=%q rate=%s, want status=%q rate=%s", r.Status, gotRate, tc.wantStatus, tc.wantRate)
			}
		})
	}
}

// An RFC 7489 null report (one record, count=0, no auth results) must come
// back as status "empty" with a null compliance rate, not as a 0% failure.
func TestGetReports_NullReport(t *testing.T) {
	storage, err := NewStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() { _ = storage.Close() }()

	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<feedback>
  <version>1.0</version>
  <report_metadata>
    <org_name>example-reporter.tld</org_name>
    <email>dmarc-support@example-reporter.tld</email>
    <report_id>1780726747.321742392</report_id>
    <date_range>
      <begin>1780610400</begin>
      <end>1780696800</end>
    </date_range>
  </report_metadata>
  <policy_published>
    <domain>example.com</domain>
    <adkim>r</adkim>
    <aspf>r</aspf>
    <p>quarantine</p>
    <sp>quarantine</sp>
    <pct>100</pct>
  </policy_published>
  <record>
    <row>
      <source_ip></source_ip>
      <count>0</count>
      <policy_evaluated>
        <disposition></disposition>
        <dkim></dkim>
        <spf></spf>
      </policy_evaluated>
    </row>
    <identifiers>
      <header_from>example.com</header_from>
    </identifiers>
    <auth_results></auth_results>
  </record>
</feedback>`

	feedback, err := parser.ParseReport([]byte(xmlData))
	if err != nil {
		t.Fatalf("Failed to parse report: %v", err)
	}
	if err := storage.SaveReport(feedback); err != nil {
		t.Fatalf("Failed to save report: %v", err)
	}

	reports, err := storage.GetReports(10, 0)
	if err != nil {
		t.Fatalf("Failed to get reports: %v", err)
	}
	if len(reports) != 1 {
		t.Fatalf("Expected 1 report, got %d", len(reports))
	}
	if reports[0].Status != "empty" || reports[0].ComplianceRate != nil {
		t.Fatalf("got status=%q rate=%v, want status=\"empty\" rate=nil", reports[0].Status, reports[0].ComplianceRate)
	}

	// The API contract the dashboard and MCP clients see.
	out, err := json.Marshal(reports[0])
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"compliance_rate":null`, `"status":"empty"`} {
		if !strings.Contains(string(out), want) {
			t.Errorf("JSON %s\nmissing %s", out, want)
		}
	}
}

// Health is judged over delivered mail only: spoofing the receiver already
// blocked under the published policy is the policy working, not a gap.
func TestStatisticsClassify(t *testing.T) {
	cases := map[string]struct {
		reports, total, compliant, enforced int
		wantHealth                          string
		wantRate                            string // "nil" or the formatted rate
	}{
		"no reports":            {0, 0, 0, 0, "nodata", "nil"},
		"all delivered pass":    {1, 100, 100, 0, "secure", "100.0"},
		"blocked spoofing only": {1, 90, 0, 90, "secure", "nil"},
		"3 pass, 90 blocked":    {1, 93, 3, 90, "secure", "100.0"},
		"warning":               {1, 100, 85, 0, "warning", "85.0"},
		"critical":              {1, 100, 50, 0, "critical", "50.0"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			st := Statistics{
				TotalReports: tc.reports, TotalMessages: tc.total,
				CompliantMessages: tc.compliant, EnforcedMessages: tc.enforced,
				HasData: tc.reports > 0,
			}
			st.classify()

			gotRate := "nil"
			if st.DeliveredComplianceRate != nil {
				gotRate = fmt.Sprintf("%.1f", *st.DeliveredComplianceRate)
			}
			if st.Health != tc.wantHealth || gotRate != tc.wantRate {
				t.Fatalf("got health=%q rate=%s, want health=%q rate=%s", st.Health, gotRate, tc.wantHealth, tc.wantRate)
			}
		})
	}
}

// reportXML builds a one-domain aggregate report from (count, disposition,
// dkim, spf) rows so tests can synthesise traffic mixes compactly.
func reportXML(id, policy string, rows ...[4]string) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<?xml version="1.0"?><feedback><version>1.0</version>
<report_metadata><org_name>google.com</org_name><report_id>%s</report_id>
<date_range><begin>1609459200</begin><end>1609545600</end></date_range></report_metadata>
<policy_published><domain>example.com</domain><p>%s</p><pct>100</pct></policy_published>`, id, policy)
	for i, r := range rows {
		fmt.Fprintf(&b, `<record><row><source_ip>192.0.2.%d</source_ip><count>%s</count>
<policy_evaluated><disposition>%s</disposition><dkim>%s</dkim><spf>%s</spf></policy_evaluated></row>
<identifiers><header_from>example.com</header_from></identifiers></record>`, i+1, r[0], r[1], r[2], r[3])
	}
	b.WriteString("</feedback>")
	return b.String()
}

// A p=reject domain that blocks a botnet must read as healthy, and an
// unenforced domain letting failures through must still drag health down.
func TestGetStatistics_Health(t *testing.T) {
	storage, err := NewStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() { _ = storage.Close() }()

	save := func(xml string) {
		t.Helper()
		fb, err := parser.ParseReport([]byte(xml))
		if err != nil {
			t.Fatal(err)
		}
		if err := storage.SaveReport(fb); err != nil {
			t.Fatal(err)
		}
	}
	check := func(wantEnforced int, wantRate float64, wantHealth string) {
		t.Helper()
		stats, err := storage.GetStatistics()
		if err != nil {
			t.Fatal(err)
		}
		if stats.EnforcedMessages != wantEnforced || stats.Health != wantHealth ||
			stats.DeliveredComplianceRate == nil || fmt.Sprintf("%.1f", *stats.DeliveredComplianceRate) != fmt.Sprintf("%.1f", wantRate) {
			t.Fatalf("got enforced=%d rate=%v health=%q, want enforced=%d rate=%.1f health=%q",
				stats.EnforcedMessages, stats.DeliveredComplianceRate, stats.Health, wantEnforced, wantRate, wantHealth)
		}
	}

	// p=reject: 10 legitimate messages pass, 40 spoofed ones are rejected.
	save(reportXML("r1", "reject",
		[4]string{"40", "reject", "fail", "fail"},
		[4]string{"10", "none", "pass", "pass"},
	))
	check(40, 100, "secure")

	// p=none: 50 of 100 fail and are delivered anyway.
	save(reportXML("r2", "none",
		[4]string{"50", "none", "fail", "fail"},
		[4]string{"50", "none", "pass", "pass"},
	))
	// delivered = 150 - 40 = 110, compliant = 60
	check(40, 60.0/110*100, "critical")

	stats, _ := storage.GetStatistics()
	if stats.ComplianceRate != 40 { // raw rate over all mail is unchanged for metrics
		t.Errorf("ComplianceRate = %v, want 40", stats.ComplianceRate)
	}
	out, _ := json.Marshal(stats)
	for _, want := range []string{`"enforced_messages":40`, `"health":"critical"`} {
		if !strings.Contains(string(out), want) {
			t.Errorf("JSON %s\nmissing %s", out, want)
		}
	}
}

func TestIPWhois_UpsertAndGet(t *testing.T) {
	storage, err := NewStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() { _ = storage.Close() }()

	entry := &IPWhois{
		IP: "203.0.114.5", Org: "Example Org", Network: "EXAMPLE-NET",
		CIDR: "203.0.114.0/24", Country: "DE", Hostname: "mail.example.com",
		Source: "rdap", LookedUpAt: 1700000000, ExpiresAt: 1700600000,
	}
	if err := storage.UpsertIPWhois(entry); err != nil {
		t.Fatalf("UpsertIPWhois: %v", err)
	}

	got, err := storage.GetIPWhois([]string{"203.0.114.5", "203.0.114.6"})
	if err != nil {
		t.Fatalf("GetIPWhois: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d entries, want 1 (an uncached IP must simply be absent)", len(got))
	}
	if got["203.0.114.5"].Org != "Example Org" {
		t.Errorf("Org = %q", got["203.0.114.5"].Org)
	}

	// A refreshed lookup must overwrite in place, not accumulate rows.
	entry.Org = "Renamed Org"
	entry.LookedUpAt = 1700700000
	if err := storage.UpsertIPWhois(entry); err != nil {
		t.Fatalf("UpsertIPWhois (refresh): %v", err)
	}
	got, err = storage.GetIPWhois([]string{"203.0.114.5"})
	if err != nil {
		t.Fatalf("GetIPWhois: %v", err)
	}
	if len(got) != 1 || got["203.0.114.5"].Org != "Renamed Org" {
		t.Errorf("got %+v, want a single updated row", got)
	}
}

// Enrichment is additive: an unenriched source must serialize exactly as it
// did before the feature existed, so older clients see no change.
func TestGetTopSourceIPs_Whois(t *testing.T) {
	storage, err := NewStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() { _ = storage.Close() }()

	fb, err := parser.ParseReport([]byte(reportXML("r1", "none",
		[4]string{"10", "none", "pass", "pass"},
		[4]string{"5", "none", "fail", "fail"},
		[4]string{"1", "none", "fail", "fail"},
	)))
	if err != nil {
		t.Fatalf("ParseReport: %v", err)
	}
	if err := storage.SaveReport(fb); err != nil {
		t.Fatalf("SaveReport: %v", err)
	}

	mustUpsert := func(w *IPWhois) {
		t.Helper()
		if err := storage.UpsertIPWhois(w); err != nil {
			t.Fatalf("UpsertIPWhois: %v", err)
		}
	}
	mustUpsert(&IPWhois{
		IP: "192.0.2.1", Org: "Example Org", Network: "EXAMPLE-NET",
		Country: "DE", Hostname: "mail.example.com", Source: "rdap",
		LookedUpAt: 1700000000, ExpiresAt: 1700600000,
	})
	// A failed lookup is cached to avoid retrying it, but has nothing to show.
	mustUpsert(&IPWhois{
		IP: "192.0.2.3", Source: "error", LastError: "rdap status 500",
		LookedUpAt: 1700000000, ExpiresAt: 1700003600,
	})

	sources, err := storage.GetTopSourceIPs(10)
	if err != nil {
		t.Fatalf("GetTopSourceIPs: %v", err)
	}
	if len(sources) != 3 {
		t.Fatalf("got %d sources, want 3", len(sources))
	}

	byIP := make(map[string]TopSourceIP, len(sources))
	for _, s := range sources {
		byIP[s.SourceIP] = s
	}

	enriched := byIP["192.0.2.1"]
	if enriched.Whois == nil {
		t.Fatal("192.0.2.1 has no whois data")
	}
	if enriched.Whois.Org != "Example Org" || enriched.Whois.Hostname != "mail.example.com" {
		t.Errorf("whois = %+v", enriched.Whois)
	}
	if enriched.Count != 10 {
		t.Errorf("Count = %d, want 10 (the join must not disturb the aggregate)", enriched.Count)
	}

	if byIP["192.0.2.2"].Whois != nil {
		t.Error("an IP with no cache entry must have no whois object")
	}
	if byIP["192.0.2.2"].WhoisExpiresAt != 0 {
		t.Error("an IP with no cache entry must look stale to the enricher")
	}

	failed := byIP["192.0.2.3"]
	if failed.Whois != nil {
		t.Error("a failed lookup must not produce an empty whois object")
	}
	if failed.WhoisExpiresAt == 0 {
		t.Error("a failed lookup must still be remembered so it is not retried at once")
	}

	// The JSON contract: absent, not null, when there is nothing to show.
	out, err := json.Marshal(byIP["192.0.2.2"])
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "whois") {
		t.Errorf("JSON %s must omit the whois key entirely", out)
	}
	if out, err := json.Marshal(enriched); err != nil {
		t.Fatal(err)
	} else if !strings.Contains(string(out), `"whois":{`) {
		t.Errorf("JSON %s is missing the whois object", out)
	}
}

// An entry written by superseded lookup logic must neither be shown nor kept:
// the fix that stopped naming private individuals would otherwise take a week
// to reach an existing installation.
func TestGetTopSourceIPs_SupersededCacheVersion(t *testing.T) {
	storage, err := NewStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() { _ = storage.Close() }()

	fb, err := parser.ParseReport([]byte(reportXML("r1", "none",
		[4]string{"10", "none", "pass", "pass"},
	)))
	if err != nil {
		t.Fatalf("ParseReport: %v", err)
	}
	if err := storage.SaveReport(fb); err != nil {
		t.Fatalf("SaveReport: %v", err)
	}

	// Written the way an older release would have: far from expiry, but by
	// logic we have since replaced.
	future := time.Now().Add(7 * 24 * time.Hour).Unix()
	_, err = storage.db.Exec(`
		INSERT INTO ip_whois (ip, org, network, cidr, country, hostname, source,
			last_error, looked_up_at, expires_at, lookup_version)
		VALUES (?, ?, '', '', 'DE', '', 'rdap', '', ?, ?, ?)
	`, "192.0.2.1", "A Named Person", time.Now().Unix(), future, WhoisCacheVersion-1)
	if err != nil {
		t.Fatalf("insert legacy row: %v", err)
	}

	sources, err := storage.GetTopSourceIPs(10)
	if err != nil {
		t.Fatalf("GetTopSourceIPs: %v", err)
	}
	if len(sources) != 1 {
		t.Fatalf("got %d sources, want 1", len(sources))
	}

	if sources[0].Whois != nil {
		t.Errorf("superseded entry was shown: %+v", sources[0].Whois)
	}
	if !sources[0].WhoisStale(time.Now().Unix()) {
		t.Error("superseded entry was not treated as stale, so it is never refreshed")
	}

	// Once re-resolved at the current version it shows again.
	if err := storage.UpsertIPWhois(&IPWhois{
		IP: "192.0.2.1", Org: "Example Org", Source: "rdap",
		LookedUpAt: time.Now().Unix(), ExpiresAt: future,
	}); err != nil {
		t.Fatalf("UpsertIPWhois: %v", err)
	}
	sources, err = storage.GetTopSourceIPs(10)
	if err != nil {
		t.Fatalf("GetTopSourceIPs: %v", err)
	}
	if sources[0].Whois == nil || sources[0].Whois.Org != "Example Org" {
		t.Errorf("re-resolved entry not shown: %+v", sources[0].Whois)
	}
	if sources[0].WhoisStale(time.Now().Unix()) {
		t.Error("a freshly written entry must not look stale")
	}
}

// A database created before the column existed must gain it on start.
func TestMigrateAddsLookupVersion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy.sqlite")

	storage, err := NewStorage(path)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	if _, err := storage.db.Exec(`ALTER TABLE ip_whois DROP COLUMN lookup_version`); err != nil {
		t.Skipf("this SQLite build cannot drop a column: %v", err)
	}
	if err := storage.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	// Reopening runs the migration; a second open must not fail on the column
	// already being there.
	for attempt := range 2 {
		reopened, err := NewStorage(path)
		if err != nil {
			t.Fatalf("reopen %d: %v", attempt, err)
		}
		if err := reopened.UpsertIPWhois(&IPWhois{
			IP: "192.0.2.9", Org: "Example Org", Source: "rdap",
			LookedUpAt: 1700000000, ExpiresAt: 1700600000,
		}); err != nil {
			t.Fatalf("upsert after migration: %v", err)
		}
		if err := reopened.Close(); err != nil {
			t.Fatalf("close: %v", err)
		}
	}
}

// Display decides what the report detail may show for one address.
func TestIPWhoisDisplay(t *testing.T) {
	usable := IPWhois{
		IP: "203.0.114.5", Org: "Example Org", Network: "EXAMPLE-NET",
		Country: "DE", Hostname: "mail.example.com", Source: "rdap",
		LookedUpAt: 1700000000, ExpiresAt: 1700600000,
		LookupVersion: WhoisCacheVersion,
	}

	tests := map[string]struct {
		mutate func(*IPWhois)
		want   bool
	}{
		"a registry answer is shown":      {func(*IPWhois) {}, true},
		"reverse DNS alone is shown":      {func(w *IPWhois) { w.Source = "rdns"; w.Org = ""; w.Network = ""; w.Country = "" }, true},
		"a failure is not":                {func(w *IPWhois) { w.Source = "error"; w.LastError = "rdap status 500" }, false},
		"a reserved range is not":         {func(w *IPWhois) { w.Source = "private" }, false},
		"superseded lookup logic is not":  {func(w *IPWhois) { w.LookupVersion = WhoisCacheVersion - 1 }, false},
		"an entry with no content is not": {func(w *IPWhois) { w.Org, w.Network, w.Country, w.Hostname = "", "", "", "" }, false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			entry := usable
			tt.mutate(&entry)

			shown := entry.Display()
			if (shown != nil) != tt.want {
				t.Fatalf("Display() = %+v, want shown=%v", shown, tt.want)
			}
			if shown != nil && shown.Source != entry.Source {
				t.Errorf("Source = %q, want %q: a reader must see where it came from",
					shown.Source, entry.Source)
			}
		})
	}
}

func TestIPWhoisStale(t *testing.T) {
	now := int64(1700000000)

	fresh := IPWhois{ExpiresAt: now + 60, LookupVersion: WhoisCacheVersion}
	if fresh.Stale(now) {
		t.Error("an unexpired current entry must not be stale")
	}

	expired := IPWhois{ExpiresAt: now - 1, LookupVersion: WhoisCacheVersion}
	if !expired.Stale(now) {
		t.Error("an expired entry must be stale")
	}

	superseded := IPWhois{ExpiresAt: now + 86400, LookupVersion: WhoisCacheVersion - 1}
	if !superseded.Stale(now) {
		t.Error("an entry from older lookup logic must be stale however long it has left")
	}
}
