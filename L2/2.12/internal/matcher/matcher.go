package matcher

import (
	"strings"
)

// Matcher provides united interface to diff -v -i -F flags
type Matcher interface {
	Match(str, pattern string) bool
}

// fixedMatcher handles the Match with patter in -F flag situations
type fMatcher struct {
}

func (fm *fMatcher) Match(str, pattern string) bool {
	return strings.EqualFold(str, pattern)
}

// vMatcher handles the Match with patter in -v flag situations
type vMatcher struct {
}

func (vm *vMatcher) Match(str, pattern string) bool {
	return !strings.EqualFold(str, pattern)
}

// iMatcher handles the Match with patter in -i flag situations
type iMatcher struct {
}

func (im *iMatcher) Match(str, pattern string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(pattern))
}
