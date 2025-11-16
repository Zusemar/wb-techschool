package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Options structure contains all flags info. Map is bitmap showing each
// position of column that must be shown
type Options struct {
	Delimiter rune
	Separated bool
	Mask      map[int]bool
}

// Parse gets all the flags from args = os.Args[1:], fills Options struct
func Parse(args []string) (Options, error) {
	opts := Options{
		Delimiter: '\t', // default delimiter is tab
		Mask:      make(map[int]bool),
	}
	if len(args) == 0 {
		return Options{}, fmt.Errorf("usage: cut -f LIST [-d DELIM(optional)] [-s]")
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		// -f Mask may return boolean mask that contains bits that are not into stdin lines
		case "-f":
			if i+1 >= len(args) {
				return opts, errors.New("missing value for -f")
			}
			values := strings.Split(args[i+1], ",")
			for _, val := range values {
				val = strings.TrimSpace(val)
				if val == "" {
					continue
				}
				// на вход дали диапазон
				if strings.Contains(val, "-") {
					rang := strings.Split(val, "-")
					if len(rang) != 2 {
						return opts, fmt.Errorf("invalid range format: %v", val)
					}
					if rang[0] == "" || rang[1] == "" {
						return opts, fmt.Errorf("invalid range format: %v", val)
					}
					from, err := strconv.Atoi(rang[0])
					if err != nil {
						return opts, fmt.Errorf("invalid argument %v: %s", val, err)
					}
					to, err := strconv.Atoi(rang[1])
					if err != nil {
						return opts, fmt.Errorf("invalid argument %v: %s", val, err)
					}
					if from > to {
						return opts, fmt.Errorf("invalid range: %v (from > to)", val)
					}
					if from < 1 {
						return opts, fmt.Errorf("invalid field number: %d (must be >= 1)", from)
					}

					for j := from; j <= to; j++ {
						opts.Mask[j-1] = true
					}
				} else { // на вход дали просто число
					ival, err := strconv.Atoi(val)
					if err != nil {
						return opts, fmt.Errorf("invalid argument %v: %s", val, err)
					}
					if ival < 1 {
						return opts, fmt.Errorf("invalid field number: %d (must be >= 1)", ival)
					}
					opts.Mask[ival-1] = true
				}
			}
			i++
		case "-d":
			if i+1 >= len(args) {
				opts.Delimiter = '\t'
			} else {
				if len(args[i+1]) == 0 {
					return opts, errors.New("delimiter cannot be empty")
				}
				runes := []rune(args[i+1])
				if len(runes) != 1 {
					return opts, fmt.Errorf("delimiter must be a single character, got: %v", args[i+1])
				}
				opts.Delimiter = runes[0]
				i++
			}
		case "-s":
			opts.Separated = true
		default:
			return opts, fmt.Errorf("unknown flag: %v", args[i])
		}
	}

	if len(opts.Mask) == 0 {
		return opts, errors.New("field list (-f) is required")
	}

	return opts, nil
}

