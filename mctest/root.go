package mctest

import (
	"os"
	"path/filepath"
)

// FindInParentDirs searches for any of the specified file or directory names in the current and parent directories.
// If any of the given names are found the path of the directory is returned.
// Otherwise an empty string and false are returned.
func FindInParentDirs(start string, anyOf ...string) (string, bool) {
	path := start

	for {
		for _, name := range anyOf {
			candidate := filepath.Join(path, name)
			if _, err := os.Stat(candidate); err == nil {
				return path, true
			}
		}

		parent := filepath.Dir(path)
		if parent == path {
			return "", false
		}
		path = parent
	}
}

var RepositoryMarkers = []string{".git", "go.mod"}

// FindRepositoryRoot searches for the root of the repository.
// It is a shorthand for calling FindInParentDirs with RepositoryMarkers.
func FindRepositoryRoot(start string) (string, bool) {
	return FindInParentDirs(start, RepositoryMarkers...)
}
