package core

import (
	"cut/internal/parser"
	"cut/internal/reader"
	"cut/internal/writer"
	"os"
	"strings"
)

// Core эт структурка где ядро хранит все что надо для работы
type Core struct {
	r    *reader.Reader
	w    *writer.Writer
	opts parser.Options
}

// New создает новое ядро
func New(args []string) (*Core, error) {
	opts, err := parser.Parse(args)
	if err != nil {
		return nil, err
	}
	return &Core{
		r:    reader.New(os.Stdin),
		w:    writer.New(os.Stdout),
		opts: opts,
	}, nil
}

// Process обрабатывает одну строку
func (c *Core) Process(line string) (string, bool) {
	d := c.opts.Delimiter
	if c.opts.Separated && !strings.ContainsRune(line, d) {
		return "", false
	}

	fields := strings.Split(line, string(d))
	var b strings.Builder
	first := true
	for i, f := range fields {
		// Проверяем, что поле существует в маске (игнорируем поля за границами)
		if c.opts.Mask != nil && c.opts.Mask[i] {
			if !first {
				b.WriteRune(d)
			}
			b.WriteString(f)
			first = false
		}
	}

	if b.Len() == 0 {
		return "", false
	}
	return b.String(), true
}

// Run соединяет reader, core и writer
func (c *Core) Run() error {
	for {
		line, ok, err := c.r.Next()
		if err != nil {
			return err
		}
		if !ok {
			break
		}

		if out, ok2 := c.Process(line); ok2 {
			if err := c.w.Write(out); err != nil {
				return err
			}
		}
	}
	return nil
}
