package reader

import (
	"bufio"
	"io"
)

// Reader просто удобынй нейминг для сканера
type Reader struct {
	scanner *bufio.Scanner
}

// New создает новый сканер
func New(r io.Reader) *Reader {
	sc := bufio.NewScanner(r)
	buf := make([]byte, 0, 1024)
	sc.Buffer(buf, 10*1024*1024) // увеличиваем лимит до 10 МБ, вдруг там большие строки
	return &Reader{scanner: sc}
}

// Next дает следующую строку
func (r *Reader) Next() (string, bool, error) {
	if r.scanner.Scan() {
		return r.scanner.Text(), true, nil
	}
	if err := r.scanner.Err(); err != nil {
		return "", false, err
	}
	return "", false, nil
}
