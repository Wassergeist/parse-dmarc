package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/meysam81/parse-dmarc/internal/parser"
)

type Storage struct {
	db *sql.DB
}

type ReportSummary struct {
	ID                int64  `json:"id"`
	ReportID          string `json:"report_id"`
	OrgName           string `json:"org_name"`
	Domain            string `json:"domain"`
	DateBegin         int64  `json:"date_begin"`
	DateEnd           int64  `json:"date_end"`
	TotalMessages     int    `json:"total_messages"`
	CompliantMessages int    `json:"compliant_messages"`
	// ComplianceRate is null for an RFC 7489 null report (no messages):
	// 0 of 0 is not a 0% pass rate.
	ComplianceRate *float64 `json:"compliance_rate"`
	PolicyP        string   `json:"policy_p"`
	// Status is the backend's verdict for display: empty, pass, warn or fail.
	Status string `json:"status"`
}

// Compliance-rate thresholds (percent) behind Status.
const (
	passRate = 100
	warnRate = 80
)

// classify derives ComplianceRate and Status from the message counts.
func (r *ReportSummary) classify() {
	if r.TotalMessages == 0 {
		r.Status = "empty"
		return
	}
	rate := float64(r.CompliantMessages) / float64(r.TotalMessages) * 100
	r.ComplianceRate = &rate
	switch {
	case rate >= passRate:
		r.Status = "pass"
	case rate >= warnRate:
		r.Status = "warn"
	default:
		r.Status = "fail"
	}
}

type Statistics struct {
	TotalReports      int     `json:"total_reports"`
	TotalMessages     int     `json:"total_messages"`
	CompliantMessages int     `json:"compliant_messages"`
	ComplianceRate    float64 `json:"compliance_rate"`
	// EnforcedMessages counts non-compliant mail the receiver already blocked
	// (disposition reject or quarantine): the published policy working, not
	// an authentication gap.
	EnforcedMessages int `json:"enforced_messages"`
	// DeliveredComplianceRate is the pass rate over mail that was actually
	// delivered (total minus enforced). Null when nothing was delivered.
	DeliveredComplianceRate *float64 `json:"delivered_compliance_rate"`
	// Health is the backend's verdict for the dashboard: nodata, secure,
	// warning or critical.
	Health          string `json:"health"`
	UniqueSourceIPs int    `json:"unique_source_ips"`
	UniqueDomains   int    `json:"unique_domains"`
	HasData         bool   `json:"has_data"`
}

// Delivered-compliance thresholds (percent) behind Health.
const (
	secureRate  = 95
	warningRate = 80
)

// classify derives DeliveredComplianceRate and Health from the counts.
func (st *Statistics) classify() {
	if !st.HasData {
		st.Health = "nodata"
		return
	}
	delivered := st.TotalMessages - st.EnforcedMessages
	if delivered <= 0 {
		// Everything seen was blocked spoofing: nothing unauthenticated got through.
		st.Health = "secure"
		return
	}
	rate := float64(st.CompliantMessages) / float64(delivered) * 100
	st.DeliveredComplianceRate = &rate
	switch {
	case rate >= secureRate:
		st.Health = "secure"
	case rate >= warningRate:
		st.Health = "warning"
	default:
		st.Health = "critical"
	}
}

type TopSourceIP struct {
	SourceIP string `json:"source_ip"`
	Count    int    `json:"count"`
	Pass     int    `json:"pass"`
	Fail     int    `json:"fail"`
	// Whois is nil until a lookup has produced something worth showing, so a
	// client that predates enrichment sees exactly the JSON it saw before.
	Whois *SourceWhois `json:"whois,omitempty"`
	// WhoisPending marks a source whose lookup is still outstanding, so the
	// dashboard can tell "being resolved" apart from "nothing to show". It is
	// set by the API layer, which knows whether enrichment runs at all.
	WhoisPending bool `json:"whois_pending,omitempty"`
	// WhoisExpiresAt is 0 when nothing is cached for this IP. It lets the
	// enricher spot stale entries without a second query and is deliberately
	// kept out of the API contract.
	WhoisExpiresAt int64 `json:"-"`
}

// WhoisStale reports whether this source still needs a lookup at time now
// (unix seconds). Missing and expired entries are both stale; it is the one
// definition the enricher and the API both use.
func (t *TopSourceIP) WhoisStale(now int64) bool {
	return t.WhoisExpiresAt < now
}

// SourceWhois is the display-oriented subset of a cached lookup.
type SourceWhois struct {
	Org      string `json:"org,omitempty"`
	Network  string `json:"network,omitempty"`
	CIDR     string `json:"cidr,omitempty"`
	Country  string `json:"country,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	// Source names where this came from: rdap, whois or rdns. Shown so a
	// reader can judge the answer, registries disagreeing with each other
	// being a normal state of affairs.
	Source     string `json:"source,omitempty"`
	LookedUpAt int64  `json:"looked_up_at,omitempty"`
}

// IPWhois is one cached ownership lookup for a sending IP. Empty fields mean
// the registry did not publish that detail, not that the lookup failed.
type IPWhois struct {
	IP       string
	Org      string
	Network  string
	CIDR     string
	Country  string
	Hostname string
	// Source records how the entry was produced: rdap, rdns, private or error.
	Source     string
	LastError  string
	LookedUpAt int64
	ExpiresAt  int64
	// LookupVersion is the WhoisCacheVersion that produced this entry.
	LookupVersion int
}

func (s *Storage) SaveReport(feedback *parser.Feedback) error {
	rawReport, err := json.Marshal(feedback)
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	result, err := tx.Exec(`
		INSERT OR IGNORE INTO reports (
			report_id, org_name, email, domain,
			date_begin, date_end, created_at,
			policy_p, policy_sp, policy_pct,
			total_messages, compliant_messages,
			raw_report
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		feedback.ReportMetadata.ReportID,
		feedback.ReportMetadata.OrgName,
		feedback.ReportMetadata.Email,
		feedback.PolicyPublished.Domain,
		feedback.ReportMetadata.DateRange.Begin,
		feedback.ReportMetadata.DateRange.End,
		time.Now().Unix(),
		feedback.PolicyPublished.P,
		feedback.PolicyPublished.SP,
		feedback.PolicyPublished.PCT,
		feedback.GetTotalMessages(),
		feedback.GetDMARCCompliantCount(),
		rawReport,
	)

	if err != nil {
		return fmt.Errorf("failed to insert report: %w", err)
	}

	reportID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert ID: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil
	}

	for _, record := range feedback.Records {
		dkimDomains, _ := json.Marshal(record.AuthResults.DKIM)
		spfDomains, _ := json.Marshal(record.AuthResults.SPF)

		_, err := tx.Exec(`
			INSERT INTO records (
				report_id, source_ip, count,
				disposition, dkim_result, spf_result,
				header_from, envelope_from,
				dkim_domains, spf_domains
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			reportID,
			record.Row.SourceIP,
			record.Row.Count,
			record.Row.PolicyEvaluated.Disposition,
			record.Row.PolicyEvaluated.DKIM,
			record.Row.PolicyEvaluated.SPF,
			record.Identifiers.HeaderFrom,
			record.Identifiers.EnvelopeFrom,
			dkimDomains,
			spfDomains,
		)

		if err != nil {
			return fmt.Errorf("failed to insert record: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *Storage) GetReports(limit, offset int) ([]ReportSummary, error) {
	rows, err := s.db.Query(`
		SELECT id, report_id, org_name, domain,
		       date_begin, date_end,
		       total_messages, compliant_messages,
		       policy_p
		FROM reports
		ORDER BY date_begin DESC
		LIMIT ? OFFSET ?
	`, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("query reports: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var reports []ReportSummary
	for rows.Next() {
		var r ReportSummary
		err := rows.Scan(
			&r.ID, &r.ReportID, &r.OrgName, &r.Domain,
			&r.DateBegin, &r.DateEnd,
			&r.TotalMessages, &r.CompliantMessages,
			&r.PolicyP,
		)
		if err != nil {
			return nil, fmt.Errorf("scan report row: %w", err)
		}

		r.classify()
		reports = append(reports, r)
	}

	return reports, nil
}

func (s *Storage) GetReportByID(id int64) (*parser.Feedback, error) {
	var rawReport string
	err := s.db.QueryRow("SELECT raw_report FROM reports WHERE id = ?", id).Scan(&rawReport)
	if err != nil {
		return nil, fmt.Errorf("query report %d: %w", id, err)
	}

	var feedback parser.Feedback
	if err := json.Unmarshal([]byte(rawReport), &feedback); err != nil {
		return nil, fmt.Errorf("unmarshal report %d: %w", id, err)
	}

	return &feedback, nil
}

func (s *Storage) GetStatistics() (*Statistics, error) {
	var stats Statistics

	err := s.db.QueryRow(`
		SELECT
			COUNT(*) as total_reports,
			COALESCE(SUM(total_messages), 0) as total_messages,
			COALESCE(SUM(compliant_messages), 0) as compliant_messages
		FROM reports
	`).Scan(&stats.TotalReports, &stats.TotalMessages, &stats.CompliantMessages)

	if err != nil {
		return nil, fmt.Errorf("query report statistics: %w", err)
	}

	stats.HasData = stats.TotalReports > 0

	if stats.TotalMessages > 0 {
		stats.ComplianceRate = float64(stats.CompliantMessages) / float64(stats.TotalMessages) * 100
	}

	err = s.db.QueryRow(`
		SELECT COALESCE(SUM(count), 0)
		FROM records
		WHERE disposition IN ('reject', 'quarantine')
		  AND dkim_result != 'pass'
		  AND spf_result != 'pass'
	`).Scan(&stats.EnforcedMessages)
	if err != nil {
		return nil, fmt.Errorf("query enforced messages: %w", err)
	}

	err = s.db.QueryRow("SELECT COUNT(DISTINCT source_ip) FROM records").Scan(&stats.UniqueSourceIPs)
	if err != nil {
		return nil, fmt.Errorf("query unique source IPs: %w", err)
	}

	err = s.db.QueryRow("SELECT COUNT(DISTINCT domain) FROM reports").Scan(&stats.UniqueDomains)
	if err != nil {
		return nil, fmt.Errorf("query unique domains: %w", err)
	}

	stats.classify()
	return &stats, nil
}

// GetTopSourceIPs returns the busiest sending IPs, each with whatever
// ownership data is cached for it. Stale cache entries are still returned:
// a slightly old owner beats a blank row while the enricher refreshes it.
func (s *Storage) GetTopSourceIPs(limit int) ([]TopSourceIP, error) {
	rows, err := s.db.Query(`
		SELECT
			t.source_ip, t.total_count, t.pass_count, t.fail_count,
			COALESCE(w.org, ''), COALESCE(w.network, ''), COALESCE(w.cidr, ''),
			COALESCE(w.country, ''), COALESCE(w.hostname, ''), COALESCE(w.source, ''),
			COALESCE(w.looked_up_at, 0),
			-- An entry from an older lookup version is expired by definition,
			-- so an upgrade re-resolves it instead of serving it for days.
			CASE WHEN COALESCE(w.lookup_version, 0) < ? THEN 0 ELSE COALESCE(w.expires_at, 0) END,
			COALESCE(w.lookup_version, 0)
		FROM (
			SELECT
				source_ip,
				SUM(count) as total_count,
				SUM(CASE WHEN (dkim_result = 'pass' OR spf_result = 'pass') THEN count ELSE 0 END) as pass_count,
				SUM(CASE WHEN (dkim_result != 'pass' AND spf_result != 'pass') THEN count ELSE 0 END) as fail_count
			FROM records
			GROUP BY source_ip
			ORDER BY total_count DESC
			LIMIT ?
		) t
		LEFT JOIN ip_whois w ON w.ip = t.source_ip
		ORDER BY t.total_count DESC
	`,
		// Bind order follows the query text: the version in the SELECT list
		// comes before the subquery's LIMIT.
		WhoisCacheVersion, limit)

	if err != nil {
		return nil, fmt.Errorf("query top source IPs: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []TopSourceIP
	for rows.Next() {
		var r TopSourceIP
		var w SourceWhois
		var source string
		var version int
		if err := rows.Scan(
			&r.SourceIP, &r.Count, &r.Pass, &r.Fail,
			&w.Org, &w.Network, &w.CIDR, &w.Country, &w.Hostname, &source,
			&w.LookedUpAt, &r.WhoisExpiresAt, &version,
		); err != nil {
			return nil, fmt.Errorf("scan source IP row: %w", err)
		}
		if version < WhoisCacheVersion {
			// Superseded logic: show nothing rather than an answer we would
			// no longer produce. The row above already marks it stale.
			source = ""
		}
		// A failed or private lookup is cached to stop us retrying it, but it
		// has nothing to display: leave Whois nil rather than emit an empty
		// object the frontend would have to special-case.
		if source != "" && source != whoisSourceError && source != whoisSourcePrivate {
			if w.Org != "" || w.Network != "" || w.Hostname != "" || w.Country != "" {
				w.Source = source
				r.Whois = &w
			}
		}
		results = append(results, r)
	}

	return results, nil
}

// Lookup outcomes stored in ip_whois.source. Kept here rather than in the
// whois package so storage does not depend on it (whois depends on storage).
const (
	whoisSourceRDAP    = "rdap"
	whoisSourceWHOIS   = "whois"
	whoisSourcePrivate = "private"
	whoisSourceError   = "error"
)

// UpsertIPWhois writes or refreshes the cache entry for one IP.
func (s *Storage) UpsertIPWhois(w *IPWhois) error {
	_, err := s.db.Exec(`
		INSERT INTO ip_whois (
			ip, org, network, cidr, country, hostname, source, last_error,
			looked_up_at, expires_at, lookup_version
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(ip) DO UPDATE SET
			org = excluded.org,
			network = excluded.network,
			cidr = excluded.cidr,
			country = excluded.country,
			hostname = excluded.hostname,
			source = excluded.source,
			last_error = excluded.last_error,
			looked_up_at = excluded.looked_up_at,
			expires_at = excluded.expires_at,
			lookup_version = excluded.lookup_version
	`, w.IP, w.Org, w.Network, w.CIDR, w.Country, w.Hostname, w.Source, w.LastError,
		w.LookedUpAt, w.ExpiresAt, WhoisCacheVersion)
	if err != nil {
		return fmt.Errorf("upsert ip whois %s: %w", w.IP, err)
	}
	return nil
}

// GetFreshWhoisRanges returns cached registry answers that still carry a usable
// address range, so a restart does not have to re-query a range it already
// knows. Reverse-DNS-only entries are excluded: a PTR record says nothing
// about the addresses next to it.
func (s *Storage) GetFreshWhoisRanges(now int64) ([]IPWhois, error) {
	rows, err := s.db.Query(`
		SELECT ip, org, network, cidr, country, hostname, source, last_error,
			looked_up_at, expires_at, lookup_version
		FROM ip_whois
		WHERE cidr != '' AND expires_at > ? AND lookup_version = ?
			AND source IN (?, ?)
		ORDER BY looked_up_at
	`, now, WhoisCacheVersion, whoisSourceRDAP, whoisSourceWHOIS)
	if err != nil {
		return nil, fmt.Errorf("query whois ranges: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var result []IPWhois
	for rows.Next() {
		var w IPWhois
		if err := rows.Scan(&w.IP, &w.Org, &w.Network, &w.CIDR, &w.Country,
			&w.Hostname, &w.Source, &w.LastError, &w.LookedUpAt, &w.ExpiresAt,
			&w.LookupVersion); err != nil {
			return nil, fmt.Errorf("scan whois range row: %w", err)
		}
		result = append(result, w)
	}
	return result, nil
}

// whoisChunkSize keeps the generated IN (...) list well under the SQLite
// parameter limit on older builds.
const whoisChunkSize = 500

// GetIPWhois returns the cached entries for the given IPs. IPs with nothing
// cached are simply absent from the map.
func (s *Storage) GetIPWhois(ips []string) (map[string]IPWhois, error) {
	result := make(map[string]IPWhois, len(ips))
	for start := 0; start < len(ips); start += whoisChunkSize {
		end := min(start+whoisChunkSize, len(ips))
		chunk := ips[start:end]

		args := make([]interface{}, len(chunk))
		for i, ip := range chunk {
			args[i] = ip
		}
		query := `SELECT ip, org, network, cidr, country, hostname, source, last_error,
			looked_up_at, expires_at, lookup_version
			FROM ip_whois WHERE ip IN (?` + strings.Repeat(", ?", len(chunk)-1) + `)`

		rows, err := s.db.Query(query, args...)
		if err != nil {
			return nil, fmt.Errorf("query ip whois: %w", err)
		}
		for rows.Next() {
			var w IPWhois
			if err := rows.Scan(&w.IP, &w.Org, &w.Network, &w.CIDR, &w.Country,
				&w.Hostname, &w.Source, &w.LastError, &w.LookedUpAt, &w.ExpiresAt,
				&w.LookupVersion); err != nil {
				_ = rows.Close()
				return nil, fmt.Errorf("scan ip whois row: %w", err)
			}
			result[w.IP] = w
		}
		_ = rows.Close()
	}
	return result, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

// DomainStats holds statistics for a single domain
type DomainStats struct {
	Domain            string  `json:"domain"`
	TotalMessages     int     `json:"total_messages"`
	CompliantMessages int     `json:"compliant_messages"`
	ComplianceRate    float64 `json:"compliance_rate"`
}

// OrgStats holds statistics for a reporting organization
type OrgStats struct {
	OrgName string `json:"org_name"`
	Reports int    `json:"reports"`
}

// DispositionStats holds statistics for a disposition type
type DispositionStats struct {
	Disposition string `json:"disposition"`
	Count       int    `json:"count"`
}

// AuthResultStats holds authentication result statistics
type AuthResultStats struct {
	Result string `json:"result"`
	Count  int    `json:"count"`
}

// GetDomainStats returns statistics grouped by domain
func (s *Storage) GetDomainStats() ([]DomainStats, error) {
	rows, err := s.db.Query(`
		SELECT domain,
		       COALESCE(SUM(total_messages), 0) as total_messages,
		       COALESCE(SUM(compliant_messages), 0) as compliant_messages
		FROM reports
		GROUP BY domain
	`)
	if err != nil {
		return nil, fmt.Errorf("query domain stats: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var stats []DomainStats
	for rows.Next() {
		var ds DomainStats
		if err := rows.Scan(&ds.Domain, &ds.TotalMessages, &ds.CompliantMessages); err != nil {
			return nil, fmt.Errorf("scan domain stats row: %w", err)
		}
		if ds.TotalMessages > 0 {
			ds.ComplianceRate = float64(ds.CompliantMessages) / float64(ds.TotalMessages) * 100
		}
		stats = append(stats, ds)
	}
	return stats, nil
}

// GetOrgStats returns statistics grouped by reporting organization
func (s *Storage) GetOrgStats() ([]OrgStats, error) {
	rows, err := s.db.Query(`
		SELECT org_name, COUNT(*) as reports
		FROM reports
		GROUP BY org_name
	`)
	if err != nil {
		return nil, fmt.Errorf("query org stats: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var stats []OrgStats
	for rows.Next() {
		var os OrgStats
		if err := rows.Scan(&os.OrgName, &os.Reports); err != nil {
			return nil, fmt.Errorf("scan org stats row: %w", err)
		}
		stats = append(stats, os)
	}
	return stats, nil
}

// GetDispositionStats returns message counts grouped by disposition
func (s *Storage) GetDispositionStats() ([]DispositionStats, error) {
	rows, err := s.db.Query(`
		SELECT COALESCE(disposition, 'unknown') as disposition,
		       SUM(count) as total_count
		FROM records
		GROUP BY disposition
	`)
	if err != nil {
		return nil, fmt.Errorf("query disposition stats: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var stats []DispositionStats
	for rows.Next() {
		var ds DispositionStats
		if err := rows.Scan(&ds.Disposition, &ds.Count); err != nil {
			return nil, fmt.Errorf("scan disposition stats row: %w", err)
		}
		stats = append(stats, ds)
	}
	return stats, nil
}

// GetSPFStats returns SPF authentication result statistics
func (s *Storage) GetSPFStats() ([]AuthResultStats, error) {
	rows, err := s.db.Query(`
		SELECT COALESCE(spf_result, 'unknown') as result,
		       SUM(count) as total_count
		FROM records
		GROUP BY spf_result
	`)
	if err != nil {
		return nil, fmt.Errorf("query SPF stats: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var stats []AuthResultStats
	for rows.Next() {
		var as AuthResultStats
		if err := rows.Scan(&as.Result, &as.Count); err != nil {
			return nil, fmt.Errorf("scan SPF stats row: %w", err)
		}
		stats = append(stats, as)
	}
	return stats, nil
}

// GetDKIMStats returns DKIM authentication result statistics
func (s *Storage) GetDKIMStats() ([]AuthResultStats, error) {
	rows, err := s.db.Query(`
		SELECT COALESCE(dkim_result, 'unknown') as result,
		       SUM(count) as total_count
		FROM records
		GROUP BY dkim_result
	`)
	if err != nil {
		return nil, fmt.Errorf("query DKIM stats: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var stats []AuthResultStats
	for rows.Next() {
		var as AuthResultStats
		if err := rows.Scan(&as.Result, &as.Count); err != nil {
			return nil, fmt.Errorf("scan DKIM stats row: %w", err)
		}
		stats = append(stats, as)
	}
	return stats, nil
}

// busyTimeout is how long SQLite waits for a lock before giving up. Without
// it, a dashboard request reading while the enricher writes fails outright
// with SQLITE_BUSY; the writes involved take milliseconds, so waiting is
// always the better answer.
const busyTimeout = 5 * time.Second

// withDSNParam appends a driver-specific parameter to a database path.
func withDSNParam(dbPath, param string) string {
	if strings.Contains(dbPath, "?") {
		return dbPath + "&" + param
	}
	return dbPath + "?" + param
}

// WhoisCacheVersion identifies the lookup logic that produced a cached entry.
// Bump it whenever a change would make an older entry wrong or unwanted, and
// every entry written before the change is treated as expired: an upgrade then
// re-resolves in the background instead of serving stale answers for days.
//
// 2: contacts who are natural persons are never named as the owner.
const WhoisCacheVersion = 2

// migrate applies schema changes to databases created by an earlier version.
// CREATE TABLE IF NOT EXISTS leaves an existing table alone, so a new column
// has to be added explicitly.
func (s *Storage) migrate() error {
	// SQLite has no ADD COLUMN IF NOT EXISTS; adding one that is already there
	// is the expected outcome on every start but the first.
	_, err := s.db.Exec(`ALTER TABLE ip_whois ADD COLUMN lookup_version INTEGER NOT NULL DEFAULT 0`)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column") {
		return fmt.Errorf("add ip_whois.lookup_version: %w", err)
	}
	return nil
}
