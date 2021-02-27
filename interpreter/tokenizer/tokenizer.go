package tokenizer

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

// TokenError is error
type TokenError struct {
	msg  string
	line int
}

func (e TokenError) Error() string {
	return fmt.Sprintf("%s at %d", e.msg, e.line)
}

func emitError(s string, l int) error {
	return &TokenError{msg: s, line: l}
}

type tokenizer struct {
	source  *bytes.Buffer
	index   int
	lineNum int
}

// Tokenizer is interface for token generation
type Tokenizer interface {
	GetToken() (Token, error)
}

// NewTokenizer creates new instance of tokenizer
func NewTokenizer(source io.Reader) Tokenizer {
	buf := bytes.NewBuffer([]byte{})
	buf.ReadFrom(source)
	return &tokenizer{
		source: buf,
	}
}

// GetToken returns next token
func (t *tokenizer) GetToken() (Token, error) {
Loop:
	rune, _, err := t.source.ReadRune()
	if err != nil {
		if err == io.EOF {
			return Token{
				Type:  EOF,
				Value: "EOF",
				Line:  t.lineNum,
			}, nil
		}
		return NullToken, emitError("IOError", t.lineNum)
	}

	switch string(rune) {
	case "(", ")", "{", "}", ",", ".", "-", "+", ";", "*":
		return Token{
			Type:  _tokenMap[string(rune)],
			Value: string(rune),
			Line:  t.lineNum,
		}, nil
	case "\t", " ", "\n":
		if string(rune) == "\n" {
			t.lineNum++
		}
		goto Loop
	case "!", "=", ">", "<":
		var nextChar int32
		nextChar, _, err = t.source.ReadRune()
		if err != nil || string(nextChar) != "=" {
			if err == nil {
				t.source.UnreadRune()
			}
			return Token{
				Type:  _tokenMap[string(rune)],
				Value: string(rune),
				Line:  t.lineNum,
			}, nil
		}
		return Token{
			Type:  _tokenMap[string(rune)+string(nextChar)],
			Value: string(rune) + string(nextChar),
			Line:  t.lineNum,
		}, nil
	case "/":
		nextChar, _, err := t.source.ReadRune()
		if err != nil || string(nextChar) != "/" {
			t.source.UnreadRune()
			return Token{
				Type:  _tokenMap[string(rune)],
				Value: string(rune),
				Line:  t.lineNum,
			}, nil
		}
		var b bytes.Buffer
		for {
			nextChar, _, err := t.source.ReadRune()
			if err != nil || string(nextChar) == "\n" {
				t.lineNum++
				return Token{
					Type:  COMMENT,
					Value: b.String(),
					Line:  t.lineNum - 1,
				}, nil
			}
			b.WriteRune(nextChar)
		}
	case `"`:
		var b bytes.Buffer
		for {
			nextChar, _, err := t.source.ReadRune()
			switch {
			case err != nil:
				return NullToken, emitError("TokenError", t.lineNum)
			case string(nextChar) == `"`:
				return Token{
					Type:  STRING,
					Value: b.String(),
					Line:  t.lineNum,
				}, nil
			case string(nextChar) == "\n":
				t.lineNum++
			}
			b.WriteRune(nextChar)
		}
	default:
		t.source.UnreadRune()
		tk, err := t.getComplexToken()
		return tk, err
	}
	// Should not reach here
	return NullToken, emitError("TokenError", t.lineNum)
}

func (t *tokenizer) getComplexToken() (Token, error) {
	var b bytes.Buffer
	firstChar := unicode.MaxRune + 1
	for t.source.Len() > 0 {
		r, _, err := t.source.ReadRune()
		if err != nil {
			break
		}
		if firstChar == (unicode.MaxRune + 1) { // Store first char for valid identifier check
			firstChar = r
		}
		if unicode.IsSpace(r) || unicode.IsSymbol(r) || unicode.IsPunct(r) {
			if string(r) != "_" {
				t.source.UnreadRune()
				break
			}
		}
		b.WriteRune(r)
	}

	if _, ok := _tokenMap[b.String()]; ok {
		return Token{
			Type:  _tokenMap[b.String()],
			Value: b.String(),
			Line:  t.lineNum,
		}, nil
	}

	if _, err := strconv.ParseFloat(b.String(), 64); err == nil {
		return Token{
			Type:  NUMBER,
			Value: b.String(),
			Line:  t.lineNum,
		}, nil
	}
	if _, err := strconv.ParseInt(b.String(), 10, 64); err == nil {
		return Token{
			Type:  NUMBER,
			Value: b.String(),
			Line:  t.lineNum,
		}, nil
	}

	if !unicode.IsLetter(firstChar) {
		return NullToken, emitError("Error Token", t.lineNum)
	}
	return Token{
		Type:  IDENTIFIER,
		Value: b.String(),
		Line:  t.lineNum,
	}, nil
}
