package internal

type Token struct {
	Type  TokenType
	Value string
}

type TokenType int

const (
	TOKEN_WORD TokenType = iota
	TOKEN_PIPE
	TOKEN_AND
	TOKEN_OR
	TOKEN_REDIR_IN
	TOKEN_REDIR_OUT
	TOKEN_EOF
)

type State int

const (
	STATE_NORMAL State = iota
	STATE_IN_DOUBLE
	STATE_IN_SINGLE
	STATE_ESCAPED
)

type Tokenizer struct {
	input     []rune
	state     State
	prevState State
	pos       int
	next      rune
	buf       []rune
	tokens    []Token
}

/* ------------------------- PUBLIC TOKENIZE ---------------------------- */

func Tokenize(str string) ([]Token, error) {
	t := &Tokenizer{}
	return t.Tokenize(str)
}

func (t *Tokenizer) Tokenize(str string) ([]Token, error) {
	t.input = []rune(str)
	t.state = STATE_NORMAL
	t.prevState = STATE_NORMAL
	t.pos = 0
	t.buf = t.buf[:0]
	t.tokens = t.tokens[:0]
	//fmt.Println(t.tokens)
	for t.pos < len(t.input) {
		//fmt.Println(t.state)
		//fmt.Println(t.tokens)
		ch := t.input[t.pos]

		if t.pos+1 < len(t.input) {
			t.next = t.input[t.pos+1]
		} else {
			t.next = 0
		}

		switch t.state {

		case STATE_NORMAL:
			t.handleNormal(ch)

		case STATE_IN_DOUBLE:
			t.handleDouble(ch)

		case STATE_IN_SINGLE:
			t.handleSingle(ch)

		case STATE_ESCAPED:
			t.handleEscaped(ch)
		}

		t.pos++
	}

	// закрываем последний WORD, если он есть
	t.flushToken()

	// EOF-токен
	t.tokens = append(t.tokens, Token{Type: TOKEN_EOF, Value: ""})

	return t.tokens, nil
}

/* ---------------------------- HELPERS -------------------------------- */

func (t *Tokenizer) flushToken() {
	if len(t.buf) == 0 {
		return
	}
	t.tokens = append(t.tokens, Token{
		Type:  TOKEN_WORD,
		Value: string(t.buf),
	})
	t.buf = t.buf[:0]
}

func (t *Tokenizer) addToken(tt TokenType, val string) {
	t.tokens = append(t.tokens, Token{
		Type:  tt,
		Value: val,
	})
}

/* ---------------------------- STATE HANDLERS -------------------------- */

func (t *Tokenizer) handleNormal(ch rune) {
	switch ch {

	case ' ':
		t.flushToken()
		return

	case '"':
		t.state = STATE_IN_DOUBLE
		return

	case '\'':
		t.state = STATE_IN_SINGLE
		return

	case '\\':
		t.prevState = STATE_NORMAL
		t.state = STATE_ESCAPED
		return

	case '|':
		t.flushToken()
		if t.next == '|' {
			t.addToken(TOKEN_OR, "||")
			t.pos++ // пропускаем второй |
		} else {
			t.addToken(TOKEN_PIPE, "|")
		}
		return

	case '&':
		t.flushToken()
		if t.next == '&' {
			t.addToken(TOKEN_AND, "&&")
			t.pos++ // пропускаем второй &
		}
		// одиночный & игнорируем
		return

	case '<':
		t.flushToken()
		t.addToken(TOKEN_REDIR_IN, "<")
		return

	case '>':
		t.flushToken()
		t.addToken(TOKEN_REDIR_OUT, ">")
		return

	default:
		t.buf = append(t.buf, ch)
	}
}

func (t *Tokenizer) handleDouble(ch rune) {
	switch ch {
	case '"':
		t.state = STATE_NORMAL
		return
	case '\\':
		t.prevState = STATE_IN_DOUBLE
		t.state = STATE_ESCAPED
		return
	default:
		t.buf = append(t.buf, ch)
	}
}

func (t *Tokenizer) handleSingle(ch rune) {
	switch ch {
	case '\'':
		t.state = STATE_NORMAL
		return
	case '\\':
		t.prevState = STATE_IN_SINGLE
		t.state = STATE_ESCAPED
		return
	default:
		t.buf = append(t.buf, ch)
	}
}

func (t *Tokenizer) handleEscaped(ch rune) {
	t.buf = append(t.buf, ch)
	t.state = t.prevState
}
