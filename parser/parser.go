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
		token.NUMBER:     p.parseNum,
		token.STRING:     p.parseStr,
		token.TRUE:       p.parseBool,
		token.FALSE:      p.parseBool,
		token.BANG:       p.parsePrefixExpr,
		token.MINUS:      p.parsePrefixExpr,
	}
	p.infixParseFns = map[token.TokenType]infixParseFn{
		token.PLUS:          p.parseInfixExpr,
		token.MINUS:         p.parseInfixExpr,
		token.SLASH:         p.parseInfixExpr,
		token.STAR:          p.parseInfixExpr,
		token.EQUAL_EQUAL:   p.parseInfixExpr,
		token.BANG_EQUAL:    p.parseInfixExpr,
		token.LESS:          p.parseInfixExpr,
		token.GREATER:       p.parseInfixExpr,
		token.LESS_EQUAL:    p.parseInfixExpr,
		token.GREATER_EQUAL: p.parseInfixExpr,
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
		return p.parseExprStmt()
		// p.errors = append(p.errors,
		// 	ParserError{msg: fmt.Sprintf(
		// 		"Expected the beginning of a statement, like 'var x = 100' at line %d:%d. Got=%s",
		// 		p.curToken.Line, p.curToken.LineOffset, p.curToken.Type.String())})
		// p.advancePast(token.SEMICOLON)
		// return nil
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
	p.nextToken()
	stmt.Value = p.parseExpr(LOWEST)
	p.matchPeek(token.SEMICOLON)
	return stmt
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{Token: p.curToken}
	if p.peekToken.Type == token.SEMICOLON {
		stmt.ReturnValue = nil
	} else {
		p.nextToken()
		preExprLine := p.curToken.Line
		preExprLineOffset := p.curToken.LineOffset
		stmt.ReturnValue = p.parseExpr(LOWEST)
		if stmt.ReturnValue == nil {
			// if parseExpr failed, then we should report an error and move past semicolon
			p.advancePast(token.SEMICOLON)
			postExprLine := p.curToken.Line
			postExprLineOffset := p.curToken.LineOffset
			p.errors = append(p.errors,
				ParserError{fmt.Sprintf("Invalid expression from %d:%d to %d:%d",
					preExprLine, preExprLineOffset,
					postExprLine, postExprLineOffset)})
		} else if p.peekToken.Type != token.SEMICOLON {
			p.errors = append(p.errors,
				ParserError{fmt.Sprintf("Expected ';' after %q at line %d:%d",
					p.curToken.Lexeme, p.peekToken.Line, p.peekToken.LineOffset)})
			p.advancePast(token.SEMICOLON)
		} else {
			p.nextToken()
		}
	}
	return stmt
}

func (p *Parser) advancePast(toktype token.TokenType) {
	for p.peekToken.Type != token.SEMICOLON {
		if p.peekToken.Type == token.EOF {
			p.errors = append(p.errors,
				ParserError{fmt.Sprintf("Expected to find %q before EOF", toktype)})
			break
		}
		p.nextToken()
	}
	p.nextToken()
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

var precedences = map[token.TokenType]Prec{
	token.EQUAL_EQUAL:   EQUALS,
	token.BANG_EQUAL:    EQUALS,
	token.LESS:          LESSGREATER,
	token.GREATER:       LESSGREATER,
	token.LESS_EQUAL:    LESSGREATER,
	token.GREATER_EQUAL: LESSGREATER,
	token.PLUS:          SUM,
	token.MINUS:         SUM,
	token.SLASH:         PRODUCT,
	token.STAR:          PRODUCT,
}

func (p *Parser) peekPrec() Prec {
	if pr, ok := precedences[p.peekToken.Type]; ok {
		return pr
	}
	return LOWEST
}

func (p *Parser) curPrec() Prec {
	if pr, ok := precedences[p.curToken.Type]; ok {
		return pr
	}
	return LOWEST
}

func (p *Parser) parseExprStmt() *ast.ExprStmt {
	stmt := &ast.ExprStmt{Token: p.curToken}
	stmt.Expr = p.parseExpr(LOWEST)
	if p.peekToken.Type != token.SEMICOLON {
		p.errors = append(p.errors,
			ParserError{fmt.Sprintf("Expected ';' after %q at line %d:%d",
				p.curToken.Lexeme, p.peekToken.Line, p.peekToken.LineOffset)})
		p.advancePast(token.SEMICOLON)
	} else {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpr(prec Prec) ast.Expr {
	prefix, found := p.prefixParseFns[p.curToken.Type]
	if !found {
		msg := fmt.Sprintf("no prefix parse function for %s found", p.curToken.Type)
		p.errors = append(p.errors, ParserError{msg})
		return nil
	}
	leftExp := prefix()

	for p.peekToken.Type != token.SEMICOLON && prec < p.peekPrec() {
		infix, found := p.infixParseFns[p.peekToken.Type]
		if !found {
			msg := fmt.Sprintf("no infix parse function for %s found", p.curToken.Type)
			p.errors = append(p.errors, ParserError{msg})
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIdent() ast.Expr {
	return ast.Identifier{Token: p.curToken}
}

func (p *Parser) parseNum() ast.Expr {
	return ast.NumExpr{Token: p.curToken}
}

func (p *Parser) parseStr() ast.Expr {
	return ast.StrExpr{Token: p.curToken}
}

func (p *Parser) parseBool() ast.Expr {
	return ast.BoolExpr{Token: p.curToken}
}

func (p *Parser) parsePrefixExpr() ast.Expr {
	expr := &ast.PrefixExpr{Token: p.curToken}
	p.nextToken()
	expr.Right = p.parseExpr(PREFIX)
	return expr
}
func (p *Parser) parseInfixExpr(left ast.Expr) ast.Expr {
	expr := &ast.InfixExpr{
		Token: p.curToken,
		Left:  left,
	}
	prec := p.curPrec()
	p.nextToken()
	expr.Right = p.parseExpr(prec)
	return expr
}
