package main

import (
	"fmt"

	"github.com/Miklakapi/gometrum/internal/cli"
)

func main() {
	flags := cli.ParseFlags()
	fmt.Printf("%+v\n", flags)
}
