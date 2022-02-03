package lexer

import (
	"fmt"
	"strconv"
    . "github.com/candrewlee14/golox/token"
    "github.com/candrewlee14/golox/report"
)

var keywords map[string]TokenType = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Lexer struct {
	source     string
	tokens     []Token
	lexStart   int
	current    int
	lineOffset int
	line       int
}

func NewLexer(source string) Lexer {
    return Lexer{source, nil, 0, 0, 0, 0}
}

func (s *Lexer) isAtEnd() bool {
	return s.current >= len(s.source)
}

// Consumes and returns current char
func (s *Lexer) advance() byte {
	c := s.source[s.current]
	s.current += 1
	s.lineOffset += 1
	return c
}

// Returns current char without consuming
// Returns null char if EOF
func (s *Lexer) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

// Returns next char without consuming
// Returns null char if EOF
func (s *Lexer) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

// Checks if current char matches given char and advances if so
func (s *Lexer) match(expected byte) bool {
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

func (s *Lexer) addToken(toktype TokenType) {
	lex := s.source[s.lexStart:s.current]
	s.tokens = append(s.tokens, NewToken(toktype, lex, s.line, s.lineOffset, nil))
}

func (s *Lexer) addTokenWithLiteral(toktype TokenType, val interface{}) {
	lex := s.source[s.lexStart:s.current]
	s.tokens = append(s.tokens, NewToken(toktype, lex, s.line, s.lineOffset, val))
}

func (s *Lexer) ScanToken() {
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
		s.current = s.current + 0
	case '\n':
		s.line += 1
		s.lineOffset = 0
	default:
		if isAlpha(c) {
			s.takeIdentifier()
		} else if isDigit(c) {
			s.takeNumber()
		} else {
			report.Error(s.line, s.lineOffset, fmt.Sprintf("Unexpected character: '%c'", c))
		}
	}
}

func (s *Lexer) takeIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	txt := s.source[s.lexStart:s.current]
	toktype, found := keywords[txt]
	if found {
		s.addToken(toktype)
	} else {
		s.addToken(IDENTIFIER)
	}
}

func (s *Lexer) takeString() {
	for !s.isAtEnd() && s.peek() != '"' {
		if s.peek() == '\n' {
			s.line += 1
			s.lineOffset = 0
		}
		s.advance()
	}
	if s.isAtEnd() {
		report.Error(s.line, s.lineOffset, "Unterminated string.")
		return
	}
	s.advance()
	str := s.source[s.lexStart+1 : s.current-1]
	s.addTokenWithLiteral(STRING, str)
}

func (s *Lexer) takeNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	f, err := strconv.ParseFloat(s.source[s.lexStart:s.current], 64)
	if err != nil {
		report.Error(s.line, s.lineOffset, "Invalid number.")
		return
	}
	s.addTokenWithLiteral(NUMBER, f)
}

func (s *Lexer) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.lexStart = s.current
		s.ScanToken()
	}
	s.tokens = append(s.tokens, NewToken(EOF, "", s.line, s.lineOffset + 1, nil))
	return s.tokens
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isDigit(c) || isAlpha(c)
}
