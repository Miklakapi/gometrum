package config

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed example.yaml
var ExampleYAML string

type Config struct {
	host string
}

func SaveExample(path string) error {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, []byte(ExampleYAML), 0644)
}

func load(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func LoadString(path string) (string, error) {
	data, err := load(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func LoadConfig(path string) (Config, error) {
	panic("TODO")
}

func Validate(path string) error {
	panic("TODO")
}
