package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Miklakapi/gometrum/internal/cli"
	"github.com/Miklakapi/gometrum/internal/config"
	"github.com/Miklakapi/gometrum/internal/logger"
)

func main() {
	_, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	flags, err := cli.ParseFlags()
	if err != nil {
		log.Fatalln(err)
	}

	logger.SetupLogger(flags.LogLevel)

	switch {
	case flags.GenerateConfig:
		// generate config + exit
		// available flags: configPath + logLevel
		if flags.ConfigPath == "" {
			fmt.Print(config.ExampleYAML)
			return
		}
		// TODO: write to file flags.ConfigPath
		fmt.Print(config.ExampleYAML)
		return
	case flags.PrintConfig:
		// print config + exit
		// available flags: configPath + logLevel
		return
	case flags.Validate:
		// validate config + exit
		// available flags: configPath + logLevel
		return
	}

	// Handle dry run + one + normal run
}
