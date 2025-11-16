package writer

import (
	"bufio"
	"io"
)

// Writer просто удобный нейминг для врайтера
type Writer struct {
	w *bufio.Writer
}

// New создаёт буферизированный writer поверх stdout.
func New(out io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(out)}
}

// Write записывает строку в stdout с переводом строки.
func (wr *Writer) Write(line string) error {
	if _, err := wr.w.WriteString(line + "\n"); err != nil {
		return err
	}
	// Flush сразу, чтобы результат был виден немедленно.
	return wr.w.Flush()
}
