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
		token.NIL:        p.parseNil,
		token.STRING:     p.parseStr,
		token.TRUE:       p.parseBool,
		token.FALSE:      p.parseBool,
		token.BANG:       p.parsePrefixExpr,
		token.MINUS:      p.parsePrefixExpr,
		token.LEFT_PAREN: p.parseGroupedExpr,
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
		token.LEFT_PAREN:    p.parseCallExpr,
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
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead",
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

// func (p *Parser) sync() {
// 	for p.curToken.Type != token.EOF {
// 		if p.curToken.Type == token.SEMICOLON {
// 			p.nextToken()
// 			return
// 		}
// 		switch p.curToken.Type {
// 		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
// 			return
// 		}
// 		p.nextToken()
// 	}
// }

func (p *Parser) parseStatement() ast.Stmt {
	// if p.curToken.Type == token.IDENTIFIER && p.peekToken.Type == token.LEFT_PAREN {
	//     return p.parseCallExpr()
	// }
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStmt()
	case token.FUN:
		return p.parseFuncDeclStmt()
	case token.LEFT_BRACE:
		return p.parseBlockStmt()
	case token.IF:
		return p.parseIfStmt()
	case token.WHILE:
		return p.parseWhileStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	case token.PRINT:
		return p.parsePrintStmt()
	default:
		if p.curToken.Type == token.IDENTIFIER && p.peekToken.Type == token.EQUAL {
			return p.parseAssignStmt()
		}
		return p.parseExprStmt()
	}
}

func (p *Parser) matchPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) parseCallExpr(funcExpr ast.Expr) ast.Expr {
	// because Lox doesn't have lambdas, we know the function call is a identifier
	funcIdent, ok := funcExpr.(ast.Identifier)
	if !ok {
		p.errors = append(p.errors,
			ParserError{fmt.Sprintf("Expected function identifier, got %t", funcExpr)})
		p.advancePast(token.RIGHT_BRACE)
		return nil
	}
	callExpr := &ast.CallExpr{Token: funcIdent.Token, Function: &funcIdent}
	// we know curToken is LEFT_PAREN because the map infixParseFns
	p.nextToken()
	for p.curToken.Type != token.RIGHT_PAREN {
		if p.curToken.Type == token.EOF {
			p.errors = append(p.errors,
				ParserError{fmt.Sprintf("Expected \")\", found end of file instead.")})
			return nil
		}
		arg := p.parseExpr(LOWEST)
		callExpr.Args = append(callExpr.Args,
			arg)

		p.nextToken()
		if p.curToken.Type == token.COMMA {
			p.nextToken()
			continue
		} else if p.curToken.Type == token.RIGHT_PAREN {
			break
		} else {
			p.errors = append(p.errors,
				ParserError{fmt.Sprintf("Expected comma separating argument identifiers, found %s", p.curToken.Type)})
			p.advancePast(token.SEMICOLON)
			return nil
		}
	}
	return callExpr
}

func (p *Parser) parsePrintStmt() *ast.PrintStmt {
	ps := &ast.PrintStmt{Token: p.curToken}
	p.nextToken()
	ps.Expr = p.parseExpr(LOWEST)
	if !p.matchPeek(token.SEMICOLON) {
		p.addError(token.SEMICOLON)
		return nil
	}
	return ps
}

func (p *Parser) parseFuncDeclStmt() *ast.FuncDeclStmt {
	stmt := &ast.FuncDeclStmt{Token: p.curToken}
	p.nextToken()
	if p.curToken.Type != token.IDENTIFIER {
		p.errors = append(p.errors,
			ParserError{fmt.Sprintf("Expected function name identifier, got %s", p.curToken.Type)})
		p.advancePast(token.RIGHT_BRACE)
		return nil
	}
	ident := p.parseIdent().(ast.Identifier)
	stmt.Name = &ident
	p.nextToken()
	if p.curToken.Type != token.LEFT_PAREN {
		p.errors = append(p.errors,
			ParserError{fmt.Sprintf("Expected \"(\" after \"fun\".")})
		p.advancePast(token.RIGHT_BRACE)
		return nil
	}
	p.nextToken()
	paramNames := make(map[string]struct{})
	exists := struct{}{}
	for p.curToken.Type != token.RIGHT_PAREN {
		if p.curToken.Type == token.EOF {
			p.errors = append(p.errors,
				ParserError{fmt.Sprintf("Expected \")\", found end of file instead.")})
			return nil
		}
		// look for identifier
		if p.curToken.Type == token.IDENTIFIER {
			param := p.parseIdent().(ast.Identifier)
			stmt.Params = append(stmt.Params,
				&param)
			_, dup := paramNames[param.String()]
			if dup {
				p.errors = append(p.errors,
					ParserError{fmt.Sprintf("Found duplicate parameter identifier %q for function %q", param, stmt.Name)})
			} else {
				paramNames[param.String()] = exists
			}

		} else {
			p.errors = append(p.errors,
				ParserError{fmt.Sprintf("Expected parameter identifier, found %s", p.curToken.Type)})
			p.advancePast(token.RIGHT_BRACE)
			return nil
		}
		p.nextToken()
		// Look for comma or closing paren
		if p.curToken.Type == token.COMMA {
			p.nextToken()
			continue
		} else if p.curToken.Type == token.RIGHT_PAREN {
			p.nextToken()
			break
		} else {
			p.errors = append(p.errors,
				ParserError{fmt.Sprintf("Expected comma separating parameter identifiers, found %s", p.curToken.Type)})
			p.advancePast(token.RIGHT_BRACE)
			return nil
		}
	}
	fmt.Println(p.curToken.Lexeme)

	blockStmt := p.parseBlockStmt()
	stmt.Body = blockStmt

	fmt.Println(p.curToken.Lexeme)

	return stmt
}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	p.nextToken()

	block := &ast.BlockStmt{}
	block.Statements = []ast.Stmt{}

	for p.curToken.Type != token.RIGHT_BRACE {
		if p.curToken.Type == token.EOF {
			fmt.Print(block)
			p.addError(token.RIGHT_BRACE)
			return nil
		}
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		if p.curToken.Type != token.RIGHT_BRACE {
			p.nextToken()
		}
	}
	// if p.curToken.Type == token.RIGHT_BRACE {
	//     p.nextToken()
	// }
	return block
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
	if !p.matchPeek(token.SEMICOLON) {
		p.addError(token.SEMICOLON)
		return nil
	}
	return stmt
}

func (p *Parser) parseIfStmt() *ast.IfStmt {
	stmt := &ast.IfStmt{Token: p.curToken}
	p.nextToken()
	stmt.Cond = p.parseExpr(LOWEST)
	p.nextToken()
	stmt.OnTrue = p.parseBlockStmt()
	p.nextToken()

	fmt.Println(p.curToken.Type)
	if p.curToken.Type == token.ELSE {
		p.nextToken()
		stmt.OnFalse = p.parseBlockStmt()
	}
	return stmt
}

func (p *Parser) parseWhileStmt() *ast.WhileStmt {
	stmt := &ast.WhileStmt{Token: p.curToken}
	p.nextToken()
	stmt.Cond = p.parseExpr(LOWEST)
	p.nextToken()
	stmt.Body = p.parseBlockStmt()
	return stmt
}

func (p *Parser) parseAssignStmt() *ast.AssignStmt {
	stmt := &ast.AssignStmt{Name: ast.Identifier{Token: p.curToken}}
	p.nextToken()
	p.nextToken() // pass over the EQUAL token
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

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{Token: p.curToken}
	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
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
				ParserError{fmt.Sprintf("Expected \";\" after %q at line %d:%d",
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
				ParserError{fmt.Sprintf("Expected to find %s before EOF", toktype)})
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
	token.LEFT_PAREN:    CALL,
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

// Parsing any type of expression should end on the last token of the expression
// So peekToken after parsing all of a exprStmt should be semicolon

func (p *Parser) parseExprStmt() *ast.ExprStmt {
	stmt := &ast.ExprStmt{Token: p.curToken}
	stmt.Expr = p.parseExpr(LOWEST)
	// This no longer works correctly because function calls
	// will both need a prefix function after the semicolon
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

func (p *Parser) parseNil() ast.Expr {
	return ast.NilExpr{Token: p.curToken}
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

func (p *Parser) parseGroupedExpr() ast.Expr {
	p.nextToken()
	exp := p.parseExpr(LOWEST)
	if p.peekToken.Type != token.RIGHT_PAREN {
		msg := fmt.Sprintf("expected ')', found %s", p.curToken.Type)
		p.errors = append(p.errors, ParserError{msg})
		p.advancePast(token.RIGHT_PAREN)
		p.nextToken()
		return nil
	}
	p.nextToken()
	return exp
}
