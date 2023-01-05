package tokenizer

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
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
			return _token{
				_type:  EOF,
				_value: "EOF",
				_line:  t.lineNum,
			}, nil
		}
		return NullToken, emitError("Unknown IOError", t.lineNum)
	}

	switch string(rune) {
	case "(", ")", "{", "}", ",", ".", "-", "+", ";", "*":
		return _token{
			_type:  _tokenMap[string(rune)],
			_value: string(rune),
			_line:  t.lineNum,
		}, nil
	case "\t", " ", "\n":
		if string(rune) == "\n" {
			t.lineNum++
		}
		goto Loop
	case "!", "=", ">", "<":
		nextChar, _, err := t.source.ReadRune()
		if err != nil || string(nextChar) != "=" {
			if err == nil {
				t.source.UnreadRune()
			}
			return _token{
				_type:  _tokenMap[string(rune)],
				_value: string(rune),
				_line:  t.lineNum,
			}, nil
		}
		return _token{
			_type:  _tokenMap[string(rune)+string(nextChar)],
			_value: string(rune) + string(nextChar),
			_line:  t.lineNum,
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
		return _token{
			_type:  _tokenMap[string(rune)],
			_value: string(rune),
			_line:  t.lineNum,
		}, nil

	case `"`:
		var b bytes.Buffer
		for {
			nextChar, _, err := t.source.ReadRune()
			switch {
			case err != nil:
				return NullToken, emitError("Unterminated \"", t.lineNum)
			case string(nextChar) == `"`:
				return _token{
					_type:  STRING,
					_value: b.String(),
					_line:  t.lineNum,
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
		if unicode.IsDigit(r) || unicode.IsLetter(r) || string(r) == "_" || string(r) == "." {
			b.WriteRune(r)
		} else {
			t.source.UnreadRune()
			break
		}
	}

	if _, ok := _tokenMap[b.String()]; ok {
		return _token{
			_type:  _tokenMap[b.String()],
			_value: b.String(),
			_line:  t.lineNum,
		}, nil
	}

	fmt.Println("Checking for numbers: ", b.String())
	if strings.ContainsAny(b.String(), ".") {
		fmt.Println("Parsing floating numbers")
		if _, err := strconv.ParseFloat(b.String(), 64); err == nil {
			return _token{
				_type:  FLOAT,
				_value: b.String(),
				_line:  t.lineNum,
			}, nil
		}
	} else if _, err := strconv.ParseInt(b.String(), 10, 64); err == nil {
		fmt.Println("Parsing integer numbers")
		return _token{
			_type:  INTEGER,
			_value: b.String(),
			_line:  t.lineNum,
		}, nil
	}

	var validIdRegEx = regexp.MustCompile(`^[a-zA-Z_]+[a-zA-Z0-9]*$`)
	result := validIdRegEx.MatchString(b.String())

	if result {
		return _token{
			_type:  IDENTIFIER,
			_value: b.String(),
			_line:  t.lineNum,
		}, nil
	}

	return NullToken, emitError(fmt.Sprintf("Invalid identifier \"%s\"", b.String()), t.lineNum)
}

func (t *tokenizer) singleLineComment() (Token, error) {
	var b bytes.Buffer
	for {
		nextChar, _, err := t.source.ReadRune()
		if err != nil || string(nextChar) == "\n" {
			t.lineNum++
			return _token{
				_type:  COMMENT,
				_value: b.String(),
				_line:  t.lineNum - 1,
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
					return _token{
						_type:  COMMENT,
						_value: b.String(),
						_line:  lineno,
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
