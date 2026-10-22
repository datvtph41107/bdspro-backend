package controller

import "strings"

// normalizeClientLocationCatalogJSON shapes /v2/tqd/client/provinces|wards payloads:
// expose display "name", drop legacy "shortName".
func normalizeClientLocationCatalogJSON(data []any) []any {
	if len(data) == 0 {
		return data
	}
	out := make([]any, len(data))
	for i, item := range data {
		out[i] = normalizeClientLocationCatalogItem(item)
	}
	return out
}

func normalizeClientLocationCatalogItem(item any) any {
	m, ok := item.(map[string]any)
	if !ok {
		return item
	}
	name := firstNonEmptyJSONString(m, "fullName", "officialName", "shortName")
	delete(m, "shortName")
	if name != "" {
		m["name"] = name
	}
	return m
}

func firstNonEmptyJSONString(m map[string]any, keys ...string) string {
	for _, key := range keys {
		if v, ok := m[key].(string); ok {
			if s := strings.TrimSpace(v); s != "" {
				return s
			}
		}
	}
	return ""
}
