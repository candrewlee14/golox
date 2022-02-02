package main

import (
	"bufio"
	//"strings"
	//"text/scanner"
	"fmt"
	"os"
)

var hadError bool = false

func Report(line int, lineOffset int, where string, message string) {
	fmt.Printf("[line %d:%d] Error %s: %s\n", line, lineOffset, where, message)
}

func Error(line int, lineOffset int, message string) {
	hadError = true
	Report(line, lineOffset, "", message)
}

func Run(source string) {
	scanner := Scanner{source, nil, 0, 0, 0, 0}
	toks := scanner.ScanTokens()
	for _, tok := range toks {
		fmt.Println(tok)
	}
}

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
		hadError = false
	}
}

func RunFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(64)
	}
	Run(string(bytes))
	if hadError {
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
