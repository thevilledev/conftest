package jsonnet

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestJsonnetParser(t *testing.T) {
	parser := &Parser{}

	// Example from Jsonnet(https://jsonnet.org/)
	sample := `// Edit me!
{
  person1: {
    name: "Alice",
    welcome: "Hello " + self.name + "!",
  },
  person2: self.person1 { name: "Bob" },
}`

	var input interface{}
	if err := parser.Unmarshal([]byte(sample), &input); err != nil {
		t.Fatalf("parser should not have thrown an error: %v", err)
	}

	if input == nil {
		t.Error("there should be information parsed but its nil")
	}

	item := input.(map[string]interface{})

	if len(item) == 0 {
		t.Error("there should be at least one item defined in the parsed file, but none found")
	}
}

func TestUnmarshalWithImports(t *testing.T) {
	tests := []struct {
		name       string
		files      map[string]string // map of relative path to content
		input      string            // the main file to evaluate
		wantData   map[string]interface{}
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "basic import and importstr",
			files: map[string]string{
				"main.jsonnet": `{
					local data = import './data.libsonnet',
					local content = importstr 'content.txt',
					data: data,
					content: content,
				}`,
				"data.libsonnet": `{
					hello: "world",
					nested: import 'nested/more.libsonnet'
				}`,
				"nested/more.libsonnet": `{
					more: "data"
				}`,
				"content.txt": "Hello from content.txt",
			},
			input: "main.jsonnet",
			wantData: map[string]interface{}{
				"data": map[string]interface{}{
					"hello": "world",
					"nested": map[string]interface{}{
						"more": "data",
					},
				},
				"content": "Hello from content.txt",
			},
			wantErr: false,
		},
		{
			name: "attempt to read system file",
			files: map[string]string{
				"malicious.jsonnet": `{
					local secret = importstr '/foo/bar/baz',
					data: secret
				}`,
			},
			input:      "malicious.jsonnet",
			wantData:   nil,
			wantErr:    true,
			wantErrMsg: "access denied: /foo/bar/baz is outside of allowed directories",
		},
		{
			name: "attempt to escape via relative path",
			files: map[string]string{
				"relative_escape.jsonnet": `{
					local data = import '../outside.libsonnet',
					result: data
				}`,
				"../outside/outside.libsonnet": `{ "secret": "data" }`,
			},
			input:      "relative_escape.jsonnet",
			wantData:   nil,
			wantErr:    true,
			wantErrMsg: "access denied: ../outside.libsonnet is outside of allowed directories",
		},
	}

	parser := &Parser{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory structure for testing
			tmpDir := t.TempDir()

			// Create test files for this test case
			for path, content := range tt.files {
				fullPath := filepath.Join(tmpDir, path)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
					t.Fatalf("Failed to create directories for %s: %v", path, err)
				}
				if err := os.WriteFile(fullPath, []byte(content), 0600); err != nil {
					t.Fatalf("Failed to write file %s: %v", path, err)
				}
			}

			// Change to the temp directory for testing
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("Failed to change to temp directory: %v", err)
			}

			// Run the test
			var got interface{}
			err := parser.Unmarshal([]byte(tt.files[tt.input]), &got)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("error message mismatch\ngot:  %v\nwant: %v", err, tt.wantErrMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			gotMap := got.(map[string]interface{})
			if !reflect.DeepEqual(gotMap, tt.wantData) {
				t.Errorf("Unmarshal() = %v, want %v", gotMap, tt.wantData)
			}
		})
	}
}
