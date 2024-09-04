package ast

import (
	"fmt"
	lx"golox/lexer"
)


// type Expr interface {NUMBER | STRING | BOOL | NULL | Unary | Binary | Grouping}


type Visitor interface {
	Visit(Expr) (string, error)
}

type Expr interface {
	Accept(Visitor) (string, error)
}

type PrintVisitor struct {}

func (p *PrintVisitor) Visit(e Expr) (string, error) {
	switch expr := e.(type) {
	case *Grouping:
		result, _ := expr.Expr.Accept(p)
		return fmt.Sprintf("(group %s)", result), nil
	case *Unary:
		result, _ := expr.Expr.Accept(p)
		return fmt.Sprintf("(%s %s)", expr.Operator.Lexeme, result), nil
	case *Binary:
		rLeft, _ := expr.Left.Accept(p)
		rRight, _ := expr.Right.Accept(p)
		return fmt.Sprintf("(%s (%s %s))", expr.Operator.Lexeme, rLeft, rRight), nil
	case *NUMBER:
		result, _ := expr.Accept(p)
		return result, nil
	case *STRING:
		result, _ := expr.Accept(p)
		return result, nil
	case *BOOL:
		result, _ := expr.Accept(p)
		return result, nil
	case *NULL:
		result, _ := expr.Accept(p)
		return result, nil
	default:
		return "", fmt.Errorf("UNKNOWN EXPRESSION")
	}
}

type Grouping struct {
	Expr  Expr
}

func (g *Grouping) Accept(v Visitor) (string, error) {
	return v.Visit(g)
}
func (g *Grouping) ToExpr() Expr {return Expr(g)}

type Unary struct {
	Operator lx.Token
	Expr Expr
}

func (g *Unary) Accept(v Visitor) (string, error) {
	return v.Visit(g)
}
func (g *Unary) ToExpr() Expr {return Expr(g)}

type Binary struct {
	Operator lx.Token
	Left Expr
	Right Expr
}
func (g *Binary) Accept(v Visitor) (string, error) {
	return v.Visit(g)
}
func (g *Binary) ToExpr() Expr {return Expr(g)}

type NUMBER struct {Value float32}
type STRING struct {Value string}
type NULL struct {Value any}
type BOOL struct {Value bool}

func (g *NUMBER) Accept(v Visitor) (string, error) {
	return fmt.Sprintf("%f", g.Value), nil
}
func (g *NUMBER) ToExpr() Expr {
	return Expr(g)
}

func (g *STRING) Accept(v Visitor) (string, error) {
	return g.Value, nil
}
func (g *STRING) ToExpr() Expr {
	return Expr(g)
}
func (g *NULL) Accept(v Visitor) (string, error) {
	return "null", nil
}
func (g *NULL) ToExpr() Expr {
	return Expr(g)
}
func (g *BOOL) Accept(v Visitor) (string, error) {
	return fmt.Sprintf("%t", g.Value), nil
}
func (g *BOOL) ToExpr() Expr {
	return Expr(g)
}
