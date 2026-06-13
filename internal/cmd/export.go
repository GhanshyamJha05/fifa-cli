package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

func exportJSON(path string, fetch func() (any, error)) error {
	data, err := fetch()
	if err != nil {
		return err
	}
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	fmt.Printf("Exported to %s\n", path)
	return nil
}
