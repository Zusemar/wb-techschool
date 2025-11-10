package matcher

import (
	"regexp"
	"strings"
)

// Matcher decides if a given string matches the pattern.
// provides united interface to diff -v -i -F flags
type Matcher interface {
	Match(str string) bool
}

// matcher
type matcher struct {
	pattern    string
	fixed      bool
	ignoreCase bool
	invert     bool
	re         *regexp.Regexp
}

// Match
func (m *matcher) Match(str string) bool {
	var matched bool

	if m.fixed {
		if m.ignoreCase {
			matched = strings.Contains(strings.ToLower(str), strings.ToLower(m.pattern))
		} else {
			matched = strings.Contains(str, m.pattern)
		}
	} else {
		matched = m.re.MatchString(str)
	}
	if m.invert {
		matched = !matched
	}

	return matched
}

// New is a fabric of matchers
// TODO: if args will become too many its better to change signature
// to New(pattern string, cfg confix) where confix contains all flags
func New(pattern string, fixed, ignoreCase, invert bool) (Matcher, error) {
	m := &matcher{
		pattern:    pattern,
		fixed:      fixed,
		ignoreCase: ignoreCase,
		invert:     invert,
	}

	// если не fixed, значит это регулярка
	if !fixed {
		reStr := pattern
		if ignoreCase {
			reStr = "(?i)" + reStr
		}
		re, err := regexp.Compile(reStr)
		if err != nil {
			return nil, err
		}
		m.re = re
	}

	return m, nil
}
