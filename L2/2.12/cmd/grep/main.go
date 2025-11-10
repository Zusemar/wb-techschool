package main

import (
	"fmt"
	"os"

	"grep/internal/engine"
	"grep/internal/matcher"
	"grep/internal/output"
	"grep/internal/parser"
	"grep/internal/reader"
)

func main() {
	opts, err := parser.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	// choose input source
	var src reader.Reader
	if opts.Filename == "" {
		src = reader.NewFromStdin()
	} else {
		src, err = reader.NewFromFile(opts.Filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot open file:", err)
			os.Exit(1)
		}
	}
	defer src.Close()

	m, err := matcher.New(opts.Pattern, opts.Fixed, opts.IgnoreCase, opts.Invert)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid pattern:", err)
		os.Exit(1)
	}

	e := engine.New(m, opts)
	f := output.NewFormatter(os.Stdout, opts)

	if err := e.Run(src, f); err != nil {
		fmt.Fprintln(os.Stderr, "grep failed:", err)
		os.Exit(1)
	}
}
