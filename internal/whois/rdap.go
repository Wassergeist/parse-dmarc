package whois

import (
	"errors"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
)

// rdapIPNetwork is the subset of an RFC 9083 IP network object we display.
// Every field is optional in practice: the five RIRs each publish a different
// mix, so nothing here may be relied upon.
type rdapIPNetwork struct {
	Handle string `json:"handle"`
	Name   string `json:"name"`
	// Country is set by RIPE, APNIC, AFRINIC and LACNIC; ARIN usually omits it.
	Country      string       `json:"country"`
	StartAddress string       `json:"startAddress"`
	EndAddress   string       `json:"endAddress"`
	Cidr0        []rdapCIDR   `json:"cidr0_cidrs"`
	Entities     []rdapEntity `json:"entities"`
}

type rdapCIDR struct {
	V4Prefix string `json:"v4prefix"`
	V6Prefix string `json:"v6prefix"`
	Length   int    `json:"length"`
}

type rdapEntity struct {
	Handle string   `json:"handle"`
	Roles  []string `json:"roles"`
	// VcardArray is jCard: ["vcard", [[name, params, type, value], ...]].
	// It stays raw so that one registry sending something else entirely does
	// not fail the whole response; it is decoded per entity, defensively.
	VcardArray json.RawMessage `json:"vcardArray"`
	Entities   []rdapEntity    `json:"entities"`
}

// errNotFound marks an IP the registry does not have an entry for.
var errNotFound = errors.New("not found in registry")

// entityRolePriority ranks the roles we accept as "the owner", best first.
// ARIN and RIPE publish a registrant; APNIC frequently ships only admin or
// tech contacts, and an abuse desk is still a better answer than nothing.
var entityRolePriority = []string{"registrant", "administrative", "technical", "abuse"}

// abuseRole is treated specially when naming an owner: see trimAbuseSuffix.
const abuseRole = "abuse"

// parseRDAP extracts the displayable fields from an RDAP IP network response.
func parseRDAP(body []byte) (Info, error) {
	var nw rdapIPNetwork
	if err := json.Unmarshal(body, &nw); err != nil {
		return Info{}, errors.New("rdap: response is not a JSON object")
	}

	info := Info{
		Network: strings.TrimSpace(nw.Name),
		CIDR:    formatCIDR(nw),
	}
	if info.Network == "" {
		info.Network = strings.TrimSpace(nw.Handle)
	}

	owner, ownerEntity := pickOwner(nw.Entities)
	info.Org = owner
	if info.Org == "" {
		info.Org = info.Network
	}

	info.Country = normalizeCountry(nw.Country)
	if info.Country == "" && ownerEntity != nil {
		info.Country = countryFromVcard(ownerEntity.VcardArray)
	}

	return info, nil
}

// formatCIDR prefers the cidr0 extension and falls back to the address range.
func formatCIDR(n rdapIPNetwork) string {
	for _, c := range n.Cidr0 {
		prefix := c.V4Prefix
		if prefix == "" {
			prefix = c.V6Prefix
		}
		if prefix != "" {
			return prefix + "/" + strconv.Itoa(c.Length)
		}
	}
	if n.StartAddress != "" && n.EndAddress != "" {
		return n.StartAddress + " - " + n.EndAddress
	}
	return ""
}

// pickOwner walks the entity tree depth-first and returns the best name to
// show as the owner, along with the entity it came from.
//
// Ranking is by vCard kind first and role second: a registry may list a whole
// company under an organisation entity while the administrative contact for
// the same range is a named employee, and the company is what a reader wants.
func pickOwner(entities []rdapEntity) (string, *rdapEntity) {
	best := ownerCandidate{kind: kindRankUnusable, role: len(entityRolePriority)}
	var bestEntity *rdapEntity

	var walk func(list []rdapEntity)
	walk = func(list []rdapEntity) {
		for i := range list {
			e := &list[i]
			if c, ok := ownerCandidateFor(e); ok && c.betterThan(best) {
				best, bestEntity = c, e
			}
			walk(e.Entities)
		}
	}
	walk(entities)

	if bestEntity == nil {
		return "", nil
	}
	return best.name, bestEntity
}

// ownerCandidate is one entity considered as the owner, with its ranking.
type ownerCandidate struct {
	name string
	kind int
	role int
}

// betterThan orders candidates: kind first, then role.
func (c ownerCandidate) betterThan(other ownerCandidate) bool {
	if c.kind != other.kind {
		return c.kind < other.kind
	}
	return c.role < other.role
}

// Ranking of the vCard kind. An organisation is what we are looking for; a
// group (a role account or an abuse desk) names the company too. An
// individual never qualifies: see ownerCandidateFor.
const (
	kindRankOrg = iota
	kindRankGroup
	kindRankUnknown
	kindRankUnusable
)

// ownerCandidateFor judges whether an entity may name the owner.
func ownerCandidateFor(e *rdapEntity) (ownerCandidate, bool) {
	kind := strings.ToLower(vcardValue(e.VcardArray, "kind"))
	if kind == "individual" {
		// This is a named human being, usually the technical contact for the
		// range. Publishing their name on a dashboard would be both wrong
		// (the reader wants the company) and a needless exposure of personal
		// data, so they are never a candidate.
		return ownerCandidate{}, false
	}

	name := vcardValue(e.VcardArray, "fn")
	if name == "" {
		return ownerCandidate{}, false
	}
	// RIPE assignment objects often name themselves after their own handle
	// ("MNT-DOMAINFACTORY"), while the readable name, if any, sits elsewhere.
	// Such a placeholder tells a reader nothing, so keep looking.
	if strings.EqualFold(name, e.Handle) {
		return ownerCandidate{}, false
	}

	role := roleRank(e.Roles)
	if role == len(entityRolePriority) {
		return ownerCandidate{}, false
	}

	// An abuse desk is named "<company> Abuse"; the desk part is noise in a
	// column answering "who is this?". If nothing is left, the caller falls
	// back to the network name.
	if entityRolePriority[role] == abuseRole {
		if name = trimAbuseSuffix(name); name == "" {
			return ownerCandidate{}, false
		}
	}

	return ownerCandidate{name: name, kind: kindRank(kind), role: role}, true
}

func kindRank(kind string) int {
	switch kind {
	case "org":
		return kindRankOrg
	case "group":
		return kindRankGroup
	default:
		// Not every registry sets kind; such an entity is still usable, it
		// just loses to one that says what it is.
		return kindRankUnknown
	}
}

// trimAbuseSuffix drops a trailing "abuse" word from an abuse desk's name.
func trimAbuseSuffix(name string) string {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) < len(abuseRole) {
		return trimmed
	}
	if strings.EqualFold(trimmed[len(trimmed)-len(abuseRole):], abuseRole) {
		return strings.TrimSpace(trimmed[:len(trimmed)-len(abuseRole)])
	}
	return trimmed
}

func roleRank(roles []string) int {
	best := len(entityRolePriority)
	for _, role := range roles {
		for i, want := range entityRolePriority {
			if strings.EqualFold(role, want) && i < best {
				best = i
			}
		}
	}
	return best
}

// vcardProperties decodes the jCard property list, tolerating every shape a
// registry might send in place of it.
func vcardProperties(vcard json.RawMessage) [][]json.RawMessage {
	if len(vcard) == 0 {
		return nil
	}
	var outer []json.RawMessage
	if err := json.Unmarshal(vcard, &outer); err != nil || len(outer) < 2 {
		return nil
	}
	var props [][]json.RawMessage
	if err := json.Unmarshal(outer[1], &props); err != nil {
		return nil
	}
	return props
}

// vcardValue returns the string value of the first property called name.
func vcardValue(vcard json.RawMessage, name string) string {
	for _, prop := range vcardProperties(vcard) {
		if len(prop) < 4 {
			continue
		}
		var propName string
		if err := json.Unmarshal(prop[0], &propName); err != nil || !strings.EqualFold(propName, name) {
			continue
		}
		var value string
		if err := json.Unmarshal(prop[3], &value); err != nil {
			continue
		}
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

// countryFromVcard reads the country element of a structured jCard address.
// ARIN instead ships the address as free text in the property parameters,
// which is not an ISO code, so those responses simply yield no country.
func countryFromVcard(vcard json.RawMessage) string {
	for _, prop := range vcardProperties(vcard) {
		if len(prop) < 4 {
			continue
		}
		var propName string
		if err := json.Unmarshal(prop[0], &propName); err != nil || !strings.EqualFold(propName, "adr") {
			continue
		}
		var parts []string
		if err := json.Unmarshal(prop[3], &parts); err != nil {
			continue
		}
		// RFC 6350 structured address: the country is the seventh element.
		const countryIndex = 6
		if len(parts) > countryIndex {
			if c := normalizeCountry(parts[countryIndex]); c != "" {
				return c
			}
		}
	}
	return ""
}

// normalizeCountry accepts only a two-letter ISO 3166-1 alpha-2 code; a full
// country name is left out rather than guessed at.
func normalizeCountry(s string) string {
	s = strings.ToUpper(strings.TrimSpace(s))
	if len(s) != 2 {
		return ""
	}
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return ""
		}
	}
	return s
}
