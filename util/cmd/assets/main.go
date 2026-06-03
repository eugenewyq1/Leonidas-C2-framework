package main

import (
	"fmt"
	"os"

	"github.com/leonidas-c2/leonidas/util/assets"
)

func main() {
	if err := assets.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
