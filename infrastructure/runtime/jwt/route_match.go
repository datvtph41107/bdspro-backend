package _jwt

import "strings"

/**
 * isPublicRoute matches one path against the public route contract.
 */
func isPublicRoute(pubRoutes []string, path string) bool {
	for _, route := range pubRoutes {
		if matchRoute(route, path) {
			return true
		}
	}
	return false
}

/**
 * matchRoute supports exact segments, `:` parameters and `*` wildcards.
 */
func matchRoute(pattern, path string) bool {
	if !strings.Contains(pattern, "*") && !strings.Contains(pattern, ":") {
		if pattern == "/" {
			return path == "/"
		}
		pattern = strings.TrimSuffix(pattern, "/")
		path = strings.TrimSuffix(path, "/")
		return path == pattern || strings.HasPrefix(path, pattern+"/")
	}

	patternSegments := splitRouteSegments(pattern)
	pathSegments := splitRouteSegments(path)
	for index, patternSegment := range patternSegments {
		if isWildcardSegment(patternSegment) {
			return index == len(patternSegments)-1 && len(pathSegments) >= index
		}
		if index >= len(pathSegments) {
			return false
		}
		if isParameterSegment(patternSegment) {
			if pathSegments[index] == "" {
				return false
			}
			continue
		}
		if patternSegment != pathSegments[index] {
			return false
		}
	}
	return len(pathSegments) == len(patternSegments)
}

func splitRouteSegments(value string) []string {
	value = strings.Trim(value, "/")
	if value == "" {
		return nil
	}
	segments := strings.Split(value, "/")
	result := segments[:0]
	for _, segment := range segments {
		if segment != "" {
			result = append(result, segment)
		}
	}
	return result
}

func isParameterSegment(segment string) bool {
	return segment == ":" || (strings.HasPrefix(segment, ":") && len(segment) > 1)
}

func isWildcardSegment(segment string) bool {
	return segment == "*" || (strings.HasPrefix(segment, "*") && len(segment) > 1)
}
