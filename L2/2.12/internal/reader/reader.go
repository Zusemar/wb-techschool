package reader

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

// Line represents a single line from the input stream along with its metadata.
type Line struct {
	Num  int    // line number (starting from 1)
	Text string // line content without newline characters
	File string // file name (empty if reading from stdin)
}

// Reader defines an interface for reading lines from any source.
// It returns one line per Next() call.
type Reader interface {
	Next() (Line, bool, error) // returns a line, a continuation flag, and an error
	Close() error              // closes the source (noop for stdin)
}

// fileReader is a universal Reader implementation for file or stdin.
type fileReader struct {
	file *os.File
	buf  *bufio.Reader
	num  int
	name string
}

// NewFromFile creates a Reader for reading from the specified file.
func NewFromFile(path string) (Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &fileReader{
		file: f,
		buf:  bufio.NewReader(f),
		name: path,
	}, nil
}

// NewFromStdin creates a Reader for reading from standard input.
func NewFromStdin() Reader {
	return &fileReader{
		file: os.Stdin,
		buf:  bufio.NewReader(os.Stdin),
		name: "",
	}
}

// Next reads the next line.
// Returns (Line{}, false, nil) when the stream ends.
func (r *fileReader) Next() (Line, bool, error) {
	b, err := r.buf.ReadBytes('\n')

	if err != nil {
		if errors.Is(err, io.EOF) {
			// EOF without data to end of stream
			if len(b) == 0 {
				return Line{}, false, nil
			}
			// last line without trailing newline
			r.num++
			text := strings.TrimRight(string(b), "\r\n")
			return Line{Num: r.num, Text: text, File: r.name}, true, nil
		}
		return Line{}, false, err // actual I/O error
	}

	// normalize line endings (\r\n to \n)
	text := strings.TrimRight(string(b), "\r\n")
	r.num++
	return Line{Num: r.num, Text: text, File: r.name}, true, nil
}

// Close closes the source.
// For stdin this is a no-op.
func (r *fileReader) Close() error {
	if r.file == os.Stdin {
		return nil
	}
	return r.file.Close()
}
