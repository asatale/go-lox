package interpreter

import (
  "fmt"
  "github.com/asatale/go-lox/interpreter/tokenizer"
  "io"
)

// Run is top level exec routine
func Run(source io.Reader) error {

  tk := tokenizer.NewTokenizer(source)

  for tk.Scan() {
    token, err := tk.GetToken()
    switch err {
    case nil:
      fmt.Println("Token: ", token)
    case tokenizer.SourceEmpty:
      break
    default:
      fmt.Println("Error: ", err)
      return err
    }
  }
  return nil
}
