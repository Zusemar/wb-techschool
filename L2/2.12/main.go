package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type grepArgs struct {
	pattern  string
	filename string
	an       int
	bn       int
	cn       int
	A        bool
	B        bool
	C        bool
	c        bool
	i        bool
	v        bool
	F        bool
	n        bool
}

func parseArgs(args []string) grepArgs {
	opts := grepArgs{}

argsLoop:
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-A":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "missing N argument in -A call")
				break argsLoop
			}
			an, err := strconv.Atoi(args[i+1])
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				break argsLoop
			}
			opts.A = true
			opts.an = an
			i++
		case "-B":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "missing N argument in -B call")
				break argsLoop
			}
			bn, err := strconv.Atoi(args[i+1])
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				break argsLoop
			}
			opts.B = true
			opts.bn = bn
			i++
		case "-C":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "missing N argument in -C call")
				break argsLoop
			}
			cn, err := strconv.Atoi(args[i+1])
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				break argsLoop
			}
			opts.C = true
			opts.cn = cn
			i++
		case "-c":
			opts.c = true
		case "-i":
			opts.i = true
		case "-v":
			opts.v = true
		case "-F":
			opts.F = true
		case "-n":
			opts.n = true
		default:
			// Всё, что не флаг — это либо шаблон, либо файл
			if opts.pattern == "" {
				opts.pattern = args[i]

			} else {
				opts.filename = args[i]
			}
		}
	}

	return opts
}

func grep(lines []string, opts grepArgs) []string {

	if opts.i {
		for i := 0; i < len(lines); i++ {
			lines[i] = strings.ToLower(lines[i])
		}
	}

	// TODO:
	if opts.n {

	}

	//  TODO:
	if opts.F {

	}

	// TODO:
	if opts.v {

	}

	// TODO:
	if opts.A {

		out := []string{}

		for i := 0; i < len(lines); i++ {
			if opts.F {
				if strings.EqualFold(lines[i], opts.pattern) {
					// Добавляем найденную строку
					out = append(out, lines[i])

					// Добавляем N строк после неё (контекст)
					for j := 1; j <= opts.an && i+j < len(lines); j++ {
						out = append(out, lines[i+j])
					}
				}
			} else {

				if strings.Contains(lines[i], opts.pattern) {
					// Добавляем найденную строку
					out = append(out, lines[i])

					// Добавляем N строк после неё (контекст)
					for j := 1; j <= opts.an && i+j < len(lines); j++ {
						out = append(out, lines[i+j])
					}
				}
			}

		}
		return out
	}

	// TODO: тут же только первое нахождение строки будет обрабатываться
	if opts.B {
		out := []string{}
		idx := 0

		if opts.F {
			for i := 0; i < len(lines); i++ {
				out = append(out, lines[i])
				if strings.EqualFold(lines[i], opts.pattern) {
					// Добавляем найденную строку
					out = append(out, lines[i])
					idx = i
					break
				}
			}
			return out[idx-opts.bn : idx]
		} else {
			for i := 0; i < len(lines); i++ {
				out = append(out, lines[i])
				if strings.Contains(lines[i], opts.pattern) {
					// Добавляем найденную строку
					out = append(out, lines[i])
					idx = i
					break
				}
			}
			return out[idx-opts.bn : idx]
		}
	}
	// TODO:
	if opts.C {
		out := []string{}
		idx := 0
		for i := 0; i < len(lines); i++ {
			if opts.F {
				if strings.EqualFold(lines[i], opts.pattern) {
					// Добавляем найденную строку
					out = append(out, lines[i])
					idx = i
					// Добавляем N строк после неё (контекст)
					for j := 1; j <= opts.an && i+j < len(lines); j++ {
						out = append(out, lines[i+j])
					}
				}
			} else {

				if strings.Contains(lines[i], opts.pattern) {
					// Добавляем найденную строку
					out = append(out, lines[i])
					idx = i
					// Добавляем N строк после неё (контекст)
					for j := 1; j <= opts.an && i+j < len(lines); j++ {
						out = append(out, lines[i+j])
					}
				}
			}

		}
		return out[idx-opts.cn : idx+opts.cn]
	}

	// TODO:
	if opts.c {
		ctr := 0
		out := []string{}

		if opts.F {
			for i := 0; i < len(lines); i++ {
				out = append(out, lines[i])
				if strings.EqualFold(lines[i], opts.pattern) {
					// Добавляем найденную строку
					ctr++
					out = append(out, string(ctr))
				}
			}
			return out
		} else {
			for i := 0; i < len(lines); i++ {
				out = append(out, lines[i])
				if strings.Contains(lines[i], opts.pattern) {
					// Добавляем найденную строку
					ctr++
					out = append(out, string(ctr))
				}
			}
			return out
		}
	}

	return []string{}
}

func main() {
	args := os.Args[1:]
	opts := parseArgs(args)
	// Определяем источник данных
	var reader io.Reader
	if len(args) > 0 && !strings.HasPrefix(args[len(args)-1], "-") {
		file, err := os.Open(args[len(args)-1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка открытия файла:", err)
			return
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	// Читаем построчно
	scanner := bufio.NewScanner(reader)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения:", err)
		return
	}

	grep(lines, opts)
}
