package jsonnet

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImporter(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir := t.TempDir()

	// Create test directories
	allowedDir := filepath.Join(tmpDir, "allowed")
	if err := os.Mkdir(allowedDir, os.FileMode(0755)); err != nil {
		t.Fatalf("Failed to create allowed dir: %v", err)
	}

	// Create test files
	files := map[string]string{
		filepath.Join(allowedDir, "valid.libsonnet"):          `{ "key": "value" }`,
		filepath.Join(allowedDir, "subdir", "file.libsonnet"): `{ "sub": "value" }`,
		filepath.Join(tmpDir, "outside.libsonnet"):            `{ "outside": true }`,
	}

	for path, content := range files {
		if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0755)); err != nil {
			t.Fatalf("Failed to create directories for %s: %v", path, err)
		}
		if err := os.WriteFile(path, []byte(content), os.FileMode(0600)); err != nil {
			t.Fatalf("Failed to write file %s: %v", path, err)
		}
	}

	tests := []struct {
		name         string
		importedFrom string
		importedPath string
		jpaths       []string
		wantErr      bool
		errContains  string
		wantContent  string
	}{
		{
			name:         "valid import within allowed directory",
			importedFrom: filepath.Join(allowedDir, "main.jsonnet"),
			importedPath: "./valid.libsonnet",
			jpaths:       []string{allowedDir},
			wantErr:      false,
			wantContent:  `{ "key": "value" }`,
		},
		{
			name:         "attempt to access file outside allowed directory",
			importedFrom: filepath.Join(allowedDir, "main.jsonnet"),
			importedPath: "../outside.libsonnet",
			jpaths:       []string{allowedDir},
			wantErr:      true,
			errContains:  "access denied",
		},
		{
			name:         "import from jpaths",
			importedFrom: "",
			importedPath: "valid.libsonnet",
			jpaths:       []string{allowedDir},
			wantErr:      false,
			wantContent:  `{ "key": "value" }`,
		},
		{
			name:         "file does not exist",
			importedFrom: filepath.Join(allowedDir, "main.jsonnet"),
			importedPath: "./nonexistent.libsonnet",
			jpaths:       []string{allowedDir},
			wantErr:      false,
			wantContent:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importer := newImporter(allowedDir, tt.jpaths)
			contents, foundAt, err := importer.Import(tt.importedFrom, tt.importedPath)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q but got %v", tt.errContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// For nonexistent files, we expect empty contents and empty foundAt
			if tt.wantContent == "" {
				if foundAt != "" {
					t.Error("expected empty foundAt for nonexistent file")
				}
				return
			}

			// For existing files, verify the content
			if contents.String() != tt.wantContent {
				t.Errorf("expected content %q but got %q", tt.wantContent, contents.String())
			}
		})
	}
}
