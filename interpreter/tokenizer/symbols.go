package tokenizer

import (
	"bytes"
	"errors"
	"io"
)

type symReader struct {
	source    *bytes.Buffer
	lookahead []string
	index     int
	lineNum   int
}

var (
	EOFError   = errors.New("EOF")
	UnknownErr = errors.New("Unknown error")
	ArgErr     = errors.New("Invalid argument")
)

func initSymReader(s io.Reader) *symReader {
	buf := bytes.NewBuffer([]byte{})
	buf.ReadFrom(s)
	return &symReader{
		source: buf,
	}
}

func updateLookahead(t *symReader) error {
	if t.source.Len() == 0 {
		return EOFError
	}
	r, _, err := t.source.ReadRune()
	if err != nil {
		return UnknownErr
	}
	if string(r) != "" {
		t.lookahead = append(t.lookahead, string(r))
	} else {
		return EOFError
	}
	return nil
}

func (t *symReader) getNext() (string, error) {
	if len(t.lookahead) == 0 {
		err := updateLookahead(t)
		if err != nil {
			return "", err
		}
	}
	sym := t.lookahead[0]
	t.lookahead = t.lookahead[1:]
	return sym, nil
}

func (t *symReader) peekNext(ahead int) (string, error) {
	if ahead <= 0 {
		return "", ArgErr
	}
	if len(t.lookahead) < ahead {
		err := updateLookahead(t)
		if err != nil {
			return "", err
		}
	}
	return t.lookahead[ahead-1], nil
}

func (t *symReader) consumeNext(ahead int) error {
	if ahead > len(t.lookahead) {
		return ArgErr
	}
	t.lookahead = t.lookahead[ahead:]
	return nil
}
