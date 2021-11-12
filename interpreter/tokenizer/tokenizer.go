package tokenizer

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"unicode"
)

const (
	SINGLE_LINE_COMMENT = 0 // Single line comment
	MULTI_LINE_COMMENT  = 1 // /*.....*/
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
		if err == nil {
			switch {
			case string(nextChar) == "/":
				return t.singleLineComment()
			case string(nextChar) == "*":
				return t.multiLineComment()
			default:
				t.source.UnreadRune()
			}
		}
		return Token{
			Type:  _tokenMap[string(rune)],
			Value: string(rune),
			Line:  t.lineNum,
		}, nil

	case `"`:
		var b bytes.Buffer
		for {
			nextChar, _, err := t.source.ReadRune()
			switch {
			case err != nil:
				return NullToken, emitError("Unterminated \"", t.lineNum)
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

	for t.source.Len() > 0 {
		r, _, err := t.source.ReadRune()
		if err != nil {
			break
		}

		// End scanning for any space/symbols/punct except for "_"
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

	var validIdRegEx = regexp.MustCompile(`^[a-zA-Z]+_?[a-zA-Z0-9]*$`)
	idString := b.String()

	result := validIdRegEx.MatchString(idString)

	if result {
		return Token{
			Type:  IDENTIFIER,
			Value: idString,
			Line:  t.lineNum,
		}, nil
	}

	return NullToken, emitError(fmt.Sprintf("Invalid identifier \"%s\"", idString), t.lineNum)
}

func (t *tokenizer) singleLineComment() (Token, error) {
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
	return NullToken, emitError("TokenError", t.lineNum)
}

func (t *tokenizer) multiLineComment() (Token, error) {
	var b bytes.Buffer
	lineno := t.lineNum
	for {
		nextChar, _, err := t.source.ReadRune()
		switch {
		case err != nil:
			return NullToken, emitError("Unterminated block comment", t.lineNum)
		case string(nextChar) == "\n":
			t.lineNum++
		case string(nextChar) == "*":
			nextChar, _, err := t.source.ReadRune()
			if err == nil {
				if string(nextChar) == "/" {
					return Token{
						Type:  COMMENT,
						Value: b.String(),
						Line:  lineno,
					}, nil
				} else {
					t.source.UnreadRune()
				}
			}
		}
		b.WriteRune(nextChar)
	}
	return NullToken, emitError("Unterminated block comment", t.lineNum)
}
