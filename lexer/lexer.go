package lexer

import (
	"errors"
	"io"
	"unicode"
)

type TokenType int

const (
	TokenError TokenType = iota
	TokenEOF
	TokenStartTag
	TokenEndTag
	TokenAttributeName
	TokenAttributeValue
	TokenText
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input  io.RuneScanner
	buffer []rune
	Tokens []Token
}

type stateFn func(*Lexer) stateFn

func NewLexer(r io.RuneScanner) *Lexer {
	return &Lexer{
		input: r,
	}
}

func (l *Lexer) Run() {
	for state := lexData; state != nil; {
		state = state(l)
	}
}

func lexData(l *Lexer) stateFn {
	err := l.skipWhitespace()
	if errors.Is(err, io.EOF) {
		l.emit(TokenEOF)
		return nil
	}
	for {
		r, err := l.readNext()
		if err == io.EOF {
			l.emit(TokenEOF)
			return nil
		}

		switch r {
		case '<':
			if len(l.buffer) > 0 {
				l.emitToken(TokenText, string(l.buffer))
				l.clearRuneBuffer()
			}
			return lexTagOpen
		default:
			l.bufferRune(r)
		}
	}
}

func lexTagOpen(l *Lexer) stateFn {
	r, err := l.readNext()
	if err != nil {
		l.emit(TokenError)
		return nil
	}
	switch r {
	case '/':
		l.clearRuneBuffer()
		return lexEndTag
	default:
		l.clearRuneBuffer()
		l.bufferRune(r)
		return lexTagName
	}
}

func lexTagName(l *Lexer) stateFn {
	for {
		r, err := l.readNext()
		if err != nil {
			l.emit(TokenError)
			return nil
		}

		if unicode.IsUpper(r) {
			r += 0x0020
		}

		if unicode.IsSpace(r) || r == '>' || r == '/' {
			l.emitToken(TokenStartTag, string(l.buffer))
			l.clearRuneBuffer()
			switch r {
			case '>':
				return lexData
			case '/':
				return nil
			default:
				return lexBeforeAttributeName
			}
		}
		l.bufferRune(r)
	}
}

func lexEndTag(l *Lexer) stateFn {
	for {
		r, err := l.readNext()
		if err != nil {
			l.emit(TokenError)
			return nil
		}

		if r == '>' {
			l.emitToken(TokenEndTag, string(l.buffer))
			l.clearRuneBuffer()
			return lexData
		}

		l.bufferRune(r)
	}
}

func lexBeforeAttributeName(l *Lexer) stateFn {
	l.skipWhitespace()
	for {
		r, err := l.readNext()
		if err != nil {
			l.emit(TokenError)
			return nil
		}

		switch {
		case r == '>':
			return lexData
		case unicode.IsLetter(r):
			l.bufferRune(r)
			return lexAttributeName
		default:
			l.emit(TokenError)
			return nil
		}
	}
}

func lexAttributeName(l *Lexer) stateFn {
	for {
		r, err := l.readNext()
		if err != nil {
			l.emit(TokenError)
			return nil
		}

		if unicode.IsSpace(r) || r == '=' {
			l.emitToken(TokenAttributeName, string(l.buffer))
			l.clearRuneBuffer()
			return lexAttributeValue
		}
		l.bufferRune(r)
	}
}

func lexAttributeValue(l *Lexer) stateFn {
	quoteChar, err := l.readNext()
	if err != nil || (quoteChar != '"' && quoteChar != '\'') {
		l.emit(TokenError)
		return nil
	}

	for {
		r, err := l.readNext()
		if err != nil || r == quoteChar {
			l.emitToken(TokenAttributeValue, string(l.buffer))
			l.clearRuneBuffer()
			return lexBeforeAttributeName
		}
		l.bufferRune(r)
	}
}

func (l *Lexer) emit(t TokenType) {
	l.Tokens = append(l.Tokens, Token{Type: t})
}

func (l *Lexer) emitToken(t TokenType, value string) {
	l.Tokens = append(l.Tokens, Token{Type: t, Value: value})
}

func (l *Lexer) bufferRune(r rune) {
	l.buffer = append(l.buffer, r)
}

func (l *Lexer) clearRuneBuffer() {
	l.buffer = make([]rune, 0)
}

func (l *Lexer) readNext() (r rune, err error) {
	r, _, err = l.input.ReadRune()
	return
}

func (l *Lexer) skipWhitespace() error {
	for {
		r, err := l.readNext()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(r) {
			break
		}
	}
	l.input.UnreadRune()
	return nil
}
