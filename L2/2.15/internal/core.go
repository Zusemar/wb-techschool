package internal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

// Execute выполняет AST-узел и возвращает exit code.
func Execute(node Node, stdin io.Reader, stdout, stderr io.Writer) int {
	switch n := node.(type) {
	case *Command:
		return execCommand(n, stdin, stdout, stderr)
	case *Pipeline:
		return execPipeline(n, stdin, stdout, stderr)
	case *LogicalChain:
		return execLogical(n, stdin, stdout, stderr)
	default:
		fmt.Fprintln(stderr, "unknown AST node")
		return 1
	}
}

/* ---------------- logical ---------------- */

func execLogical(chain *LogicalChain, stdin io.Reader, stdout, stderr io.Writer) int {
	leftCode := Execute(chain.Left, stdin, stdout, stderr)

	switch chain.Op {
	case TOKEN_AND:
		if leftCode == 0 {
			return Execute(chain.Right, stdin, stdout, stderr)
		}
		return leftCode
	case TOKEN_OR:
		if leftCode != 0 {
			return Execute(chain.Right, stdin, stdout, stderr)
		}
		return leftCode
	default:
		fmt.Fprintln(stderr, "invalid logical op")
		return 1
	}
}

/* ---------------- pipeline (через буферы) ---------------- */

func execPipeline(p *Pipeline, stdin io.Reader, stdout, stderr io.Writer) int {
	var in io.Reader = stdin
	var exitCode int

	for i, cmd := range p.Commands {
		isLast := i == len(p.Commands)-1

		var outBuf bytes.Buffer
		var out io.Writer
		if isLast {
			out = stdout
		} else {
			out = &outBuf
		}

		exitCode = execCommand(cmd, in, out, stderr)

		if !isLast {
			in = &outBuf
		}
	}

	return exitCode
}

/* ---------------- commands ---------------- */

func execCommand(cmd *Command, stdin io.Reader, stdout, stderr io.Writer) int {
	// builtins
	if handled, code := tryBuiltin(cmd, stdin, stdout, stderr); handled {
		return code
	}

	// внешняя команда
	return execExternal(cmd, stdin, stdout, stderr)
}

func execExternal(cmd *Command, stdin io.Reader, stdout, stderr io.Writer) int {
	if cmd.Name == "" {
		fmt.Fprintln(stderr, "empty command")
		return 1
	}

	c := exec.Command(cmd.Name, cmd.Args...)

	// stdin
	var in io.Reader = stdin
	if cmd.RedirectIn != "" {
		f, err := os.Open(cmd.RedirectIn)
		if err != nil {
			fmt.Fprintln(stderr, "open input:", err)
			return 1
		}
		defer f.Close()
		in = f
	}
	c.Stdin = in

	// stdout
	var out io.Writer = stdout
	if cmd.RedirectOut != "" {
		f, err := os.Create(cmd.RedirectOut)
		if err != nil {
			fmt.Fprintln(stderr, "create output:", err)
			return 1
		}
		defer f.Close()
		out = f
	}
	c.Stdout = out
	c.Stderr = stderr

	if err := c.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok2 := exitErr.Sys().(syscall.WaitStatus); ok2 {
				return status.ExitStatus()
			}
			return 1
		}
		fmt.Fprintln(stderr, "exec error:", err)
		return 1
	}
	return 0
}

/* ---------------- builtins ---------------- */

func tryBuiltin(cmd *Command, stdin io.Reader, stdout, stderr io.Writer) (bool, int) {
	switch cmd.Name {
	case "cd":
		return true, builtinCd(cmd, stderr)
	case "pwd":
		return true, builtinPwd(stdout, stderr)
	case "echo":
		return true, builtinEcho(cmd, stdout)
	case "kill":
		return true, builtinKill(cmd, stderr)
	case "ps":
		return true, builtinPs(stdout, stderr)
	default:
		return false, 0
	}
}

func builtinCd(cmd *Command, stderr io.Writer) int {
	dir := ""
	if len(cmd.Args) == 0 {
		dir = os.Getenv("HOME")
	} else {
		dir = cmd.Args[0]
	}
	if dir == "" {
		fmt.Fprintln(stderr, "cd: empty dir")
		return 1
	}
	if err := os.Chdir(dir); err != nil {
		fmt.Fprintln(stderr, "cd:", err)
		return 1
	}
	return 0
}

func builtinPwd(stdout, stderr io.Writer) int {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(stderr, "pwd:", err)
		return 1
	}
	fmt.Fprintln(stdout, wd)
	return 0
}

func builtinEcho(cmd *Command, stdout io.Writer) int {
	for i, a := range cmd.Args {
		if i > 0 {
			fmt.Fprint(stdout, " ")
		}
		fmt.Fprint(stdout, a)
	}
	fmt.Fprintln(stdout)
	return 0
}

func builtinKill(cmd *Command, stderr io.Writer) int {
	if len(cmd.Args) == 0 {
		fmt.Fprintln(stderr, "kill: missing pid")
		return 1
	}
	pid, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		fmt.Fprintln(stderr, "kill: invalid pid:", err)
		return 1
	}
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		fmt.Fprintln(stderr, "kill:", err)
		return 1
	}
	return 0
}

func builtinPs(stdout, stderr io.Writer) int {
	c := exec.Command("ps", "aux")
	c.Stdout = stdout
	c.Stderr = stderr
	if err := c.Run(); err != nil {
		fmt.Fprintln(stderr, "ps:", err)
		return 1
	}
	return 0
}
