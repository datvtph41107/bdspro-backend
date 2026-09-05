package filehttp

// storedPathResolver is the narrow storage capability required by HTTP readers
// to translate stable File wire paths into deployment-specific filesystem paths.
type storedPathResolver interface {
	ResolveStoredPath(storedPath string) (string, error)
}
