package engine

import (
	"container/list"

	"grep/internal/matcher"
	"grep/internal/output"
	"grep/internal/parser"
	"grep/internal/reader"
)

// Engine orchestrates the grep process — reading, matching, and outputting lines.
type Engine struct {
	m    matcher.Matcher
	opts parser.Options
}

// New creates a new Engine with the given matcher and parsed options.
func New(m matcher.Matcher, opts parser.Options) *Engine {
	return &Engine{m: m, opts: opts}
}

// Run executes the main grep logic using the provided Reader and Formatter.
// It supports flags: -A, -B, -C, -c, -n.
func (e *Engine) Run(src reader.Reader, out output.Formatter) error {
	defer src.Close()

	var (
		matchCount   int
		afterCounter int
		beforeBuffer = list.New() // ring buffer for -B
	)

	for {
		line, ok, err := src.Next()
		if err != nil {
			return err
		}
		if !ok {
			break
		}

		matched := e.m.Match(line.Text)

		if matched {
			matchCount++

			// Handle count-only mode (-c)
			if e.opts.CountOnly {
				continue
			}

			// Print lines before match (-B)
			e.flushBefore(out, beforeBuffer)

			// Print matched line
			out.PrintMatch(line)

			// Activate after context (-A)
			if e.opts.After > 0 {
				afterCounter = e.opts.After
			}
		} else {
			// Save previous lines for -B
			e.rememberBefore(beforeBuffer, line)

			// Print "after" context if active
			if afterCounter > 0 {
				out.PrintContext(line)
				afterCounter--
			}
		}
	}

	// Final count output if -c flag is active
	if e.opts.CountOnly {
		out.PrintCount(matchCount)
	}

	return nil
}

// rememberBefore stores the current line into a fixed-size buffer for -B context.
func (e *Engine) rememberBefore(buf *list.List, line reader.Line) {
	if e.opts.Before == 0 {
		return
	}

	if buf.Len() >= e.opts.Before {
		buf.Remove(buf.Front()) // discard oldest
	}
	buf.PushBack(line)
}

// flushBefore prints all stored "before" context lines.
func (e *Engine) flushBefore(out output.Formatter, buf *list.List) {
	for elem := buf.Front(); elem != nil; elem = elem.Next() {
		line := elem.Value.(reader.Line)
		out.PrintContext(line)
	}
	buf.Init() // clear buffer after printing
}
