package internal

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// Reader — интерфейс абстрактного потокового ридера shell.
// Возвращает одну строку без символов \n и \r.
// При достижении EOF возвращает ("", io.EOF).
type Reader interface {
	Next() (string, error)
}

// reader — обёртка над bufio.Reader для удобного чтения строк из stdin или pipe.
type reader struct {
	r *bufio.Reader
}

// NewReader создаёт новый ридер, принимающий любой io.Reader (stdin, файл, pipe).
func NewReader(src io.Reader) Reader {
	return &reader{r: bufio.NewReader(src)}
}

// Next читает одну строку до '\n' и возвращает её без завершающих переносов.
// Поведение:
//   - обычная строка: "cmd arg1 arg2"
//   - пустая строка: ""
//   - Ctrl+D (EOF): возвращает ("", io.EOF)
//   - любая другая ошибка чтения: возвращает ошибку
func (r *reader) Next() (string, error) {
	line, err := r.r.ReadString('\n')

	if err != nil {
		// EOF: пользователь нажал Ctrl+D или поток закончился
		if errors.Is(err, io.EOF) {
			// Если что-то было прочитано до EOF — вернуть эту строку
			if len(line) > 0 {
				return strings.TrimRight(line, "\r\n"), nil
			}
			// Если вообще ничего не было — это сигнал выхода из shell
			return "", io.EOF
		}
		// Любая другая ошибка — фатальная
		return "", err
	}

	// Убираем \n и \r
	return strings.TrimRight(line, "\r\n"), nil
}
