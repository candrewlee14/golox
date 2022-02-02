package main

import (
	"fmt"
)

type Scanner struct {
	source     string
	tokens     []Token
	lexStart   int
	current    int
	lineOffset int
	line       int
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// Consumes and returns next char
func (s *Scanner) advance() byte {
	c := s.source[s.current]
	s.current += 1
	s.lineOffset += 1
	return c
}

// Returns next char without consuming
func (s *Scanner) peek() byte {
	return s.source[s.current]
}

// Checks if next char matches given char and advances if so
func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	c := s.source[s.current]
	if c != expected {
		return false
	}
	s.advance()
	return true
}

func (s *Scanner) addToken(toktype TokenType) {
	lex := s.source[s.lexStart:s.current]
	s.tokens = append(s.tokens, Token{toktype, lex, s.line, s.lineOffset})
}

func (s *Scanner) ScanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		var toktype TokenType
		if s.match('=') {
			toktype = BANG_EQUAL
		} else {
			toktype = BANG
		}
		s.addToken(toktype)
	case '=':
		var toktype TokenType
		if s.match('=') {
			toktype = EQUAL_EQUAL
		} else {
			toktype = EQUAL
		}
		s.addToken(toktype)
	case '<':
		var toktype TokenType
		if s.match('=') {
			toktype = GREATER_EQUAL
		} else {
			toktype = GREATER
		}
		s.addToken(toktype)
	case '>':
		var toktype TokenType
		if s.match('=') {
			toktype = LESS_EQUAL
		} else {
			toktype = EQUAL
		}
		s.addToken(toktype)
	case '/':
		if s.match('/') {
			for !s.isAtEnd() && s.peek() != '\n' {
				s.current += 1
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ', '\t', '\r':
		// Do nothing here, because they aren't unexpected characters
		s.current = s.current
	case '\n':
		s.line += 1
		s.lineOffset = 0
	default:
		Error(s.line, s.lineOffset, fmt.Sprintf("Unexpected character: '%c'", c))
	}
}

func isWhitespace(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\f', '\r':
		return true
	}
	return false
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.lexStart = s.current
		s.ScanToken()
	}
	s.tokens = append(s.tokens, Token{EOF, "", s.line, s.lineOffset+1})
	return s.tokens
}
