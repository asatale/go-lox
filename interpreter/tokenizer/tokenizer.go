package tokenizer

import (
  "bytes"
  "errors"
  "io"
)

var TokenizeError = errors.New("Error in Tokenizer")
var SourceEmpty = errors.New("No more source left to tokenize")

type Tokenizer struct {
  source  *bytes.Buffer
  index   int
  lineNum int
}

// NewTokenizer creates new instance of tokenizer
func NewTokenizer(source io.Reader) *Tokenizer {
  buf := bytes.NewBuffer([]byte{})
  buf.ReadFrom(source)
  return &Tokenizer{
    source: buf,
  }
}

// Scan returns bool to indicate if more tokens are available
func (t *Tokenizer) Scan() bool {
  if t.source.Len() > 0 {
    return true
  }
  return false
}

// GetToken returns next token
func (t *Tokenizer) GetToken() (Token, error) {

  for t.Scan() {
    rune, _, err := t.source.ReadRune()
    if err != nil {
      return NullToken, TokenizeError
    }

    switch string(rune) {
    case "(", ")", "{", "}", ",", ".", "-", "+", ";", "*":
      return Token{
        Type:  CharTokensMap[string(rune)],
        Value: string(rune),
        Line:  t.lineNum,
      }, nil
    case "\t", " ", "\n":
      if string(rune) == "\n" {
        t.lineNum += 1
      }
      continue
    case "!", "=", ">", "<":
      var nextChar int32
      for {
        nextChar, _, err = t.source.ReadRune()
        if err == nil && string(nextChar) != " " {
          break
        }
      }
      if string(nextChar) == "=" {
        return Token{
          Type:  CharTokensMap[string(rune)+string(nextChar)],
          Value: string(rune) + string(nextChar),
          Line:  t.lineNum,
        }, nil
      }
      t.source.UnreadRune()
      return Token{
        Type:  CharTokensMap[string(rune)],
        Value: string(rune),
        Line:  t.lineNum,
      }, nil
    case "&", "|":
      nextChar, _, err := t.source.ReadRune()
      if err != nil {
        return NullToken, TokenizeError
      }
      if string(nextChar) == "&" && string(rune) == "&" {
        return Token{
          Type:  CharTokensMap["&&"],
          Value: "&&",
          Line:  t.lineNum,
        }, nil
      }
      if string(nextChar) == "|" && string(rune) == "|" {
        return Token{
          Type:  CharTokensMap["||"],
          Value: "||",
          Line:  t.lineNum,
        }, nil
      }
      return NullToken, TokenizeError
    case "/":
      nextChar, _, err := t.source.ReadRune()
      if err != nil {
        return NullToken, TokenizeError
      }
      if string(nextChar) != "/" {
        t.source.UnreadRune()
        return Token{
          Type:  CharTokensMap[string(rune)],
          Value: string(rune),
          Line:  t.lineNum,
        }, nil
      }
      for t.Scan() {
        nextChar, _, err := t.source.ReadRune()
        if err != nil {
          return NullToken, TokenizeError
        }
        if string(nextChar) == "\n" {
          t.lineNum += 1
          break
        }
      }

    }

  }
  return NullToken, SourceEmpty
}
