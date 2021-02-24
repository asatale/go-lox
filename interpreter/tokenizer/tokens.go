package tokenizer

import (
  "fmt"
)

type TokenType int

const (
  LEFTPAREN    TokenType = iota // "("
  RIGHTPAREN                    // ")"
  LEFTBRACE                     // "{"
  RIGHTBRACE                    // "}"
  COMMA                         // ","
  DOT                           // "."
  MINUS                         // "-"
  PLUS                          // "+"
  SEMICOLON                     // ";"
  DIVIDE                        // "/"
  MULTIPLY                      // "*"
  BANG                          // "!"
  BANGEQUAL                     // "!="
  EQUAL                         // "="
  DOUBLEEQUAL                   // "=="
  GREATER                       // ">"
  GREATEREQUAL                  // ">="
  LESS                          // "<"
  LESSEQUAL                     // "<="
  IDENTIFIER                    // E.g. "i"
  STRING                        // E.g. "Hello"
  NUMBER                        // E.g. 2
  AND                           // &&
  OR                            // "||"
  CLASS                         // "class"
  ELSE                          // "else"
  FALSE                         // "false"
  FUN                           // "fun"
  FOR                           // "for"
  IF                            // "if"
  NIL                           // "nil"
  PRINT                         // "print"
  RETURN                        // "return"
  SUPER                         // "super"
  THIS                          // "this"
  TRUE                          // "true"
  VAR                           // "var"
  WHILE                         // "while"
  EOF                           // EOF
)

type Token struct {
  Type  TokenType
  Value string
  Line  int
}

var NullToken = Token{}

func (t Token) String() string {
  return fmt.Sprintf("Token<%v, %v, %v>", t.Type, t.Value, t.Line)
}

var CharTokensMap = map[string]TokenType{
  "(":  LEFTPAREN,
  ")":  RIGHTPAREN,
  "{":  LEFTBRACE,
  "}":  RIGHTBRACE,
  ",":  COMMA,
  ".":  DOT,
  "-":  MINUS,
  "+":  PLUS,
  ";":  SEMICOLON,
  "/":  DIVIDE,
  "*":  MULTIPLY,
  "!":  BANG,
  "=":  EQUAL,
  ">":  GREATER,
  "<":  LESS,
  "!=": BANGEQUAL,
  "==": DOUBLEEQUAL,
  ">=": GREATEREQUAL,
  "<=": LESSEQUAL,
  "&&": AND,
  "||": OR,
}
