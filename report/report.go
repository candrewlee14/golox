package report

import "fmt"

var HadError bool = false

func Report(line int, lineOffset int, where string, message string) {
	fmt.Printf("[line %d:%d] Error %s: %s\n", line, lineOffset, where, message)
}

func Error(line int, lineOffset int, message string) {
	HadError = true
	Report(line, lineOffset, "", message)
}
