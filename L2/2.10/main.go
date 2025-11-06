package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

// sortOptionsструктура для хранения флагов
type sortOptions struct {
	N int
	k bool
	n bool
	r bool
	u bool
	m bool
	b bool
	c bool
	h bool
}

// -------------------------
// Парсинг флагов
// -------------------------
func parseFlags(args []string) sortOptions {
	opts := sortOptions{}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if !strings.HasPrefix(arg, "-") {
			continue // не флаг — возможно имя файла
		}

		switch {
		case arg == "-k":
			opts.k = true
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "Ошибка: флаг -k требует аргумент (номер колонки)")
				break
			}
			n, err := strconv.Atoi(args[i+1])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка: некорректное значение для -k:", args[i+1])
				break
			}
			opts.N = n
			i++ // пропускаем значение после -k

		case strings.HasPrefix(arg, "-") && len(arg) > 2:
			for _, r := range arg[1:] {
				switch r {
				case 'n':
					opts.n = true
				case 'r':
					opts.r = true
				case 'u':
					opts.u = true
				case 'M':
					opts.m = true
				case 'b':
					opts.b = true
				case 'c':
					opts.c = true
				case 'h':
					opts.h = true
				default:
					fmt.Fprintf(os.Stderr, "Неизвестный флаг: -%c\n", r)
				}
			}

		default:
			switch arg {
			case "-n":
				opts.n = true
			case "-r":
				opts.r = true
			case "-u":
				opts.u = true
			case "-M":
				opts.m = true
			case "-b":
				opts.b = true
			case "-c":
				opts.c = true
			case "-h":
				opts.h = true
			default:
				fmt.Fprintf(os.Stderr, "Неизвестный флаг: %s\n", arg)
			}
		}
	}

	return opts
}

// -------------------------
// Основная функция сортировки
// -------------------------
func sortStr(text []string, opts sortOptions) []string {
	if opts.c {
		// режим проверки сортировки (-c)
		for i := 1; i < len(text); i++ {
			if text[i-1] > text[i] {
				fmt.Println("Данные не отсортированы.")
				os.Exit(1)
			}
		}
		fmt.Println("Данные отсортированы.")
		os.Exit(0)
	}

	out := append([]string(nil), text...) // копия входных строк

	sort.Slice(out, func(i, j int) bool {
		a := out[i]
		b := out[j]

		// игнорируем хвостовые пробелы (-b)
		if opts.b {
			a = strings.TrimSpace(a)
			b = strings.TrimSpace(b)
		}

		// сортировка по колонке (-k)
		if opts.k {
			fieldsA := strings.Split(a, "\t")
			fieldsB := strings.Split(b, "\t")
			col := opts.N - 1
			if col < len(fieldsA) {
				a = fieldsA[col]
			}
			if col < len(fieldsB) {
				b = fieldsB[col]
			}
		}

		// числовая сортировка (-n)
		if opts.n {
			na, errA := strconv.ParseFloat(a, 64)
			nb, errB := strconv.ParseFloat(b, 64)
			if errA == nil && errB == nil {
				if opts.r {
					return na > nb
				}
				return na < nb
			}
		}

		// строковая сортировка (по умолчанию)
		if opts.r {
			return a > b
		}
		return a < b
	})

	// уникальные строки (-u)
	if opts.u {
		out = unique(out)
	}

	return out
}

// -------------------------
// Удаление дубликатов
// -------------------------
func unique(s []string) []string {
	if len(s) == 0 {
		return s
	}
	res := []string{s[0]}
	for i := 1; i < len(s); i++ {
		if s[i] != s[i-1] {
			res = append(res, s[i])
		}
	}
	return res
}

// -------------------------
// Основной main()
// -------------------------
func main() {
	args := os.Args[1:]
	opts := parseFlags(args)

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

	// Выполняем сортировку
	result := sortStr(lines, opts)

	// Печатаем результат
	for _, line := range result {
		fmt.Println(line)
	}
}
