package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"golox/interp"
	"golox/lexer"
	"golox/parser"
	"golox/report"
	"os"
)

// Run interprets source code
func Run(source string, intp *interp.Interpreter) {
	scanner := lexer.NewLexer(source)
	p := parser.New(&scanner)
	prog := p.ParseProgram()
	es := p.Errors()
	if len(es) > 0 {
		fmt.Printf("%s\n", color.MagentaString("%d parsing errors encountered.", len(es)))
		for _, e := range p.Errors() {
			fmt.Printf("%s %s\n", color.RedString("Error:"), e)
		}
	} else {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(color.RedString("Runtime Error:"), err)
			}
		}()
		obj := intp.Eval(prog)
		if obj != nil {
			fmt.Println(color.BlueString("%s", prog), "->", color.GreenString("%s", obj))
		} else {
			fmt.Println(color.BlueString("%s", prog))
		}
		intp.PrintEnv()
	}
}

// RunPrompt interprets lines in a REPL
func RunPrompt() {
	reader := bufio.NewReader(os.Stdin)
	intp := interp.New()
	fmt.Print("> ")
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			os.Exit(64)
		}
		Run(string(line), &intp)
		fmt.Print("> ")
		report.HadError = false
	}
}

// RunFile interprets a file
func RunFile(path string) {
	intp := interp.New()
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(64)
	}
	Run(string(bytes), &intp)
	if report.HadError {
		os.Exit(65)
	}
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		RunFile(os.Args[1])
	} else {
		RunPrompt()
	}
}
