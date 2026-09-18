package whois

import (
	"net"
	"strings"
	"sync"
)

// rangeEntry is one registry answer together with the address range it covers.
type rangeEntry struct {
	first net.IP
	last  net.IP
	info  Info
}

// rangeIndex remembers which range each answer covered, so the many addresses
// a single provider sends from cost one registry request rather than one each.
// A DMARC report routinely lists a dozen relays out of the same assignment.
type rangeIndex struct {
	mu      sync.RWMutex
	entries []rangeEntry
}

// Add records an answer under the range it came with. Answers without a range
// are dropped: they cannot be reused for anything but their own address.
func (r *rangeIndex) Add(info Info) {
	first, last, ok := parseRange(info.CIDR)
	if !ok {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, rangeEntry{first: first, last: last, info: info})
}

// Lookup returns the answer covering ip, if one is known.
func (r *rangeIndex) Lookup(ip net.IP) (Info, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Later entries are more specific in practice: a covering allocation is
	// usually seen before the assignment inside it.
	for i := len(r.entries) - 1; i >= 0; i-- {
		e := r.entries[i]
		if withinRange(ip, e.first, e.last) {
			return e.info, true
		}
	}
	return Info{}, false
}

func (r *rangeIndex) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}

// parseRange accepts both shapes a registry answer carries: a CIDR prefix, and
// the "first - last" form RIPE and ARIN use for an inetnum.
func parseRange(s string) (net.IP, net.IP, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil, false
	}

	if strings.Contains(s, "/") {
		_, network, err := net.ParseCIDR(s)
		if err != nil {
			return nil, nil, false
		}
		return network.IP, lastAddress(network), true
	}

	before, after, found := strings.Cut(s, "-")
	if !found {
		return nil, nil, false
	}
	first := net.ParseIP(strings.TrimSpace(before))
	last := net.ParseIP(strings.TrimSpace(after))
	if first == nil || last == nil {
		return nil, nil, false
	}
	return first, last, true
}

// lastAddress returns the broadcast-equivalent address of a network.
func lastAddress(network *net.IPNet) net.IP {
	last := make(net.IP, len(network.IP))
	copy(last, network.IP)
	for i := range last {
		last[i] |= ^network.Mask[i]
	}
	return last
}

// withinRange reports whether ip falls between first and last inclusive.
func withinRange(ip, first, last net.IP) bool {
	ip, first, last = normalizeIP(ip), normalizeIP(first), normalizeIP(last)
	if ip == nil || first == nil || last == nil || len(ip) != len(first) || len(ip) != len(last) {
		return false
	}
	return bytesCompare(ip, first) >= 0 && bytesCompare(ip, last) <= 0
}

// normalizeIP reduces an address to its canonical 4- or 16-byte form, so an
// IPv4 address written as IPv4-in-IPv6 compares equal to the plain one.
func normalizeIP(ip net.IP) net.IP {
	if v4 := ip.To4(); v4 != nil {
		return v4
	}
	return ip.To16()
}

func bytesCompare(a, b net.IP) int {
	for i := range a {
		switch {
		case a[i] < b[i]:
			return -1
		case a[i] > b[i]:
			return 1
		}
	}
	return 0
}
