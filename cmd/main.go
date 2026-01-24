package main

import (
	"fmt"
	"log"

	"github.com/Miklakapi/gometrum/internal/cli"
	"github.com/Miklakapi/gometrum/internal/config"
)

func main() {
	flags, err := cli.ParseFlags()
	if err != nil {
		log.Fatalln(err)
	}

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

	fmt.Printf("%+v\n", flags)
	fmt.Println(config.ExampleYAML)
}
