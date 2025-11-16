package main

import (
	"fmt"
	"os"

	"cut/internal/core"
)

func main() {
	c, err := core.New(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := c.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
