package tokenizer

import (
	"bytes"
	"testing"
)

type testCase struct {
	description   string
	source        *bytes.Buffer
	result        []TokenType
	errorExpected bool
}

func runTestcases(testCases []testCase, t *testing.T) {
	for _, testCase := range testCases {
		tk := NewTokenizer(testCase.source)
		for _, tokenType := range testCase.result {
			token, err := tk.GetToken()
			switch {
			case err == nil && testCase.errorExpected:
				t.Errorf("%v: Error expected. Got nil", testCase.description)
			case err == nil && token.Type != tokenType:
				t.Errorf("%v: Token expected: %v, Got: %v ", testCase.description, tokenType, token.Type)
			case err != nil && !testCase.errorExpected:
				t.Errorf("%v: Error expected: nil, Got: %v", testCase.description, err)
			}
		}
	}
}

func TestSingleTokensTests(t *testing.T) {
	testCases := []testCase{
		{
			description:   "Test for symbols ",
			source:        bytes.NewBufferString(`( ) { } , . - + ; / * ! != = == > >= < <= `),
			result:        []TokenType{LEFTPAREN, RIGHTPAREN, LEFTBRACE, RIGHTBRACE, COMMA, DOT, MINUS, PLUS, SEMICOLON, DIVIDE, MULTIPLY, BANG, BANGEQUAL, EQUAL, DOUBLEEQUAL, GREATER, GREATEREQUAL, LESS, LESSEQUAL},
			errorExpected: false,
		},
		{
			description:   "Test for Identifier",
			source:        bytes.NewBufferString("i count Count countMin test_iteration i32"),
			result:        []TokenType{IDENTIFIER, IDENTIFIER, IDENTIFIER, IDENTIFIER, IDENTIFIER, IDENTIFIER},
			errorExpected: false,
		},
		{
			description:   `Test for token string "Hello World"  "Good Bye"`,
			source:        bytes.NewBufferString(`"Hello World"`),
			result:        []TokenType{STRING},
			errorExpected: false,
		},
		{
			description:   "Test for token integer number",
			source:        bytes.NewBufferString("42 3.14"),
			result:        []TokenType{NUMBER, NUMBER},
			errorExpected: false,
		},
		{
			description:   "Test for reserved keywords",
			source:        bytes.NewBufferString("and or class else false fun for if nil print return super this true var while"),
			result:        []TokenType{AND, OR, CLASS, ELSE, FALSE, FUN, FOR, IF, NIL, PRINT, RETURN, SUPER, THIS, TRUE, VAR, WHILE},
			errorExpected: false,
		},
		{
			description:   "Test for empty source ",
			source:        bytes.NewBufferString(``),
			result:        []TokenType{EOF},
			errorExpected: false,
		},
		{
			description:   "Test for Line Comment ",
			source:        bytes.NewBufferString(`// This is a comment`),
			result:        []TokenType{COMMENT},
			errorExpected: false,
		},
	}
	runTestcases(testCases, t)
}

func TestMixedTokensTests(t *testing.T) {
	testCases := []testCase{
		{
			description: "Test print statements",
			source: bytes.NewBufferString(`
        print "Hello, world!";
        print breakfast; // "bagels".
        print 2+3;
      `),
			result: []TokenType{
				PRINT, STRING, SEMICOLON,
				PRINT, IDENTIFIER, SEMICOLON, COMMENT,
				PRINT, NUMBER, PLUS, NUMBER, SEMICOLON,
			},
			errorExpected: false,
		},
		{
			description: "Testing operations + - * /",
			source: bytes.NewBufferString(`
         add + me;
         subtract - me;
         multiply * me;
         divide / me;
      `),
			result: []TokenType{
				IDENTIFIER, PLUS, IDENTIFIER, SEMICOLON,
				IDENTIFIER, MINUS, IDENTIFIER, SEMICOLON,
				IDENTIFIER, MULTIPLY, IDENTIFIER, SEMICOLON,
				IDENTIFIER, DIVIDE, IDENTIFIER, SEMICOLON},
			errorExpected: false,
		},
		{
			description: "Testing unary operator ",
			source: bytes.NewBufferString(`
         -negateMe;
      `),
			result: []TokenType{MINUS, IDENTIFIER, SEMICOLON},
		},
		{
			description: "Testing comparison operators",
			source: bytes.NewBufferString(`
        less < than;
        lessThan <= orEqual;
        greater > than;
        greaterThan >= orEqual;
        1 == 2;         // false.
        "cat" != "dog"; // true.
      `),
			result: []TokenType{
				IDENTIFIER, LESS, IDENTIFIER, SEMICOLON,
				IDENTIFIER, LESSEQUAL, IDENTIFIER, SEMICOLON,
				IDENTIFIER, GREATER, IDENTIFIER, SEMICOLON,
				IDENTIFIER, GREATEREQUAL, IDENTIFIER, SEMICOLON,
				NUMBER, DOUBLEEQUAL, NUMBER, SEMICOLON, COMMENT,
				STRING, BANGEQUAL, STRING, SEMICOLON, COMMENT},
			errorExpected: false,
		},
		{
			description: "Testing logical operators ",
			source: bytes.NewBufferString(`
        !true;  // false.
        !false; // true.
        true and false; // false.
        true and true;  // true.
        false or false; // false.
        true or false;  // true.
     `),
			result: []TokenType{
				BANG, TRUE, SEMICOLON, COMMENT,
				BANG, FALSE, SEMICOLON, COMMENT,
				TRUE, AND, FALSE, SEMICOLON, COMMENT,
				TRUE, AND, TRUE, SEMICOLON, COMMENT,
				FALSE, OR, FALSE, SEMICOLON, COMMENT,
				TRUE, OR, FALSE, SEMICOLON, COMMENT,
			},
			errorExpected: false,
		},
		{
			description: "Statements Blocks",
			source: bytes.NewBufferString(`
        "some expression";
        {
           print "One statement.";
           print "Two statements.";
        }
     `),
			result: []TokenType{
				STRING, SEMICOLON,
				LEFTBRACE,
				PRINT, STRING, SEMICOLON,
				PRINT, STRING, SEMICOLON,
				RIGHTBRACE,
			},
			errorExpected: false,
		},
		{
			description: "Variables declaration",
			source: bytes.NewBufferString(`
        var imAVariable = "here is my value";
        var iAmNil;
        var breakfast = "bagels";
        var average = (min + max) / 2;
     `),
			result: []TokenType{
				VAR, IDENTIFIER, EQUAL, STRING, SEMICOLON,
				VAR, IDENTIFIER, SEMICOLON,
				VAR, IDENTIFIER, EQUAL, STRING, SEMICOLON,
				VAR, IDENTIFIER, EQUAL, LEFTPAREN, IDENTIFIER, PLUS, IDENTIFIER, RIGHTPAREN, DIVIDE, NUMBER, SEMICOLON,
			},
			errorExpected: false,
		},
		{
			description: "Control",
			source: bytes.NewBufferString(`
        if (condition) {
          print "yes";
        } else {
          print "no";
        }

        var a = 1;
        while (a < 10) {
          print a;
          a = a + 1;
        }

        for (var a = 1; a < 10; a = a + 1) {
          print a;
        }
     `),
			result: []TokenType{
				IF, LEFTPAREN, IDENTIFIER, RIGHTPAREN, LEFTBRACE,
				PRINT, STRING, SEMICOLON,
				RIGHTBRACE, ELSE, LEFTBRACE,
				PRINT, STRING, SEMICOLON,
				RIGHTBRACE,

				VAR, IDENTIFIER, EQUAL, NUMBER, SEMICOLON,
				WHILE, LEFTPAREN, IDENTIFIER, LESS, NUMBER, RIGHTPAREN, LEFTBRACE,
				PRINT, IDENTIFIER, SEMICOLON,
				IDENTIFIER, EQUAL, IDENTIFIER, PLUS, NUMBER, SEMICOLON,
				RIGHTBRACE,

				FOR, LEFTPAREN, VAR, IDENTIFIER, EQUAL, NUMBER, SEMICOLON, IDENTIFIER, LESS, NUMBER, SEMICOLON, IDENTIFIER, EQUAL, IDENTIFIER, PLUS, NUMBER, RIGHTPAREN, LEFTBRACE,
				PRINT, IDENTIFIER, SEMICOLON,
				RIGHTBRACE,
			},
			errorExpected: false,
		},
		{
			description: "Functions",
			source: bytes.NewBufferString(`
        makeBreakfast(bacon, eggs, toast);
        makeBreakfast();

        fun printSum(a, b) {
          print a + b;
        }

        fun returnSum(a, b) {
          return a + b;
        }

     `),
			result: []TokenType{
				IDENTIFIER, LEFTPAREN, IDENTIFIER, COMMA, IDENTIFIER, COMMA, IDENTIFIER, RIGHTPAREN, SEMICOLON,
				IDENTIFIER, LEFTPAREN, RIGHTPAREN, SEMICOLON,

				FUN, IDENTIFIER, LEFTPAREN, IDENTIFIER, COMMA, IDENTIFIER, RIGHTPAREN, LEFTBRACE,
				PRINT, IDENTIFIER, PLUS, IDENTIFIER, SEMICOLON,
				RIGHTBRACE,

				FUN, IDENTIFIER, LEFTPAREN, IDENTIFIER, COMMA, IDENTIFIER, RIGHTPAREN, LEFTBRACE,
				RETURN, IDENTIFIER, PLUS, IDENTIFIER, SEMICOLON,
				RIGHTBRACE,
			},
			errorExpected: false,
		},
		{
			description: "Functions",
			source: bytes.NewBufferString(`
      class Breakfast {
          cook() {
            print "Eggs a-fryin'!";
          }

          serve(who) {
            print "Enjoy your breakfast, " + who + ".";
          }
      }
     `),
			result: []TokenType{
				CLASS, IDENTIFIER, LEFTBRACE,
				IDENTIFIER, LEFTPAREN, RIGHTPAREN, LEFTBRACE,
				PRINT, STRING, SEMICOLON,
				RIGHTBRACE,

				IDENTIFIER, LEFTPAREN, IDENTIFIER, RIGHTPAREN, LEFTBRACE,
				PRINT, STRING, PLUS, IDENTIFIER, PLUS, STRING, SEMICOLON,
				RIGHTBRACE,

				RIGHTBRACE,
			},
			errorExpected: false,
		},
	}
	runTestcases(testCases, t)
}

func TestInvalidTokensTests(t *testing.T) {
	testCases := []testCase{
		{
			description: "Test unterminated string",
			source: bytes.NewBufferString(`
        "Hello, world!
      `),
			result: []TokenType{
				NULLTOKEN,
			},
			errorExpected: true,
		},
		{
			description: "Test invalid identifiers",
			source: bytes.NewBufferString(`
        2abc; // Identifier can not start with number
        a-bc; // Identifiers can not contain "-"
      `),
			result: []TokenType{
				NULLTOKEN,
			},
			errorExpected: true,
		},
	}
	runTestcases(testCases, t)
}
