package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	ast"golox/ast"
	lx "golox/lexer"
	ps "golox/parser"
	"log"
	"os"
)


var (
	inputFile string
	repl bool
)

func init() {

    flag.StringVar(&inputFile, "inputFile", "", "read and execute")
    flag.BoolVar(&repl, "repl", false, "Run the REPL???")
	flag.Parse()
}

func main() {

	if repl {
		runPrompt()
	}

	if inputFile == "" {
		log.Fatal("no files provided!!")
	}

	if _, err := os.Stat(inputFile); err!=nil {
		log.Fatal("Don't know what happened!", err)
	} else if errors.Is(err, os.ErrNotExist) {
		log.Println("file does not exist!", err)
		os.Exit(1)
	} else {
		runFile(inputFile)
	}


}

func runFile(inputFile string) error {
	content, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(content)

	lineno:=0
	gol := lx.NewGoloxExecution(nil)

	for scanner.Scan() {
		line := scanner.Text()
		lineno++
		InterpretLine(line, lineno, false, &gol)
		log.Printf("LINE => %s -- TOKENS => %d", line, len(gol.Tokens))
	}
	if len(gol.Errors)!=0 {
		for _, e := range gol.Errors {
			log.Printf("%v", e)
		}
	} else {
		log.Println(gol)
		log.Println("COMPILATION SUCCESFUL!")
	}

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

	return nil
}


func runPrompt() error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(">> ")
	gol := lx.NewGoloxExecution(nil)
	for scanner.Scan() {

		fmt.Print(">>>> ")
		switch pl := scanner.Text(); pl {
		case "exit", "q":
			log.Println("Bye!")
			os.Exit(0)
		case "es!":
			log.Printf("GOLOX EXECUTION STATE ==> %v", gol)
		default:
			InterpretLine(pl, 1, true, &gol)
		}
	if len(gol.Errors)!=0 {
		for _, e := range gol.Errors {
			log.Printf("%v", e)
		}
	}
	gol = lx.NewGoloxExecution(nil)
	fmt.Print(">> ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}


func InterpretLine(line string, lineno int, execute bool, gol *lx.GoloxExecution) {

	if line == "error" {
		// En algun momento hacer que el reporte produzca el objeto
		e := lx.NewError(lx.SYNTAX_ERROR, lineno, 0, fmt.Errorf("You typed error!"))
		gol.UpdateError(e)
		return
	}

	lx.ScanLine(line, lineno, gol)

	if len(gol.Errors)==0 && execute {
		v := &ast.PrintVisitor{}
		ParseLine(gol, v)
		ExecuteLine(line, gol)
	}

	return
}

func ParseLine(gol *lx.GoloxExecution, v *ast.PrintVisitor) (err error) {

	s := ps.AsTokenStream(gol.Tokens)
	if expr, err := ps.ParseStream(s); err!=nil {
		log.Fatal("PARSING ERROR => ", err)
	} else {
		pexp, _ := v.Visit(expr)
		log.Printf("PARSED EXPRESSION => %s", pexp)
	}

	return
}

func ExecuteLine(line string, gol *lx.GoloxExecution) {
	log.Println("Executing line ->", line, " -- Tokens were -->",
		gol.Tokens, " -- TOTAL TOKENS =>", len(gol.Tokens))
	return
}

