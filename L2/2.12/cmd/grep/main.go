package main

import (
	"fmt"
	"os"

	"../internal/parser"
)

func main() {

	opts, err := parser.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while parsing args %s", err)
	}
}
