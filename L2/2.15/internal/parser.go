package internal

import (
	"fmt"
)

type Node interface{}

type Command struct {
	Name        string
	Args        []string
	RedirectIn  string
	RedirectOut string
}

type Pipeline struct {
	Commands []*Command
}

type LogicalChain struct {
	Left  Node
	Op    TokenType // TOKEN_AND или TOKEN_OR
	Right Node
}

type parser struct {
	tokens []Token
	pos    int
}

func Parse(tokens []Token) (Node, error) {
	p := &parser{tokens: tokens, pos: 0}
	node, err := p.parseLogical()
	if err != nil {
		return nil, err
	}
	// ждём EOF
	if tok := p.current(); tok.Type != TOKEN_EOF {
		return nil, fmt.Errorf("unexpected token at end: %v", tok)
	}
	return node, nil
}

/* ---------------- low-level helpers ---------------- */

func (p *parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TOKEN_EOF}
	}
	return p.tokens[p.pos]
}

func (p *parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

/* ---------------- parsing ---------------- */

// logical := pipeline ( (AND|OR) pipeline )*
func (p *parser) parseLogical() (Node, error) {
	left, err := p.parsePipeline()
	if err != nil {
		return nil, err
	}

	for {
		tok := p.current()
		if tok.Type != TOKEN_AND && tok.Type != TOKEN_OR {
			break
		}
		op := tok.Type
		p.advance() // съели AND/OR

		right, err := p.parsePipeline()
		if err != nil {
			return nil, err
		}
		left = &LogicalChain{
			Left:  left,
			Op:    op,
			Right: right,
		}
	}

	return left, nil
}

// pipeline := command ( '|' command )*
func (p *parser) parsePipeline() (Node, error) {
	cmd, err := p.parseCommand()
	if err != nil {
		return nil, err
	}
	commands := []*Command{cmd}

	for {
		tok := p.current()
		if tok.Type != TOKEN_PIPE {
			break
		}
		p.advance() // съели '|'

		nextCmd, err := p.parseCommand()
		if err != nil {
			return nil, err
		}
		commands = append(commands, nextCmd)
	}

	if len(commands) == 1 {
		return commands[0], nil
	}
	return &Pipeline{Commands: commands}, nil
}

// command := WORD {WORD} {redir}
func (p *parser) parseCommand() (*Command, error) {
	tok := p.current()
	if tok.Type != TOKEN_WORD {
		return nil, fmt.Errorf("expected command name, got: %v", tok)
	}
	cmd := &Command{
		Name: tok.Value,
	}
	p.advance()

	for {
		tok = p.current()
		switch tok.Type {
		case TOKEN_WORD:
			cmd.Args = append(cmd.Args, tok.Value)
			p.advance()

		case TOKEN_REDIR_IN:
			p.advance()
			arg := p.current()
			if arg.Type != TOKEN_WORD {
				return nil, fmt.Errorf("expected filename after '<', got: %v", arg)
			}
			cmd.RedirectIn = arg.Value
			p.advance()

		case TOKEN_REDIR_OUT:
			p.advance()
			arg := p.current()
			if arg.Type != TOKEN_WORD {
				return nil, fmt.Errorf("expected filename after '>', got: %v", arg)
			}
			cmd.RedirectOut = arg.Value
			p.advance()

		default:
			return cmd, nil
		}
	}
}
