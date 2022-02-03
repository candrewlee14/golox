package report

import "fmt"

// HadError holds the state for if an error is encoutered during interpreting
var HadError bool = false

// Report prints out a fancy error report
func Report(line int, lineOffset int, where string, message string) {
	fmt.Printf("[line %d:%d] Error %s: %s\n", line, lineOffset, where, message)
}

// Error sets HadError to true and prints an error report
func Error(line int, lineOffset int, message string) {
	HadError = true
	Report(line, lineOffset, "", message)
}
