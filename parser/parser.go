package parser

import (
	"fmt"
	"golox/ast"
	lx "golox/lexer"
	"log"
	"strconv"
)

type Stream interface {
	Peek() lx.Token
	Next() lx.Token
	NTokens() int
}
type TokenStream []lx.Token

func (s *TokenStream) Peek() lx.Token {
	if len(*s) == 0 {
		return lx.Token{Tp:lx.EOF}
	}
	return (*s)[0]
}

func (s *TokenStream) Next() lx.Token {
	result := (*s)[0]
	*s = (*s)[1:]
	return result
}

func (s *TokenStream) NTokens() int {
	return len(*s)
}

func AsTokenStream(tks []lx.Token) Stream {
	ts:= TokenStream(tks)
	return &ts
}

func match(s Stream, args...lx.TokenType) (c lx.Token, b bool) {
	log.Printf("STREAM DURING MATCHES ===> %v", s)
	c = s.Peek()
	for _, t := range args {
		if check(t, c) {
			c = s.Next()
			return c, true
		}
	}
	return c, false
}

func checkFirst(s Stream, args...lx.TokenType) bool {
	c:=s.Peek()
	for _, t := range args {
		if check(t, c) {
			return false
		}
	}
	return true
}

func check(t lx.TokenType, c lx.Token) (cont bool) {
	if c.Tp==lx.EOF {
		return
	}
	return c.Tp == t
}

func consume(s Stream, t lx.TokenType, msg string) (err error) {

	if check(t, s.Next()) {
		return nil
	}

	return fmt.Errorf("SYNTAX ERROR! %s", msg)
}

func sanitizeTokens(s Stream) (err error) {

	if b := checkFirst(s, lx.EQUAL_EQUAL, lx.PLUS, lx.DOT, lx.SLASH, lx.SEMICOLON); !b {
		return fmt.Errorf("INVALID STARTING CHAR! => %v", s.Peek())
	}
	return
}

// expression -> equality
func expression(s Stream) (expr ast.Expr, err error) {
	if err = sanitizeTokens(s); err!=nil {
		return
	}
	expr, err = equality(s)

	if s.NTokens()>0 {
		err = fmt.Errorf("TOKENS LEFT! Parenthesize Expressions! %v", s)
	}
	return
}

// equality -> comparison (('!='|'==') comparison)*
func equality(s Stream) (expr ast.Expr, err error) {
	expr, err = comparison(s)
	b:=true
	c:=lx.Token{}
	for  {
		c, b=match(s, lx.BANG_EQUAL, lx.EQUAL_EQUAL)
		if !b {
			break
		}
		operator := c
		if right, err := comparison(s); err!=nil {
			log.Fatal("WHILE PARSING...", err)
		} else {
			expr = (&ast.Binary{Operator:operator, Left:expr, Right: right}).ToExpr()
			return expr, nil
		}
	}
	return
}

// comparison -> term (('>' | '>=' | '<' | '<=') term)*
func comparison(s Stream) (expr ast.Expr, err error) {
	expr, err = term(s)

	b:=true
	c:=lx.Token{}
	for {
		c, b=match(s, lx.GREATER, lx.GREATER_EQUAL, lx.LESS, lx.LESS_EQUAL)
		if !b {
			break
		}
		operator := c
		if right, err := term(s); err!=nil {
			log.Fatal("WHILE PARSING...", err)
		} else {
			expr = (&ast.Binary{Operator:operator, Left:expr, Right: right}).ToExpr()
			return expr, nil
		}
	}
	return
}

// term -> factor (('/'|'*') factor)*
func term(s Stream) (expr ast.Expr, err error) {
	expr, err = factor(s)
	b:=true
	c:=lx.Token{}
	for {
		c, b=match(s, lx.MINUS, lx.PLUS)
		if !b {
			break
		}
		operator := c
		if right, err := factor(s); err!=nil {
			log.Fatal("WHILE PARSING...", err)
		} else {
			expr = (&ast.Binary{Operator:operator, Left:expr, Right: right}).ToExpr()
			return expr, nil
		}
	}
	return
}

// factor -> unary (('/'|'*') unary)*
func factor(s Stream) (expr ast.Expr, err error) {
	expr, err = unary(s)
	b:=true
	c:=lx.Token{}
	for b {
		c, b=match(s, lx.SLASH, lx.STAR)
		if !b {
			break
		}
		operator := c
		if right, err := unary(s); err!=nil {
			log.Fatal("WHILE PARSING...", err)
		} else {
			expr = (&ast.Binary{Operator:operator, Left:expr, Right: right}).ToExpr()
			return expr, nil
		}
	}
	return
}

func unary(s Stream) (expr ast.Expr, err error) {
	b:=true
	c:=lx.Token{}
	for {
		c, b=match(s, lx.BANG, lx.MINUS)
		if !b {
			break
		}
		operator := c
		right, err := unary(s)
		expr = (&ast.Unary{Operator: operator, Expr: right}).ToExpr()
		return expr, err
	}
	return primary(s)
}

func primary(s Stream) (expr ast.Expr, err error) {

	if _,b:=match(s, lx.FALSE);b {
		expr = (&ast.BOOL{Value: false}).ToExpr()
		return
	}
	if _,b:=match(s, lx.TRUE);b {
		expr = (&ast.BOOL{Value: true}).ToExpr()
		return
	}
	if _,b:=match(s, lx.NULL);b {
		expr = (&ast.NULL{Value: nil}).ToExpr()
		return
	}
	if c,b:=match(s, lx.NUMBER);b {
		val, _ := strconv.ParseFloat(c.Lexeme, 32)
		expr = (&ast.NUMBER{Value: float32(val)}).ToExpr()
		return expr, nil
	}
	if c,b:=match(s, lx.STRING);b {
		expr = (&ast.STRING{Value: c.Lexeme}).ToExpr()
		return
	}
	if _,b := match(s, lx.LEFT_PAREN); b {
		expr, err = expression(s)
		if err = consume(s, lx.RIGHT_PAREN, "EXPECTED RIGHT PAREN AT THE END!"); err!=nil {
			return
		} else {
			expr = (&ast.Grouping{Expr:expr}).ToExpr()
			return
		}
	}
	err = fmt.Errorf("EXPECTED AN EXPRESSION!")
	return
}

func ParseStream(s Stream) (expr ast.Expr, err error) {
	log.Printf("STARTING TO PARSE => %v ", s)
	expr,err=expression(s)
	return
}
