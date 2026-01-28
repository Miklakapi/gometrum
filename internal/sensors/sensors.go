package sensors

import (
	"github.com/Miklakapi/gometrum/internal/config"
)

type SensorDefinition struct {
	DefaultName string
	DefaultIcon string
	DefaultUnit string
}

func Prepare(cfg *config.Config) error {
	if err := Normalize(cfg); err != nil {
		return err
	}
	if err := Validate(*cfg); err != nil {
		return err
	}
	return nil
}

func Normalize(cfg *config.Config) error {
	panic("TODO")
}

func Validate(cfg config.Config) error {
	panic("TODO")
}
