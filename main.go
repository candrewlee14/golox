package main

import (
	"bufio"
	"fmt"
	"golox/lexer"
	"golox/report"
	"os"
)

// Run interprets source code
func Run(source string) {
	scanner := lexer.NewLexer(source)
	toks := scanner.ScanTokens()
	for _, tok := range toks {
		fmt.Println(tok)
	}
}

// RunPrompt interprets lines in a REPL
func RunPrompt() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			os.Exit(64)
		}
		Run(string(line))
		fmt.Print("> ")
		report.HadError = false
	}
}

// RunFile interprets a file
func RunFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(64)
	}
	Run(string(bytes))
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
