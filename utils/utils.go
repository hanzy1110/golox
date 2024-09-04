package utils

type GoloxExecution struct {
	Errors []ErrorInfo
	Tokens []Token
}

func NewGoloxExecution(tokens []Token) GoloxExecution {
	return GoloxExecution{Errors:make([]ErrorInfo, 0), Tokens: tokens}
}

func (gol *GoloxExecution) UpdateError(e ErrorInfo) {
	gol.Errors = append(gol.Errors, e)
}
func (gol *GoloxExecution) UpdateTokens(t []Token) {
	gol.Tokens = append(gol.Tokens, t...)
}

type ErrorInfo struct {
	Etype GoloxError
	lineno int
	desc string
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
