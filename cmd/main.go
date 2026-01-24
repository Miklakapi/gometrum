package main

import (
	"fmt"

	"github.com/Miklakapi/gometrum/internal/cli"
	"github.com/Miklakapi/gometrum/internal/config"
)

func main() {
	flags := cli.ParseFlags()
	fmt.Printf("%+v\n", flags)
	fmt.Println(config.ExampleYAML)
}
