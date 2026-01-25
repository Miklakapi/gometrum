package main

import (
	"context"
	"fmt"
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
		printErrorAndExit(err, 2)
	}

	switch {
	case flags.GenerateConfig:
		if flags.ConfigPath == "" {
			fmt.Print(config.ExampleYAML)
			return
		}

		if err = config.SaveExample(flags.ConfigPath); err != nil {
			printErrorAndExit(err, 1)
		}
		return
	case flags.PrintConfig:
		conf, err := config.LoadString(flags.ConfigPath)
		if err != nil {
			printErrorAndExit(err, 1)
		}
		fmt.Print(conf)
		return
	case flags.Validate:
		if err = config.Validate(flags.ConfigPath); err != nil {
			printErrorAndExit(err, 1)
		}
		return
	}

	logger.SetupLogger(flags.LogLevel)

	// Handle normal run
}

func printErrorAndExit(err error, code int) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(code)
}
