package config

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed example.yaml
var ExampleYAML string

func SaveExample(path string) error {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, []byte(ExampleYAML), 0644)
}

func Load(path string) (string, error) {
	panic("TODO")
}

func Validate(path string) error {
	panic("TODO")
}
