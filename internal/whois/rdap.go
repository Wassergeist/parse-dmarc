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

// pickOwner walks the entity tree depth-first and returns the best-ranked
// vCard full name, along with the entity it came from.
func pickOwner(entities []rdapEntity) (string, *rdapEntity) {
	bestRank := len(entityRolePriority)
	var bestName string
	var bestEntity *rdapEntity

	var walk func(list []rdapEntity)
	walk = func(list []rdapEntity) {
		for i := range list {
			e := &list[i]
			if rank := roleRank(e.Roles); rank < bestRank {
				// RIPE assignment objects often name themselves after their own
				// handle ("HOS-GUN"), while the readable company name sits on a
				// nested organisation entity. Such a placeholder tells a reader
				// nothing, so keep looking.
				if name := vcardValue(e.VcardArray, "fn"); name != "" && !strings.EqualFold(name, e.Handle) {
					bestRank, bestName, bestEntity = rank, name, e
				}
			}
			walk(e.Entities)
		}
	}
	walk(entities)

	return bestName, bestEntity
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
