package glox

import (
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
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

func NewParser(tokens []Token) Parser {
	return Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() (Expr, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	return expr, nil
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
 * expression -> conditionalExpression ;
 */
func (p *Parser) expression() (Expr, error) {
	return p.conditionalExpression()
}

/*
 * conditionalExpression -> equality | ( equality "?" expression ":" expression )
 */
func (p *Parser) conditionalExpression() (Expr, error) {
	expr, err := p.equality()
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
 * primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
 */
func (p *Parser) primary() (Expr, error) {
	if p.match(TOKEN_FALSE) {
		return NewLiteralExpr(false), nil
	}
	if p.match(TOKEN_TRUE) {
		return NewLiteralExpr(true), nil
	}
	if p.match(TOKEN_NIL) {
		return NewLiteralExpr(nil), nil
	}
	if p.match(TOKEN_NUMBER, TOKEN_STRING) {
		return NewLiteralExpr(p.previous()), nil
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

func (p *Parser) consume(tokenType int, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return Token{}, NewParseError(message, p.peek())
}
