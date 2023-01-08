package tokenizer

import (
	"fmt"
)

// TokenError is error emitted by Tokenizer
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
