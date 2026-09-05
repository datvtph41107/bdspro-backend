package enums

import "strings"

// ResolvePoiCategoryIconPath builds full icon URL from stored icon value + domain_file.
func ResolvePoiCategoryIconPath(domainFile, icon string) string {
	icon = strings.TrimSpace(icon)
	if icon == "" {
		return ""
	}
	if strings.HasPrefix(icon, "http://") || strings.HasPrefix(icon, "https://") {
		return icon
	}

	relPath := ""
	if strings.Contains(icon, "/") {
		relPath = strings.TrimLeft(icon, "/")
	} else {
		relPath = "icons/poi-category/" + icon + ".svg"
	}

	domainFile = strings.TrimRight(domainFile, "/")
	if domainFile == "" {
		return relPath
	}
	return domainFile + "/" + relPath
}
