package main

import (
	"os"

	"github.com/pblaszczyk/sscc/cli"
)

func main() {
	cli.NewApp().Run(os.Args)
}
