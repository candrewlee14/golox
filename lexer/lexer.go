package lexer

import (
	"fmt"
	"golox/report"
	. "golox/token"
	"strconv"
)

var keywords map[string]TokenType = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
    "false":  FALSE,
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
	lexStart   int
	current    int
	lineOffset int
	line       int
	lineStart  int
}

func NewLexer(source string) Lexer {
	return Lexer{source, 0, 0, 0, 0, 0}
}

func (s *Lexer) isAtEnd() bool {
	return s.current >= len(s.source)
}

// Consumes and returns current char
func (s *Lexer) advance() byte {
	if s.isAtEnd() {
		return '\000'
	}
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

func (s *Lexer) newToken(toktype TokenType) Token {
	lex := s.source[s.lexStart:s.current]
	return NewToken(toktype, lex, s.line, s.lexStart-s.lineStart, nil)
}

func (s *Lexer) newTokenWithLiteral(toktype TokenType, val interface{}) Token {
	lex := s.source[s.lexStart:s.current]
	return NewToken(toktype, lex, s.line, s.lexStart-s.lineStart, val)
}

func (s *Lexer) ScanToken() Token {
    found := false
    loop:for !found {
        found = true
        var res Token
        s.skipWhitespace()

        s.lexStart = s.current

        c := s.advance()
        switch c {
        case '(':
            res = s.newToken(LEFT_PAREN)
        case ')':
            res = s.newToken(RIGHT_PAREN)
        case '{':
            res = s.newToken(LEFT_BRACE)
        case '}':
            res = s.newToken(RIGHT_BRACE)
        case ',':
            res = s.newToken(COMMA)
        case '.':
            res = s.newToken(DOT)
        case '-':
            res = s.newToken(MINUS)
        case '+':
            res = s.newToken(PLUS)
        case ';':
            res = s.newToken(SEMICOLON)
        case '*':
            res = s.newToken(STAR)
        case '!':
            var toktype TokenType
            if s.match('=') {
                toktype = BANG_EQUAL
            } else {
                toktype = BANG
            }
            res = s.newToken(toktype)
        case '=':
            var toktype TokenType
            if s.match('=') {
                toktype = EQUAL_EQUAL
            } else {
                toktype = EQUAL
            }
            res = s.newToken(toktype)
        case '<':
            var toktype TokenType
            if s.match('=') {
                toktype = GREATER_EQUAL
            } else {
                toktype = GREATER
            }
            res = s.newToken(toktype)
        case '>':
            var toktype TokenType
            if s.match('=') {
                toktype = LESS_EQUAL
            } else {
                toktype = LESS
            }
            res = s.newToken(toktype)
        case '/':
            if s.match('/') {
                // ignore comment until end of line
                for !s.isAtEnd() && s.peek() != '\n' {
                    s.current += 1
                }
                found = false
            } else {
                res = s.newToken(SLASH)
            }
        case '\000':
            res = s.newToken(EOF)
        case '"':
            res = s.takeString()

        default:
            if isAlpha(c) {
                res = s.takeIdentifier()
            } else if isDigit(c) {
                res = s.takeNumber()
            } else {
                report.Error(s.line, s.lineOffset, fmt.Sprintf("Unexpected character: '%c'", c))
                break loop
            }
        }
        if found {
            return res
        }
    }
    return s.newToken(INVALID)
}

func (s *Lexer) skipWhitespace() {
loop:
	for {
		switch s.peek() {
		case '\n':
			s.line += 1
			s.lineOffset = 0
			s.current += 1
			s.lineStart = s.current
		case ' ', '\r', '\t':
			// just skip regular whitespace
			s.advance()
		default:
			break loop
		}
	}
}

func (s *Lexer) takeIdentifier() Token {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	txt := s.source[s.lexStart:s.current]
	toktype, found := keywords[txt]
	if found {
		return s.newToken(toktype)
	} else {
		return s.newToken(IDENTIFIER)
	}
}

func (s *Lexer) takeString() Token {
	for !s.isAtEnd() && s.peek() != '"' {
		if s.peek() == '\n' {
			s.line += 1
			s.lineOffset = 0
		}
		s.advance()
	}
	if s.isAtEnd() {
		report.Error(s.line, s.lineOffset, "Unterminated string.")
		return s.newToken(INVALID)
	}
	s.advance()
	str := s.source[s.lexStart+1 : s.current-1]
	return s.newTokenWithLiteral(STRING, str)
}

func (s *Lexer) takeNumber() Token {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
    // Error cannot be possible here since numbers are of form x or x.x
	f, _ := strconv.ParseFloat(s.source[s.lexStart:s.current], 64)
	return s.newTokenWithLiteral(NUMBER, f)
}

func (s *Lexer) ScanTokens() []Token {
	var tokens []Token
	for !s.isAtEnd() {
		tokens = append(tokens, s.ScanToken())
	}
	tokens = append(tokens, NewToken(EOF, "", s.line, s.lineOffset+1, nil))
	return tokens
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
