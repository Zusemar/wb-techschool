package main

import (
	"fmt"
	"io"
	"os"

	"minishell/internal"
)

func main() {
	r := internal.NewReader(os.Stdin)

	for {
		fmt.Print("mysh> ")

		line, err := r.Next()
		if err == io.EOF {
			fmt.Println()
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "read error:", err)
			continue
		}
		if line == "" {
			continue
		}

		toks, err := internal.Tokenize(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, "tokenize error:", err)
			continue
		}

		ast, err := internal.Parse(toks)
		if err != nil {
			fmt.Fprintln(os.Stderr, "parse error:", err)
			continue
		}

		exitCode := internal.Execute(ast, os.Stdin, os.Stdout, os.Stderr)
		_ = exitCode
	}
}
