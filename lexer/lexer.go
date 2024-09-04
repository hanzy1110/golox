package lexer

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var KEYWORDS map[string]TokenType

func init() {
	KEYWORDS = map[string]TokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"fun":    FUN,
		"for":    FOR,
		"if":     IF,
		"null":   NULL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
		"eof":    EOF,
	}
}


type GoloxExecution struct {
	Errors []ErrorInfo
	Tokens []Token
}

func NewGoloxExecution(tokens []Token) GoloxExecution {
	return GoloxExecution{Errors: make([]ErrorInfo, 0), Tokens: tokens}
}

func (gol *GoloxExecution) UpdateError(e ErrorInfo) {
	gol.Errors = append(gol.Errors, e)
}
func (gol *GoloxExecution) UpdateTokens(t []Token) {
	gol.Tokens = append(gol.Tokens, t...)
}

type ErrorInfo struct {
	Etype  GoloxError
	lineno int
	desc   string
}

func NewError(eType GoloxError, lineno int, col int, err error) ErrorInfo {
	actualError := fmt.Sprintf("ERROR in LINE %d, COL %d: %s", lineno, col, err)
	return ErrorInfo{Etype: eType, lineno: lineno, desc: actualError}
}

type GoloxError int

const (
	LEXICAL_ERROR GoloxError = iota
	SYNTAX_ERROR
)

func (g GoloxError) String() string {
	switch g {
	case LEXICAL_ERROR:
		return "LEXICAL_ERROR"
	case SYNTAX_ERROR:
		return "SYNTAX_ERROR"
	default:
		return fmt.Sprintf("%d", int(g))
	}
}

type Token struct {
	Tp      TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("LX=>'%s'", t.Lexeme)
}

func ScanLine(line string, lineno int, gol *GoloxExecution) {

	start := 0

	cap := len(line)
	var tokens []Token
	var c rune

	keepLexing := true
	for i := 0; i < cap; {
		i, c = advance(i, line)
		// log.Printf("i: %d", i)
		switch c {
		case '(':
			start, tokens = addToken(LEFT_PAREN, start, i, line, lineno, tokens)
		case ')':
			start, tokens = addToken(RIGHT_PAREN, start, i, line, lineno, tokens)
		case '{':
			start, tokens = addToken(LEFT_BRACE, start, i, line, lineno, tokens)
		case '}':
			start, tokens = addToken(RIGHT_BRACE, start, i, line, lineno, tokens)
		case '.':
			start, tokens = addToken(DOT, start, i, line, lineno, tokens)
		case '-':
			start, tokens = addToken(MINUS, start, i, line, lineno, tokens)
		case '+':
			start, tokens = addToken(PLUS, start, i, line, lineno, tokens)
		case ';':
			start, tokens = addToken(SEMICOLON, start, i, line, lineno, tokens)
		case '*':
			start, tokens = addToken(STAR, start, i, line, lineno, tokens)
		case '!':
			nc := getNext(line, i, cap)
			if nc == '=' {
				start, tokens = addToken(BANG_EQUAL, start, i+1, line, lineno, tokens)
				i++
			} else {
				start, tokens = addToken(BANG, start, i, line, lineno, tokens)
			}
		case '<':
			nc := getNext(line, i, cap)
			if nc == '=' {
				start, tokens = addToken(LESS_EQUAL, start, i+1, line, lineno, tokens)
				i++
			} else {
				start, tokens = addToken(LESS, start, i, line, lineno, tokens)
			}
		case '>':
			nc := getNext(line, i, cap)
			if nc == '=' {
				start, tokens = addToken(GREATER_EQUAL, start, i+1, line, lineno, tokens)
				i++
			} else {
				start, tokens = addToken(GREATER, start, i, line, lineno, tokens)
			}
		case '=':
			nc := getNext(line, i, cap)
			if nc == '=' {
				start, tokens = addToken(EQUAL_EQUAL, start, i+1, line, lineno, tokens)
				i++
			} else {
				start, tokens = addToken(EQUAL, start, i, line, lineno, tokens)
			}
		case '/':
			nc := getNext(line, i, cap)
			if nc == '/' {
				log.Println("Skipping lexing!")
				keepLexing = false
				// break
			} else {
				start, tokens = addToken(SLASH, start, i, line, lineno, tokens)
			}

		case ' ', '\r', '\t', 'n':
			start = i
			continue

		case '"':
			s_end, err := parseString(start, line)
			if err != nil {
				e := NewError(SYNTAX_ERROR, lineno, s_end, err)
				gol.UpdateError(e)
				keepLexing = false
			} else {
				start, tokens = addToken(STRING, start, s_end+1, line, lineno, tokens)
				i = s_end + 2
			}

		default:

			if isDigit(c) {
				n_end, err := parseNumber(start, line)
				if err != nil {
					e := NewError(SYNTAX_ERROR, lineno, n_end, err)
					gol.UpdateError(e)
					keepLexing = false
				} else {
					start, tokens = addToken(NUMBER, start, n_end, line, lineno, tokens)
					i = n_end
				}
			} else if isAlpha(c) {
				i_end, err := parseIdentifier(start, line)
				if err != nil {
					e := NewError(SYNTAX_ERROR, lineno, i_end, err)
					gol.UpdateError(e)
					keepLexing = false
				} else {
					if kw, f := KEYWORDS[line[start:i_end]]; f {
						start, tokens = addToken(kw, start, i_end, line, lineno, tokens)
					} else {
						start, tokens = addToken(IDENTIFIER, start, i_end, line, lineno, tokens)
					}
					i = i_end
				}

			} else {
				msg := fmt.Sprintf("Invalid Character %#U", c)
				e := NewError(LEXICAL_ERROR, lineno, i, fmt.Errorf(msg))
				gol.UpdateError(e)
			}
		}

		if !keepLexing {
			break
		}

	}
	gol.UpdateTokens(tokens)
}

func parseString(start int, line string) (int, error) {
	var c rune
	var end int
	for k := start + 1; k < len(line); k++ {
		c = rune(line[k])
		end = k
		if c == '"' {
			return k, nil
		}
	}
	return end, fmt.Errorf("UNTERMINATED STRING")
}

func parseNumber(start int, line string) (int, error) {
	var c rune
	var end int
	// log.Printf("VALS %d, %s", start, line)
	for k := start; k < len(line); k++ {
		c = rune(line[k])
		// log.Printf("RUNE: %#U, IDX: %d", c, k)
		end = k
		if c == '.' {
			continue
		} else if !isDigit(c) {
			end = k - 1
			break
		}
	}
	// log.Printf("PARSED NUM => %s S,E => %d, %d", line[start:end], start, end)
	// TODO asegurarse que todos los numeros sean en coma flotante
	if strings.Contains(line[start:end+1], ".") {
		if _, err := strconv.ParseFloat(string(line[start:end+1]), 32); err != nil {
			return end + 1, fmt.Errorf("INVALID NUMBER: %s, %v", line[start:end+1], err)
		}
	} else if _, err := strconv.Atoi(string(line[start : end+1])); err != nil {
		return end + 1, fmt.Errorf("INVALID NUMBER: %s, %v", line[start:end+1], err)
	}
	return end + 1, nil
}

func parseIdentifier(start int, line string) (int, error) {

	var c rune
	var end int
	for k := start; k < len(line); k++ {
		c = rune(line[k])
		end = k
		if !isAlphaNum(c) {
			return end, nil
			// if c == ' ' {
			// 	return end, nil
			// }
			// return end + 1, fmt.Errorf("INVALID CHAR IN IDENTIFIER %#U", c)
		}
	}
	return end + 1, nil
}

func advance(i int, line string) (int, rune) {
	var c rune
	out := i + 1
	if i < len(line) {
		c = rune(line[i])
	}
	return out, c
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isAlphaNum(r rune) bool {
	return isAlpha(r) || isDigit(r)
}

func getNext(line string, i int, cap int) (nc rune) {
	log.Println(line, i, cap)
	if i < cap {
		nc = rune(line[i])
	} else {
		// Tengo que devolver basura
		nc = '\x80'
	}
	return
}

func addToken(tp TokenType,
	start int, current int,
	line string, lineno int,
	tokens []Token) (cur int, tks []Token) {

	// log.Printf("VALUES => START: %d, CURRENT: %d, LINE: %s", start, current, line)
	token := Token{Tp: tp, Lexeme: line[start:current], Literal: nil, Line: lineno}
	tks = append(tokens, token)
	cur = current
	return

}
