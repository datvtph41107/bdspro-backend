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

	patternSegments := strings.Split(strings.Trim(pattern, "/"), "/")
	pathSegments := strings.Split(strings.Trim(path, "/"), "/")

	for len(pathSegments) > 0 && pathSegments[len(pathSegments)-1] == "" {
		pathSegments = pathSegments[:len(pathSegments)-1]
	}
	for len(patternSegments) > 0 && patternSegments[len(patternSegments)-1] == "" {
		patternSegments = patternSegments[:len(patternSegments)-1]
	}

	if len(patternSegments) > 0 && patternSegments[len(patternSegments)-1] == "*" {
		patternSegments = patternSegments[:len(patternSegments)-1]
		if len(pathSegments) != len(patternSegments)+1 {
			return false
		}
		for i, patternSegment := range patternSegments {
			pathSegment := pathSegments[i]
			if patternSegment == ":" {
				if pathSegment == "" {
					return false
				}
				continue
			}
			if patternSegment != pathSegment {
				return false
			}
		}
		return true
	}

	if len(patternSegments) != len(pathSegments) {
		return false
	}
	for i, patternSegment := range patternSegments {
		pathSegment := pathSegments[i]
		switch patternSegment {
		case "*":
			continue
		case ":":
			if pathSegment == "" {
				return false
			}
		default:
			if patternSegment != pathSegment {
				return false
			}
		}
	}
	return true
}
