package lexer

import (
    "testing"
    "golox/token"
    "fmt"
)

type Expectations struct {
    expectedType token.TokenType
    expectedLexeme string
    expectedLiteral interface{}
}

func testTokens(t *testing.T, input string, tests []Expectations) {
    l := NewLexer(input)
    for i, tt := range tests {
        tok := l.ScanToken()
        if tok.Type != tt.expectedType {
            t.Fatalf("tests[%d] - tokentype wrong. expected=%s, got=%s",
                i, fmt.Sprint(tt.expectedType), fmt.Sprint(tok.Type))
        }
        if tok.Lexeme != tt.expectedLexeme {
            t.Fatalf("tests[%d] - lexeme wrong. expected=\"%s\", got=\"%s\"",
                i, tt.expectedLexeme, tok.Lexeme)
        }
        if tok.Literal != tt.expectedLiteral {
            t.Fatalf("tests[%d] - literal wrong. expected=%s, got=%s",
                i, fmt.Sprint(tt.expectedLiteral), fmt.Sprint(tok.Literal))
        }
    }
}

func TestSingleSymbols(t *testing.T){
    input := `=+(){},;.-+/*!`
    tests := []Expectations{
        {token.EQUAL, "=", nil},
        {token.PLUS, "+", nil},
        {token.LEFT_PAREN, "(", nil},
        {token.RIGHT_PAREN, ")", nil},
        {token.LEFT_BRACE, "{", nil},
        {token.RIGHT_BRACE, "}", nil},
        {token.COMMA, ",", nil},
        {token.SEMICOLON, ";", nil},
        {token.DOT, ".", nil},
        {token.MINUS, "-", nil},
        {token.PLUS, "+", nil},
        {token.SLASH, "/", nil},
        {token.STAR, "*", nil},
        {token.BANG, "!", nil},
        {token.EOF, "", nil},
    }
    testTokens(t, input, tests)
}

func TestString(t *testing.T){
    str1 := "12345a bcdef g*&24"
    quotedStr1 := fmt.Sprintf("\"%s\"", str1)
    input := fmt.Sprintf(" %s ", quotedStr1)
    tests := []Expectations{
        {token.STRING, quotedStr1, str1},
        {token.EOF, "", nil},
    }
    testTokens(t, input, tests)
}

func TestNumber(t *testing.T){
    input := "1.3413 2 3 6 12417.1"
    tests := []Expectations{
        {token.NUMBER, "1.3413", 1.3413},
        {token.NUMBER, "2", 2.0},
        {token.NUMBER, "3", 3.0},
        {token.NUMBER, "6", 6.0},
        {token.NUMBER, "12417.1", 12417.1},
        {token.EOF, "", nil},
    }
    testTokens(t, input, tests)
}

func TestDoubleSymbols(t *testing.T){
    input := ` == != >= <= `
    tests := []Expectations{
        {token.EQUAL_EQUAL, "==", nil},
        {token.BANG_EQUAL, "!=", nil},
        {token.LESS_EQUAL, ">=", nil},
        {token.GREATER_EQUAL, "<=", nil},
        {token.EOF, "", nil},
    }
    testTokens(t, input, tests)
}

func TestSymbolsWithWhitespaces(t *testing.T){
    input := "\n = + \n()   {} \t, \r;  \n"
    tests := []Expectations{
        {token.EQUAL, "=", nil},
        {token.PLUS, "+", nil},
        {token.LEFT_PAREN, "(", nil},
        {token.RIGHT_PAREN, ")", nil},
        {token.LEFT_BRACE, "{", nil},
        {token.RIGHT_BRACE, "}", nil},
        {token.COMMA, ",", nil},
        {token.SEMICOLON, ";", nil},
        {token.EOF, "", nil},
    }
    testTokens(t, input, tests)
}

func TestKeywords(t *testing.T){
    input := `and class else false for fun
            if   nil   or print return
            super this true var while blahblah`
    tests := []Expectations{
        {token.AND, "and", nil},
        {token.CLASS, "class", nil},
        {token.ELSE, "else", nil},
        {token.FALSE, "false", nil},
        {token.FOR, "for", nil},
        {token.FUN, "fun", nil},
        {token.IF, "if", nil},
        {token.NIL, "nil", nil},
        {token.OR, "or", nil},
        {token.PRINT, "print", nil},
        {token.RETURN, "return", nil},
        {token.SUPER, "super", nil},
        {token.THIS, "this", nil},
        {token.TRUE, "true", nil},
        {token.VAR, "var", nil},
        {token.WHILE, "while", nil},
        {token.IDENTIFIER, "blahblah", nil},
        {token.EOF, "", nil},
    }
    testTokens(t, input, tests)
}
