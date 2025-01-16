package jsonnet

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
)

// importer implements jsonnet.Importer interface with path restrictions
type importer struct {
	baseDir string
	jpaths  []string
	fsCache map[string]*fsCacheEntry
}

type fsCacheEntry struct {
	exists   bool
	contents jsonnet.Contents
}

// newImporter creates a new Importer that restricts file access to baseDir and jpaths
func newImporter(baseDir string, jpaths []string) *importer {
	return &importer{
		baseDir: baseDir,
		jpaths:  jpaths,
		fsCache: make(map[string]*fsCacheEntry),
	}
}

func (i *importer) isPathAllowed(path string) bool {
	// Check if path is within base directory or any of the jpaths
	if strings.HasPrefix(path, i.baseDir) {
		return true
	}
	for _, jpath := range i.jpaths {
		if strings.HasPrefix(path, jpath) {
			return true
		}
	}
	return false
}

// tryPath attempts to read a file from the given directory
// Cache is used to store the result of the file read operation
// If the file is not found, the cache is updated with a false value
// If the file is found, the cache is updated with the file contents
func (i *importer) tryPath(dir, importedPath string) (found bool, contents jsonnet.Contents, foundHere string, err error) {
	// Construct absolute path
	var absPath string
	if filepath.IsAbs(importedPath) {
		absPath = importedPath
	} else {
		absPath = filepath.Join(dir, importedPath)
	}
	absPath = filepath.Clean(absPath)

	// Verify path is within allowed directories
	if !i.isPathAllowed(absPath) {
		return false, jsonnet.Contents{}, "", fmt.Errorf("access denied: %s is outside of allowed directories", importedPath)
	}

	// Check cache
	if entry, ok := i.fsCache[absPath]; ok {
		if !entry.exists {
			return false, jsonnet.Contents{}, "", nil
		}
		return entry.exists, entry.contents, absPath, nil
	}

	// Read file
	content, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			i.fsCache[absPath] = &fsCacheEntry{exists: false}
			return false, jsonnet.Contents{}, "", nil
		}
		// Only return actual errors (permissions, etc)
		return false, jsonnet.Contents{}, "", err
	}

	// Cache and return result
	entry := &fsCacheEntry{
		exists:   true,
		contents: jsonnet.MakeContents(string(content)),
	}
	i.fsCache[absPath] = entry
	return true, entry.contents, absPath, nil
}

// Import implements jsonnet.Importer interface
// It searches for files first relative to the importing file
// then in the specified library paths (jpaths) in reverse order
func (i *importer) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	// If importedFrom is empty, use the first jpath as the base directory
	var dir string
	if importedFrom == "" && len(i.jpaths) > 0 {
		dir = i.jpaths[0]
	} else {
		dir = filepath.Dir(importedFrom)
	}

	// First try relative to the importing file
	found, content, foundHere, err := i.tryPath(dir, importedPath)
	if err != nil {
		return jsonnet.Contents{}, "", err
	}
	if found {
		return content, foundHere, nil
	}

	// Then try library paths in reverse order
	for j := len(i.jpaths) - 1; j >= 0; j-- {
		found, content, foundHere, err = i.tryPath(i.jpaths[j], importedPath)
		if err != nil {
			return jsonnet.Contents{}, "", err
		}
		if found {
			return content, foundHere, nil
		}
	}

	// Nothing found - return empty contents without error as per Importer interface
	return jsonnet.Contents{}, "", nil
}
