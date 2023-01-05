package interpreter

import (
	"errors"
	"fmt"
	"github.com/asatale/go-lox/interpreter/tokenizer"
	"io"
)

// Eof indicates end of file
var EOFError = errors.New("End of data")

// Run is top level exec routine
func Run(source io.Reader) error {

	tk := tokenizer.NewTokenizer(source)

	for {
		token, err := tk.GetToken()
		if err != nil {
			fmt.Println("Error: ", err)
			return err
		}
		if token.Type() == tokenizer.EOF {
			break
		}
		fmt.Println("Token: ", token)
	}
	return EOFError
}
