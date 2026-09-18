package whois

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
)

const (
	// defaultWhoisBootstrap is queried first to learn which registry is
	// responsible; its answer carries a refer: line pointing at the RIR.
	defaultWhoisBootstrap = "whois.iana.org:43"
	// maxReferrals bounds the referral chain. One hop covers IANA to the RIR;
	// the second covers a range one RIR has transferred to another.
	maxReferrals = 2
	// maxWhoisBytes caps a plaintext response from a server we do not control.
	maxWhoisBytes = 256 << 10
)

// whoisOrgKeys are the field names registries use for the owning organization,
// best first. There is no schema here: every registry named this differently
// long before anyone tried to parse it centrally.
var whoisOrgKeys = []string{"orgname", "organization", "org-name", "descr", "owner", "responsible"}

// whois43 queries the classic WHOIS protocol as a fallback for an IP that RDAP
// could not answer for. Only the few fields the dashboard shows are read; the
// contact names, phone numbers and postal addresses these responses also carry
// are deliberately left behind.
func (c *Client) whois43(ctx context.Context, ip string) (Info, error) {
	server := c.whoisBootstrap
	var lastErr error

	for range maxReferrals + 1 {
		body, err := c.queryWhoisServer(ctx, server, ip)
		if err != nil {
			return Info{}, err
		}

		fields := parseWhoisFields(body)
		if next := whoisReferral(fields); next != "" && !strings.EqualFold(next, server) {
			server = next
			continue
		}

		info := whoisInfo(fields)
		if info.Org == "" && info.Network == "" {
			lastErr = fmt.Errorf("whois: %s returned no usable fields", server)
			break
		}
		return info, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("whois: no registry answered for %s", ip)
	}
	return Info{}, lastErr
}

// queryWhoisServer sends one query and reads the plaintext answer.
func (c *Client) queryWhoisServer(ctx context.Context, server, query string) (string, error) {
	conn, err := c.dial(ctx, "tcp", server)
	if err != nil {
		return "", fmt.Errorf("whois dial %s: %w", server, err)
	}
	defer func() { _ = conn.Close() }()

	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}

	if _, err := conn.Write([]byte(query + "\r\n")); err != nil {
		return "", fmt.Errorf("whois write %s: %w", server, err)
	}

	body, err := io.ReadAll(io.LimitReader(conn, maxWhoisBytes))
	if err != nil {
		return "", fmt.Errorf("whois read %s: %w", server, err)
	}
	return string(body), nil
}

// whoisField is one key/value line of a response, in the order it appeared.
type whoisField struct {
	key   string
	value string
}

// parseWhoisFields splits the response into key/value pairs, dropping the
// comment and banner lines every registry wraps its answers in.
func parseWhoisFields(body string) []whoisField {
	var fields []whoisField

	scanner := bufio.NewScanner(strings.NewReader(body))
	scanner.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "%") || strings.HasPrefix(line, "#") ||
			strings.HasPrefix(line, ">>>") {
			continue
		}
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		fields = append(fields, whoisField{
			key:   strings.ToLower(strings.TrimSpace(key)),
			value: value,
		})
	}

	return fields
}

// firstValue returns the value of the first field matching any of keys. The
// first match wins because registries answer most-specific object first.
func firstValue(fields []whoisField, keys ...string) string {
	for _, key := range keys {
		for _, f := range fields {
			if f.key == key {
				return f.value
			}
		}
	}
	return ""
}

// whoisReferral finds the server this response points at, if any. ARIN spells
// it ReferralServer with a whois:// URL; IANA and the RIRs use refer.
func whoisReferral(fields []whoisField) string {
	if v := firstValue(fields, "refer"); v != "" {
		return withWhoisPort(v)
	}
	if v := firstValue(fields, "referralserver"); v != "" {
		v = strings.TrimPrefix(v, "whois://")
		v = strings.TrimPrefix(v, "rwhois://")
		return withWhoisPort(strings.TrimSuffix(v, "/"))
	}
	return ""
}

func withWhoisPort(host string) string {
	if _, _, err := net.SplitHostPort(host); err == nil {
		return host
	}
	return net.JoinHostPort(host, "43")
}

// whoisInfo maps the fields onto the same shape RDAP produces.
func whoisInfo(fields []whoisField) Info {
	info := Info{
		Network: firstValue(fields, "netname"),
		Country: normalizeCountry(firstValue(fields, "country")),
		CIDR:    firstValue(fields, "cidr", "inetnum", "inet6num", "netrange"),
	}

	info.Org = trimWhoisHandle(firstValue(fields, whoisOrgKeys...))
	if info.Org == "" {
		info.Org = info.Network
	}

	return info
}

// trimWhoisHandle drops the registry handle ARIN appends to an organization
// name, turning "Google LLC (GOGL)" into "Google LLC".
func trimWhoisHandle(org string) string {
	open := strings.LastIndex(org, " (")
	if open <= 0 || !strings.HasSuffix(org, ")") {
		return org
	}
	return strings.TrimSpace(org[:open])
}
