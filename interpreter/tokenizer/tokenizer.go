package tokenizer

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
)

// Tokenizer is interface for token generation
type Tokenizer interface {
	GetToken() (Token, error)
}

type tokenizer struct {
	src     *symReader
	lineNum int
}

// NewTokenizer creates new instance of tokenizer
func NewTokenizer(source io.Reader) Tokenizer {
	return &tokenizer{
		src: initSymReader(source),
	}
}

// GetToken returns next token
func (tk *tokenizer) GetToken() (Token, error) {
Loop:
	sym, err := tk.src.peekNext(1)
	if err != nil {
		if err == EOFError {
			return _token{
				_type:  EOF,
				_value: "EOF",
				_line:  tk.lineNum,
			}, nil
		}
		return NullToken, emitError("Unknown IOError", tk.lineNum)
	}

	switch sym {
	case "(", ")", "{", "}", ",", ".", "-", "+", ";", "*":
		tk.src.consumeNext(1)
		return _token{
			_type:  _tokenMap[sym],
			_value: sym,
			_line:  tk.lineNum,
		}, nil
	case "\t", " ", "\n":
		tk.src.consumeNext(1)
		if sym == "\n" {
			tk.lineNum++
		}
		goto Loop
	case "!", "=", ">", "<":
		nextSym, err := tk.src.peekNext(2)
		if err == nil {
			if nextSym != "=" {
				tk.src.consumeNext(1)
				return _token{
					_type:  _tokenMap[sym],
					_value: sym,
					_line:  tk.lineNum,
				}, nil
			}
		} else {
			return NullToken, emitError("Unknown IOError", tk.lineNum)
		}
		tk.src.consumeNext(2)
		return _token{
			_type:  _tokenMap[sym+nextSym],
			_value: sym + nextSym,
			_line:  tk.lineNum,
		}, nil
	case "/":
		nextSym, err := tk.src.peekNext(2)
		if err == nil {
			switch {
			case nextSym == "/":
				tk.src.consumeNext(2)
				return tk.processSingleLineComment()
			case string(nextSym) == "*":
				tk.src.consumeNext(2)
				return tk.processMultiLineComment()
			}
		}
		tk.src.consumeNext(1)
		return _token{
			_type:  _tokenMap[sym],
			_value: sym,
			_line:  tk.lineNum,
		}, nil

	case `"`:
		var strToken string
		tk.src.consumeNext(1)
		for {
			nextSym, err := tk.src.getNext()
			switch {
			case err != nil:
				return NullToken, emitError("Unterminated \"", tk.lineNum)
			case nextSym == `"`:
				return _token{
					_type:  STRING,
					_value: strToken,
					_line:  tk.lineNum,
				}, nil
			case nextSym == "\n":
				tk.lineNum++
			}
			strToken += nextSym
		}
	default:
		tk, err := tk.processComplexToken()
		return tk, err
	}
	// Should not reach here
	return NullToken, emitError("TokenError", tk.lineNum)
}

func isDigit(s string) bool {
	var regEx = regexp.MustCompile(`^[0-9]+`)
	if regEx.MatchString(s) {
		return true
	}
	return false
}

func isAlpha(s string) bool {
	var regEx = regexp.MustCompile(`^[a-zA-Z]+`)
	if regEx.MatchString(s) {
		return true
	}
	return false
}

func (tk *tokenizer) processComplexToken() (Token, error) {
	var strToken string

	for {
		sym, err := tk.src.peekNext(1)
		if err != nil {
			break
		}
		switch {
		case isDigit(sym), isAlpha(sym), sym == "_":
			tk.src.consumeNext(1)
			strToken += sym
			continue
		case sym == ".":
			nextSym, _ := tk.src.peekNext(2)
			if isDigit(nextSym) {
				strToken += (sym + nextSym)
				tk.src.consumeNext(2)
				continue
			}
		}
		break
	}

	if _, err := strconv.ParseInt(strToken, 10, 64); err == nil {
		return _token{
			_type:  INTEGER,
			_value: strToken,
			_line:  tk.lineNum,
		}, nil
	}

	if _, err := strconv.ParseFloat(strToken, 64); err == nil {
		return _token{
			_type:  FLOAT,
			_value: strToken,
			_line:  tk.lineNum,
		}, nil
	}

	if _, ok := _tokenMap[strToken]; ok {
		return _token{
			_type:  _tokenMap[strToken],
			_value: strToken,
			_line:  tk.lineNum,
		}, nil
	}

	var validIdRegEx = regexp.MustCompile(`^[a-zA-Z_]+[a-zA-Z0-9]*$`)
	result := validIdRegEx.MatchString(strToken)

	if result {
		return _token{
			_type:  IDENTIFIER,
			_value: strToken,
			_line:  tk.lineNum,
		}, nil
	}

	return NullToken, emitError(fmt.Sprintf("Invalid token \"%s\"", strToken), tk.lineNum)
}

func (tk *tokenizer) processSingleLineComment() (Token, error) {
	var strToken string
	for {
		sym, err := tk.src.getNext()
		if err != nil || sym == "\n" {
			tk.lineNum++
			return _token{
				_type:  COMMENT,
				_value: strToken,
				_line:  tk.lineNum - 1,
			}, nil
		}
		strToken += sym
	}
	return NullToken, emitError("TokenError", tk.lineNum)
}

func (tk *tokenizer) processMultiLineComment() (Token, error) {
	var strToken string
	lineno := tk.lineNum
	for {
		sym, err := tk.src.getNext()
		switch {
		case err != nil:
			return NullToken, emitError("Unterminated block comment", tk.lineNum)
		case sym == "\n":
			tk.lineNum++
		case sym == "*":
			if nextSym, err := tk.src.peekNext(1); err == nil && nextSym == "/" {
				tk.src.consumeNext(1)
				return _token{
					_type:  COMMENT,
					_value: strToken,
					_line:  lineno,
				}, nil
			}
		}
		strToken += sym
	}
	return NullToken, emitError("Unterminated block comment", tk.lineNum)
}
