package parser

import (
	"fmt"
	"golox/ast"
	"golox/lexer"
	"golox/token"
)

type (
	prefixParseFn func() ast.Expr
	infixParseFn  func(ast.Expr) ast.Expr
)

type ParserError struct {
	msg string
}

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []ParserError

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []ParserError{},
	}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = map[token.TokenType]prefixParseFn{
		token.IDENTIFIER: p.parseIdent,
	}

	return p
}

func (p *Parser) Errors() []ParserError {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.ScanToken()
}

func (p *Parser) addError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, ParserError{msg: msg})
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Stmt{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Stmt {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	default:
		p.errors = append(p.errors,
			ParserError{msg: fmt.Sprintf(
				"Expected the beginning of a statement, like 'var x = 100' at line %d:%d",
				p.curToken.Line, p.curToken.LineOffset)})
		return nil
	}
}

func (p *Parser) matchPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) parseVarStmt() *ast.VarStmt {
	stmt := &ast.VarStmt{Token: p.curToken}
	if !p.matchPeek(token.IDENTIFIER) {
		p.addError(token.IDENTIFIER)
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken}
	if !p.matchPeek(token.EQUAL) {
		p.addError(token.EQUAL)
		return nil
	}
	// TODO: Fix later. Currently skip expressions until semicolon
	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
		if p.curToken.Type == token.EOF {
			p.addError(token.SEMICOLON)
			return nil
		}
	}
	return stmt
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{Token: p.curToken}
	// TODO: Fix later. Currently skip expressions until semicolon
	fmt.Println(p.curToken)
	p.nextToken()
	fmt.Println(p.curToken)
	stmt.ReturnValue = p.parseExpr(LOWEST)
	p.matchPeek(token.SEMICOLON)
	return stmt
}

// Expression precedence definitions
type Prec uint8

const (
	_ Prec = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

func (p *Parser) parseExpr(prec Prec) ast.Expr {
	prefix, found := p.prefixParseFns[p.curToken.Type]
	if !found {
		return nil
	}
	leftExp := prefix()

	return leftExp
}

func (p *Parser) parseIdent() ast.Expr {
	return ast.Identifier{Token: p.curToken}
}
