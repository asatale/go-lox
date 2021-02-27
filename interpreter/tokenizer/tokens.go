package tokenizer

import (
	"fmt"
)

type TokenType int

const (
	NULLTOKEN    TokenType = iota // No Token
	LEFTPAREN                     // "("
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
	NUMBER                        // E.g. 42, 3.14
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
	EOF                           // "EOF"
	COMMENT                       // "// This is a comment"
)

func (t TokenType) String() string {
	switch t {
	case NULLTOKEN:
		return "Null"
	case LEFTPAREN:
		return "("
	case RIGHTPAREN:
		return ")"
	case LEFTBRACE:
		return "{"
	case RIGHTBRACE:
		return "}"
	case COMMA:
		return ","
	case DOT:
		return "."
	case MINUS:
		return "-"
	case PLUS:
		return "+"
	case SEMICOLON:
		return ";"
	case DIVIDE:
		return "/"
	case MULTIPLY:
		return "*"
	case BANG:
		return "!"
	case BANGEQUAL:
		return "!="
	case EQUAL:
		return "="
	case DOUBLEEQUAL:
		return "=="
	case GREATER:
		return ">"
	case GREATEREQUAL:
		return ">="
	case LESS:
		return "<"
	case LESSEQUAL:
		return "<="
	case IDENTIFIER:
		return "identifier"
	case STRING:
		return "string"
	case NUMBER:
		return "number"
	case AND:
		return "and"
	case OR:
		return "or"
	case CLASS:
		return "class"
	case ELSE:
		return "else"
	case FALSE:
		return "false"
	case FUN:
		return "fun"
	case FOR:
		return "for"
	case IF:
		return "if"
	case NIL:
		return "nil"
	case PRINT:
		return "print"
	case RETURN:
		return "return"
	case SUPER:
		return "super"
	case THIS:
		return "this"
	case TRUE:
		return "true"
	case VAR:
		return "var"
	case WHILE:
		return "while"
	case EOF:
		return "EOF"
	case COMMENT:
		return "Comment"
	}
	return "Unknown Token"
}

type Token struct {
	Type  TokenType
	Value string
	Line  int
}

var NullToken = Token{
	Type: NULLTOKEN,
}

func (t Token) String() string {
	return fmt.Sprintf("Token{ Type:%v, Value: %v, Line: %v}", t.Type, t.Value, t.Line)
}

var _tokenMap = map[string]TokenType{
	"(":      LEFTPAREN,
	")":      RIGHTPAREN,
	"{":      LEFTBRACE,
	"}":      RIGHTBRACE,
	",":      COMMA,
	".":      DOT,
	"-":      MINUS,
	"+":      PLUS,
	";":      SEMICOLON,
	"/":      DIVIDE,
	"*":      MULTIPLY,
	"!":      BANG,
	"=":      EQUAL,
	">":      GREATER,
	"<":      LESS,
	"!=":     BANGEQUAL,
	"==":     DOUBLEEQUAL,
	">=":     GREATEREQUAL,
	"<=":     LESSEQUAL,
	"and":    AND,
	"or":     OR,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"fun":    FUN,
	"for":    FOR,
	"if":     IF,
	"nil":    NIL,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}
