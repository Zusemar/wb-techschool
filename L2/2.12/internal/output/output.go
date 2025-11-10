package output

import (
	"fmt"
	"io"

	"grep/internal/parser"
	"grep/internal/reader"
)

// Formatter defines how grep results are printed.
// It handles formatting for matches, context lines, and counts.
type Formatter interface {
	PrintMatch(line reader.Line)
	PrintContext(line reader.Line)
	PrintCount(n int)
}

// formatter is the default implementation of Formatter.
// It writes output to an io.Writer (usually os.Stdout).
type formatter struct {
	w    io.Writer
	opts parser.Options
}

// NewFormatter creates a new Formatter with the given writer and options.
func NewFormatter(w io.Writer, opts parser.Options) Formatter {
	return &formatter{w: w, opts: opts}
}

// PrintMatch prints a line that matched the pattern.
// It respects the -n flag (show line numbers).
func (f *formatter) PrintMatch(line reader.Line) {
	if f.opts.WithLineNum {
		fmt.Fprintf(f.w, "%d:%s\n", line.Num, line.Text)
	} else {
		fmt.Fprintf(f.w, "%s\n", line.Text)
	}
}

// PrintContext prints a line that is part of before/after context (-A, -B).
// It uses a different prefix (for clarity) and supports -n as well.
func (f *formatter) PrintContext(line reader.Line) {
	if f.opts.WithLineNum {
		fmt.Fprintf(f.w, "%d-%s\n", line.Num, line.Text)
	} else {
		fmt.Fprintf(f.w, "%s\n", line.Text)
	}
}

// PrintCount prints only the total number of matches (-c flag).
func (f *formatter) PrintCount(n int) {
	fmt.Fprintf(f.w, "%d\n", n)
}
