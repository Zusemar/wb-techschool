package parser

import (
	"fmt"
	"strconv"
)

// Options holds all flags and its values
type Options struct {
	Pattern     string
	Filename    string
	After       int
	Before      int
	CountOnly   bool
	IgnoreCase  bool
	Invert      bool
	Fixed       bool
	WithLineNum bool
}

// Parse is options machine that reads all flags for os.Args to Options struct
// and returns an error if something went wrong
func Parse(argv []string) (Options, error) {
	var opts Options
	for i := 0; i < len(argv); i++ {
		switch argv[i] {
		case "-A":
			n, ni, err := parseIntFlag(argv, i)
			if err != nil {
				return opts, err
			}
			opts.After, i = n, ni
		case "-B":
			n, ni, err := parseIntFlag(argv, i)
			if err != nil {
				return opts, err
			}
			opts.Before, i = n, ni
		case "-C":
			n, ni, err := parseIntFlag(argv, i)
			if err != nil {
				return opts, err
			}
			opts.After, opts.Before, i = n, n, ni
		case "-c":
			opts.CountOnly = true
		case "-i":
			opts.IgnoreCase = true
		case "-v":
			opts.Invert = true
		case "-F":
			opts.Fixed = true
		case "-n":
			opts.WithLineNum = true
		default:
			if opts.Pattern == "" {
				opts.Pattern = argv[i]
			} else {
				opts.Filename = argv[i]
			}
		}
	}

	if opts.Pattern == "" {
		return opts, fmt.Errorf("missing search pattern")
	}

	return opts, nil
}

// parseIntFlag parses an integer argument that follows a flag such as -A, -B, or -C.
// It returns the parsed integer value, the index of the argument just processed (i+1),
// and an error if the value is missing or not a valid integer.
func parseIntFlag(args []string, i int) (int, int, error) {
	if i+1 >= len(args) {
		return 0, i, fmt.Errorf("missing int value for %s", args[i])
	}
	n, err := strconv.Atoi(args[i+1])
	if err != nil {
		return 0, i, fmt.Errorf("error value to %s: %v", args[i], err)
	}

	return n, i + 1, nil
}
