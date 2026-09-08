package httpauth

import "strings"

type RouteKind uint8

const (
	RouteExact RouteKind = iota + 1
	RoutePrefix
	RouteTemplate
	RouteCatchAll
)

// Route makes HTTP auth matching semantics explicit. A path that looks exact
// never gains prefix privileges unless the policy declares RoutePrefix.
type Route struct {
	Kind   RouteKind
	Path   string
	Method string
}

func Exact(path string) Route    { return Route{Kind: RouteExact, Path: cleanRoute(path)} }
func Prefix(path string) Route   { return Route{Kind: RoutePrefix, Path: cleanRoute(path)} }
func Template(path string) Route { return Route{Kind: RouteTemplate, Path: cleanRoute(path)} }
func CatchAll(path string) Route { return Route{Kind: RouteCatchAll, Path: cleanRoute(path)} }

func ExactMethod(method, path string) Route {
	return Route{Kind: RouteExact, Path: cleanRoute(path), Method: normalizeMethod(method)}
}

func PrefixMethod(method, path string) Route {
	return Route{Kind: RoutePrefix, Path: cleanRoute(path), Method: normalizeMethod(method)}
}

func MatchRequest(route Route, method, path string) bool {
	if route.Method != "" && route.Method != normalizeMethod(method) {
		return false
	}
	return MatchRoute(route, path)
}

func MatchRoute(route Route, path string) bool {
	path = cleanRoute(path)
	switch route.Kind {
	case RouteExact:
		return path == route.Path
	case RoutePrefix:
		return path == route.Path || strings.HasPrefix(path, route.Path+"/")
	case RouteTemplate:
		return matchTemplate(route.Path, path, false)
	case RouteCatchAll:
		return matchTemplate(route.Path, path, true)
	default:
		return false
	}
}

func cleanRoute(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "/"
	}
	if value != "/" {
		value = strings.TrimSuffix(value, "/")
	}
	return value
}

func matchTemplate(pattern, path string, allowCatchAll bool) bool {
	patterns, paths := split(pattern), split(path)
	for index, segment := range patterns {
		if index >= len(paths) {
			return false
		}
		if allowCatchAll && (strings.HasPrefix(segment, "*") || strings.Contains(segment, "=**")) {
			return index == len(patterns)-1
		}
		if isParameter(segment) {
			continue
		}
		if segment != paths[index] {
			return false
		}
	}
	return len(patterns) == len(paths)
}

func split(value string) []string {
	value = strings.Trim(value, "/")
	if value == "" {
		return nil
	}
	return strings.Split(value, "/")
}

func isParameter(segment string) bool {
	return (strings.HasPrefix(segment, ":") && len(segment) > 1) ||
		(strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") && !strings.Contains(segment, "=**"))
}

func normalizeMethod(method string) string {
	return strings.ToUpper(strings.TrimSpace(method))
}
