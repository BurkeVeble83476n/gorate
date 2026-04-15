package cli

import (
	"fmt"
	"os"
)

// writeFile writes data to the given path, creating or truncating the file.
func writeFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}
