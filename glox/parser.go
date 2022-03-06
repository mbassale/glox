package glox

import (
	"fmt"
)

type Parser struct {
	tokens            []Token
	current           int
	errorReporter     ErrorReporter
	iteratorStmtCount int
}

type ParseError struct {
	message string
	token   Token
}

func (e ParseError) Error() string {
	return e.message
}

func NewParseError(message string, token Token) error {
	parseError := ParseError{
		message: message,
		token:   token,
	}
	return fmt.Errorf("ParseError: %w", parseError)
}

func NewParser(tokens []Token, errorReporter ErrorReporter) Parser {
	return Parser{
		tokens:            tokens,
		current:           0,
		errorReporter:     errorReporter,
		iteratorStmtCount: 0,
	}
}

func (p *Parser) Parse() []Stmt {
	statements := []Stmt{}
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			p.errorReporter.Push(p.currentLine(), "Parser", err)
			p.synchronize()
		} else {
			statements = append(statements, stmt)
		}
	}
	return statements
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == TOKEN_SEMICOLON {
			return
		}
		switch p.peek().Type {
		case TOKEN_CLASS:
			fallthrough
		case TOKEN_FUN:
			fallthrough
		case TOKEN_VAR:
			fallthrough
		case TOKEN_FOR:
			fallthrough
		case TOKEN_IF:
			fallthrough
		case TOKEN_WHILE:
			fallthrough
		case TOKEN_PRINT:
			fallthrough
		case TOKEN_RETURN:
			return
		}

		p.advance()
	}
}

/*
 * program -> declaration* EOF ;
 */
func (p *Parser) declaration() (Stmt, error) {
	if p.match(TOKEN_VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

/*
 * declaration -> varDeclaration | statement ;
 */
func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(TOKEN_IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr = nil
	if p.match(TOKEN_EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(TOKEN_SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return NewVarStmt(name, initializer), nil
}

/*
 * statement -> forStatement | ifStatement | printStatement | whileStatement | breakStatement | continueStatement |
 				block | expressionStatement ;
*/
func (p *Parser) statement() (Stmt, error) {
	if p.match(TOKEN_FOR) {
		forStmt, err := p.forStatement()
		if err != nil {
			return nil, err
		}
		p.iteratorStmtCount--
		return forStmt, nil
	}
	if p.match(TOKEN_IF) {
		return p.ifStatement()
	}
	if p.match(TOKEN_PRINT) {
		return p.printStatement()
	}
	if p.match(TOKEN_WHILE) {
		whileStmt, err := p.whileStatement()
		if err != nil {
			return nil, err
		}
		p.iteratorStmtCount--
		return whileStmt, nil
	}
	if p.match(TOKEN_BREAK) {
		return p.breakStatement()
	}
	if p.match(TOKEN_CONTINUE) {
		return p.continueStatement()
	}
	if p.match(TOKEN_LEFT_BRACE) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return NewBlockStmt(statements), nil
	}
	return p.expressionStatement()
}

/*
 * block -> "{" declaration* "}" ;
 */
func (p *Parser) block() ([]Stmt, error) {
	statements := []Stmt{}

	for !p.check(TOKEN_RIGHT_BRACE) && !p.isAtEnd() {
		declStmt, err := p.declaration()
		if err != nil {
			return statements, err
		}
		statements = append(statements, declStmt)
	}

	_, err := p.consume(TOKEN_RIGHT_BRACE, "Expect '}' after block.")
	return statements, err
}

/*
 * printStatement -> "print" expression ";" ;
 */
func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TOKEN_SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return NewPrintStmt(value), nil
}

/*
 * expressionStatement -> expression ";" ;
 */
func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TOKEN_SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return NewExpressionStmt(expr), nil
}

/*
 * forStatement -> "for" "(" ( varDeclaration | expressionStatement | ";" ) expression? ";" expression? ")" statement ;
 */
func (p *Parser) forStatement() (Stmt, error) {
	_, err := p.consume(TOKEN_LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}
	var initializer Stmt = nil
	if p.match(TOKEN_SEMICOLON) {
		initializer = nil
	} else if p.match(TOKEN_VAR) {
		initializer, err = p.varDeclaration()
	} else {
		initializer, err = p.expressionStatement()
	}
	if err != nil {
		return nil, err
	}

	var condition Expr = nil
	if !p.check(TOKEN_SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(TOKEN_SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment Expr = nil
	if !p.check(TOKEN_RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(TOKEN_RIGHT_PAREN, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	p.iteratorStmtCount++

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = NewBlockStmt([]Stmt{
			body,
			NewExpressionStmt(increment),
		})
	}

	if condition == nil {
		condition = NewLiteralExpr(true, body.getLine())
	}
	body = NewWhileStmt(condition, body)

	if initializer != nil {
		body = NewBlockStmt([]Stmt{
			initializer,
			body,
		})
	}

	return body, nil
}

/*
 * ifStatement -> "if" "(" expression ")" statement ( "else" statement )? ;
 */
func (p *Parser) ifStatement() (Stmt, error) {
	_, err := p.consume(TOKEN_LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TOKEN_RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt = nil
	if p.match(TOKEN_ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return NewIfStmt(condition, thenBranch, elseBranch), nil
}

/*
 * whileStatement -> "while" "(" expression ")" statement ;
 */
func (p *Parser) whileStatement() (Stmt, error) {
	_, err := p.consume(TOKEN_LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TOKEN_RIGHT_PAREN, "Expect ')' after 'while'.")
	if err != nil {
		return nil, err
	}

	p.iteratorStmtCount++

	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return NewWhileStmt(condition, body), nil
}

/*
 * breakStatement -> "break" ";" ;
 */
func (p *Parser) breakStatement() (Stmt, error) {
	if p.iteratorStmtCount <= 0 {
		return nil, NewParseError("break without for or while statement", p.previous())
	}
	stmt := NewBreakStmt(p.previous())
	_, err := p.consume(TOKEN_SEMICOLON, "Expect ';' after break.")
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

/*
 * continueStatement -> "continue" ";" ;
 */
func (p *Parser) continueStatement() (Stmt, error) {
	if p.iteratorStmtCount <= 0 {
		return nil, NewParseError("continue without for or while statement", p.previous())
	}
	stmt := NewContinueStmt(p.previous())
	_, err := p.consume(TOKEN_SEMICOLON, "Expect ';' after continue.")
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

/*
 * expression -> assignment ;
 */
func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

/*
 * assignment -> IDENTIFIER "=" assignment | conditionalExpression ;
 */
func (p *Parser) assignment() (Expr, error) {
	expr, err := p.conditionalExpression()
	if err != nil {
		return nil, err
	}

	if p.match(TOKEN_EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		variableExpr, ok := expr.(VariableExpr)
		if ok {
			variableName := variableExpr.Name
			return NewAssignExpr(variableName, value), nil
		}
		return nil, NewParseError("Invalid assignment target.", equals)
	}

	return expr, nil
}

/*
 * conditionalExpression -> logicOr | ( logicOr "?" expression ":" expression )
 */
func (p *Parser) conditionalExpression() (Expr, error) {
	expr, err := p.logicOr()
	if err != nil {
		return nil, err
	}

	if p.match(TOKEN_QUESTION) {
		left, err := p.expression()
		if err != nil {
			return nil, err
		}
		if p.match(TOKEN_COLON) {
			right, err := p.expression()
			if err != nil {
				return nil, err
			}
			return NewConditionalExpr(expr, left, right), nil
		} else {
			// missing right side expression
			return nil, NewParseError("Expecting ':' in conditional expression", p.tokens[p.current])
		}
	}

	return expr, nil
}

/*
 * logicOr -> logicAnd ( "or" logicAnd )* ;
 */
func (p *Parser) logicOr() (Expr, error) {
	expr, err := p.logicAnd()
	if err != nil {
		return nil, err
	}

	for p.match(TOKEN_OR) {
		operator := p.previous()
		right, err := p.logicAnd()
		if err != nil {
			return nil, err
		}
		expr = NewLogicalExpr(expr, operator, right)
	}

	return expr, nil
}

/*
 * logicAnd -> equality ( "and" equality )* ;
 */
func (p *Parser) logicAnd() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(TOKEN_AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = NewLogicalExpr(expr, operator, right)
	}

	return expr, nil
}

/*
 * equality -> comparison ( ( "!=" | "==" ) comparison )* ;
 */
func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(TOKEN_BANG_EQUAL, TOKEN_EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

/*
 * comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
 */
func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(TOKEN_GREATER, TOKEN_GREATER_EQUAL, TOKEN_LESS, TOKEN_LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

/*
 * term -> factor ( ( "-" | "+" ) factor )* ;
 */
func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(TOKEN_MINUS, TOKEN_PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

/*
 * factor -> unary ( ( "/" | "*" ) unary )* ;
 */
func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(TOKEN_SLASH, TOKEN_STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

/*
 * unary -> ( "!" | "-" ) unary | primary ;
 */
func (p *Parser) unary() (Expr, error) {
	if p.match(TOKEN_BANG, TOKEN_MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return NewUnaryExpr(operator, right), nil
	}

	return p.primary()
}

/*
 * primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER ;
 */
func (p *Parser) primary() (Expr, error) {
	if p.match(TOKEN_FALSE) {
		return NewLiteralExpr(false, p.currentLine()), nil
	}
	if p.match(TOKEN_TRUE) {
		return NewLiteralExpr(true, p.currentLine()), nil
	}
	if p.match(TOKEN_NIL) {
		return NewLiteralExpr(nil, p.currentLine()), nil
	}
	if p.match(TOKEN_NUMBER, TOKEN_STRING) {
		return NewLiteralExpr(p.previous().Literal, p.currentLine()), nil
	}
	if p.match(TOKEN_LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return NewGroupingExpr(expr), nil
	}
	if p.match(TOKEN_IDENTIFIER) {
		return NewVariableExpr(p.previous()), nil
	}

	return nil, NewParseError("Expected expression.", p.peek())
}

func (p *Parser) match(types ...int) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType int) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == TOKEN_EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) currentLine() int {
	return p.previous().Line
}

func (p *Parser) consume(tokenType int, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return Token{}, NewParseError(message, p.peek())
}
