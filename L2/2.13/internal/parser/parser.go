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
	opts := Options{}
	if len(args) == 0 {
		return Options{}, fmt.Errorf("usage: cut -f LIST [-d DELIM] [-s]")
	}

	for i := 1; i < len(args); i++ {
		switch args[i] {
		// -f Mask may return boolean mask that contains bits that are not into stdin lines
		case "-f":
			if i+1 >= len(args) {
				return opts, errors.New("missing value for -f")
			}
			opts.Mask = make(map[int]bool)
			values := strings.Split(args[i+1], ",")
			for _, val := range values {
				// на вход дали диапазон
				if strings.Contains(val, "-") {
					rang := strings.Split(val, "-")
					from, err := strconv.Atoi(rang[0])
					if err != nil {
						return opts, fmt.Errorf("invalid argument %v \n error: %s", val, err)
					}
					to, err := strconv.Atoi(rang[1])
					if err != nil {
						return opts, fmt.Errorf("invalid argument %v \n error: %s", val, err)
					}

					for i := from; i <= to; i++ {
						opts.Mask[i-1] = true
					}
				} else { // на вход дали просто число
					ival, err := strconv.Atoi(val)
					if err != nil {
						return opts, fmt.Errorf("invalid argument %v \n error: %s", val, err)
					}
					opts.Mask[ival-1] = true
				}

			}
			i++
		// Не считаю ошибкой -d без самого разделителя, тогда просто пробел по дефолту
		case "-d":
			if i+1 >= len(args) {
				opts.Delimiter = '\t'
				return opts, nil
			}
			opts.Delimiter = []rune(args[i+1])[0] // аххааххахахааххаха
			i++
		case "-s":
			opts.Separated = true
		}
	}
	if opts.Mask == nil {
		return opts, fmt.Errorf("missing -f call")
	}
	return opts, nil
}
