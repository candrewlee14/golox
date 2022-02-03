package token

import (
	"fmt"
)

//go:generate stringer -type=TokenType
type TokenType uint8

const (

	// Single-char tokens
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two char tokens
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals
	IDENTIFIER
	STRING
	NUMBER

	// Keywords
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
	INVALID
)

type Token struct {
	Type       TokenType
	Lexeme     string
	Line       int
	LineOffset int
	Literal    interface{}
}

func NewToken(toktype TokenType, lexeme string, line int, lineOffset int, literal interface{}) Token {
	return Token{toktype, lexeme, line, lineOffset, literal}
}

func (tok Token) String() string {
	return fmt.Sprintf("%s \"%s\" {%s} %d:%d",
		fmt.Sprint(tok.Type),
		tok.Lexeme,
		fmt.Sprint(tok.Literal),
		tok.Line,
		tok.LineOffset)
}
