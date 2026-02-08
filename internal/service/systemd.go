package service

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed example.service
var ExampleService string

func SaveExample(path string) error {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, []byte(ExampleService), 0644)
}
