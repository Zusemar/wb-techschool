package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"wb-techschool/L2/2.16/internal"
)

var (
	commands = map[string]func([]string){
		"wget": handleWget,
		"help": handleHelp,
		"exit": handleExit,
	}
)

func main() {
	in := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("wgetlite> ")

		if !in.Scan() {
			fmt.Println()
			break // EOF (Ctrl+D)
		}

		tokens := tokenize(in.Text())
		if len(tokens) == 0 {
			continue
		}

		cmd := tokens[0]
		args := tokens[1:]

		if handle, ok := commands[cmd]; ok {
			handle(args)
		} else {
			handleUnknown(cmd)
		}
	}
}

func handleUnknown(cmd string) {
	fmt.Println("unknown command:", cmd, "Try 'help' for available commands")
}

func handleExit(args []string) {
	os.Exit(0)
}

func handleHelp(args []string) {
	fmt.Println("Available commands:")
	fmt.Println("  wget <url> [depth] - mirror site starting from URL with optional recursion depth (default 1)")
	fmt.Println("  help               - show this help")
	fmt.Println("  exit               - exit program")
}

func handleWget(args []string) {
	if len(args) == 0 {
		fmt.Println("usage: wget <url> [depth]")
		return
	}

	startURL := args[0]

	depth := 1
	if len(args) >= 2 {
		if d, err := strconv.Atoi(args[1]); err == nil && d >= 0 {
			depth = d
		}
	}

	// Корневая директория для сохранения — рядом с бинарником, папка "mirror".
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Println("error determining working directory:", err)
		return
	}

	targetDir := filepath.Join(rootDir, "mirror")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cfg := internal.MirrorConfig{
		RootDir:      targetDir,
		MaxDepth:     depth,
		RequestTTL:   20 * time.Second,
		MaxPages:     0,   // без жёсткого лимита по страницам
		SameHostOnly: true, // не выходим за пределы домена
	}

	fmt.Println("Starting mirror:")
	fmt.Println("  URL:   ", startURL)
	fmt.Println("  Depth: ", depth)
	fmt.Println("  Target:", targetDir)

	if err := internal.RunMirror(ctx, startURL, cfg); err != nil {
		fmt.Println("wget error:", err)
		return
	}

	fmt.Println("Mirror completed.")
}

func tokenize(text string) []string {
	return strings.Fields(text)
}
