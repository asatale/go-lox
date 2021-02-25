package tokenizer

import (
  "bytes"
  "errors"
  "io"
  "strconv"
  "unicode"
)

var TokenizeError = errors.New("Error in Tokenizer")

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

Again:
  if t.Scan() {
    rune, _, err := t.source.ReadRune()
    if err != nil {
      return NullToken, TokenizeError
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
      goto Again
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
    case "&", "|":
      nextChar, _, err := t.source.ReadRune()
      if err != nil {
        return NullToken, TokenizeError
      }
      if string(nextChar) == "&" && string(rune) == "&" {
        return Token{
          Type:  _tokenMap["&&"],
          Value: "&&",
          Line:  t.lineNum,
        }, nil
      }
      if string(nextChar) == "|" && string(rune) == "|" {
        return Token{
          Type:  _tokenMap["||"],
          Value: "||",
          Line:  t.lineNum,
        }, nil
      }
      return NullToken, TokenizeError
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
    default:
      t.source.UnreadRune()
      tk, err := t.getComplexToken()
      return tk, err
    }
  }
  return Token{
    Type:  EOF,
    Value: "EOF",
    Line:  t.lineNum,
  }, nil
}

func (t *Tokenizer) getComplexToken() (Token, error) {
  var b bytes.Buffer
  for t.Scan() {
    r, _, err := t.source.ReadRune()
    if err != nil {
      break
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

  return Token{
    Type:  IDENTIFIER,
    Value: b.String(),
    Line:  t.lineNum,
  }, nil
}
