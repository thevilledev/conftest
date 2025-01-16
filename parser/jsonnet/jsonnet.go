package jsonnet

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/go-jsonnet"
)

type Parser struct{}

// Unmarshal unmarshals Jsonnet files
func (p *Parser) Unmarshal(data []byte, v interface{}) error {
	vm := jsonnet.MakeVM()

	// Use current directory as the boundary for imports
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	// Add import paths with absolute paths
	jpaths := []string{
		cwd,
	}

	// Create a safe importer with search paths
	vm.Importer(newImporter(cwd, jpaths))

	snippetStream, err := vm.EvaluateAnonymousSnippet("", string(data))
	if err != nil {
		return fmt.Errorf("evaluate anonymous snippet: %w", err)
	}

	if err := json.Unmarshal([]byte(snippetStream), v); err != nil {
		return fmt.Errorf("unmarshal json failed: %w", err)
	}

	return nil
}
