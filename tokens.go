package main

import (
	"fmt"
	//"text/scanner"
)

type TokenType int

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
)

//go:generate stringer -type=TokenType

type Token struct {
	toktype    TokenType
	lexeme     string
	line       int
	lineOffset int
}

func (tok Token) String() string {
	return fmt.Sprintf("%s %s %d:%d", fmt.Sprint(tok.toktype), tok.lexeme, tok.line, tok.lineOffset)
}
