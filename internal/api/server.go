package api

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog"

	"github.com/meysam81/parse-dmarc/internal/metrics"
	"github.com/meysam81/parse-dmarc/internal/parser"
	"github.com/meysam81/parse-dmarc/internal/storage"
	"github.com/meysam81/parse-dmarc/internal/whois"
)

//go:embed dist
var distFS embed.FS

// Server represents the API server
type Server struct {
	storage *storage.Storage
	metrics *metrics.Metrics
	log     *zerolog.Logger
	addr    string
	// enricher is nil when whois enrichment is disabled.
	enricher *whois.Enricher
}

// SetEnricher attaches the background whois enricher. Requests only ever
// hand it work to do later; they never wait for a lookup.
func (s *Server) SetEnricher(e *whois.Enricher) {
	s.enricher = e
}

// NewServer creates a new API server
func NewServer(store *storage.Storage, host string, port int, m *metrics.Metrics, log *zerolog.Logger) *Server {
	return &Server{
		storage: store,
		metrics: m,
		log:     log,
		addr:    fmt.Sprintf("%s:%d", host, port),
	}
}

// Start starts the HTTP server
func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/reports", s.handleReports)
	mux.HandleFunc("/api/reports/", s.handleReportDetail)
	mux.HandleFunc("/api/statistics", s.handleStatistics)
	mux.HandleFunc("/api/top-sources", s.handleTopSources)

	// Prometheus metrics endpoint
	if s.metrics != nil {
		mux.Handle("/metrics", s.metrics.Handler())
	}

	// Serve frontend
	// Try to serve embedded files, fallback to nothing if not embedded
	distFiles, err := fs.Sub(distFS, "dist")
	if err == nil {
		mux.Handle("/", http.FileServer(http.FS(distFiles)))
	} else {
		// If dist folder is not embedded, serve a simple message
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				w.Header().Set("Content-Type", "text/html")
				_, _ = fmt.Fprintf(w, `
					<!DOCTYPE html>
					<html>
					<head><title>DMARC Dashboard</title></head>
					<body>
						<h1>DMARC Report Dashboard API</h1>
						<p>API is running. Frontend assets not embedded yet.</p>
						<ul>
							<li><a href="/api/statistics">Statistics</a></li>
							<li><a href="/api/reports">Reports</a></li>
							<li><a href="/api/top-sources">Top Sources</a></li>
							<li><a href="/metrics">Prometheus Metrics</a></li>
						</ul>
					</body>
					</html>
				`)
			} else {
				http.NotFound(w, r)
			}
		})
	}

	// Build handler chain: CORS -> Metrics -> Routes
	var handler http.Handler = mux
	if s.metrics != nil {
		handler = s.metrics.HTTPMiddleware(handler)
	}
	handler = s.corsMiddleware(handler)

	server := &http.Server{
		Addr:    s.addr,
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		s.log.Info().Msg("shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.log.Error().Err(err).Msg("server shutdown error")
		}
	}()

	s.log.Info().Str("addr", s.addr).Msg("starting server")
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server listen on %s: %w", s.addr, err)
	}
	return nil
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleReports returns a list of reports
func (s *Server) handleReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse pagination parameters
	limit := 50
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	reports, err := s.storage.GetReports(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, reports)
}

// handleReportDetail returns a single report detail
func (s *Server) handleReportDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL
	idStr := r.URL.Path[len("/api/reports/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid report ID", http.StatusBadRequest)
		return
	}

	report, err := s.storage.GetReportByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	s.writeJSON(w, s.withSourceOwners(report))
}

// reportDetail is a report plus the ownership data for the addresses in it.
// The report itself stays exactly as parsed: it mirrors the XML the receiver
// sent, and enrichment is ours, not theirs.
type reportDetail struct {
	*parser.Feedback
	// Whois maps a source address to what is known about who owns it.
	Whois map[string]storage.SourceWhois `json:"whois,omitempty"`
	// WhoisPending lists addresses whose lookup is still outstanding.
	WhoisPending []string `json:"whois_pending,omitempty"`
}

// withSourceOwners attaches cached ownership data for every address in the
// report and queues whatever is missing, the same way the dashboard list does.
func (s *Server) withSourceOwners(report *parser.Feedback) any {
	if s.enricher == nil || report == nil {
		return report
	}

	ips := make([]string, 0, len(report.Records))
	seen := make(map[string]struct{}, len(report.Records))
	for i := range report.Records {
		ip := report.Records[i].Row.SourceIP
		if ip == "" {
			continue
		}
		if _, dup := seen[ip]; dup {
			continue
		}
		seen[ip] = struct{}{}
		ips = append(ips, ip)
	}
	if len(ips) == 0 {
		return report
	}

	cached, err := s.storage.GetIPWhois(ips)
	if err != nil {
		// Ownership is an addition to the report, never a reason to fail it.
		s.log.Debug().Err(err).Msg("could not load whois data for report detail")
		return report
	}

	detail := reportDetail{Feedback: report, Whois: make(map[string]storage.SourceWhois, len(cached))}
	now := time.Now().Unix()
	var stale []string
	for _, ip := range ips {
		entry, found := cached[ip]
		if !found || entry.Stale(now) {
			stale = append(stale, ip)
		}
		if !found {
			continue
		}
		if shown := entry.Display(); shown != nil {
			detail.Whois[ip] = *shown
		}
	}

	if len(stale) > 0 {
		s.enricher.Enqueue(stale...)
		detail.WhoisPending = stale
	}
	return detail
}

// handleStatistics returns dashboard statistics
func (s *Server) handleStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := s.storage.GetStatistics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, stats)
}

// handleTopSources returns top source IPs
func (s *Server) handleTopSources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	sources, err := s.storage.GetTopSourceIPs(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Whatever is already cached goes out now; anything missing or stale is
	// looked up in the background and shows up on a later load. Those rows are
	// flagged so the dashboard can say they are being resolved rather than
	// leaving a gap that looks like an answer.
	if s.enricher != nil {
		s.enricher.EnqueueStale(sources)
		now := time.Now().Unix()
		for i := range sources {
			sources[i].WhoisPending = sources[i].WhoisStale(now)
		}
	}

	s.writeJSON(w, sources)
}

// writeJSON writes JSON response
func (s *Server) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.log.Error().Err(err).Msg("failed to encode JSON")
	}
}

// RefreshMetrics updates all Prometheus metrics from current database state
func (s *Server) RefreshMetrics() {
	if s.metrics == nil {
		return
	}

	// Update basic statistics
	stats, err := s.storage.GetStatistics()
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get statistics for metrics")
	} else {
		s.metrics.UpdateStatistics(
			stats.TotalReports,
			stats.TotalMessages,
			stats.CompliantMessages,
			stats.UniqueSourceIPs,
			stats.UniqueDomains,
			stats.ComplianceRate,
		)
	}

	// Update per-domain metrics
	domainStats, err := s.storage.GetDomainStats()
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get domain stats for metrics")
	} else {
		for _, ds := range domainStats {
			s.metrics.UpdateDomainMetrics(ds.Domain, ds.TotalMessages, ds.ComplianceRate)
		}
	}

	// Update per-organization metrics
	orgStats, err := s.storage.GetOrgStats()
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get org stats for metrics")
	} else {
		for _, os := range orgStats {
			s.metrics.UpdateOrgMetrics(os.OrgName, os.Reports)
		}
	}

	// Update disposition metrics
	dispStats, err := s.storage.GetDispositionStats()
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get disposition stats for metrics")
	} else {
		for _, ds := range dispStats {
			s.metrics.UpdateDispositionMetrics(ds.Disposition, ds.Count)
		}
	}

	// Update authentication results
	spfStats, errSpf := s.storage.GetSPFStats()
	dkimStats, errDkim := s.storage.GetDKIMStats()
	if errSpf != nil {
		s.log.Error().Err(errSpf).Msg("failed to get SPF stats for metrics")
	}
	if errDkim != nil {
		s.log.Error().Err(errDkim).Msg("failed to get DKIM stats for metrics")
	}
	if errSpf == nil && errDkim == nil {
		spfResults := make(map[string]int)
		for _, s := range spfStats {
			spfResults[s.Result] = s.Count
		}
		dkimResults := make(map[string]int)
		for _, d := range dkimStats {
			dkimResults[d.Result] = d.Count
		}
		s.metrics.UpdateAuthResults(spfResults, dkimResults)
	}
}

// GetMetrics returns the metrics instance
func (s *Server) GetMetrics() *metrics.Metrics {
	return s.metrics
}
