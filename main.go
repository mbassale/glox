package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mbassale/glox/glox"
)

const EXIT_BAD_ARGS = 64
const EXIT_ERROR = 65
const EXIT_RUNTIME_ERROR = 70

var hadError bool = false
var hadRuntimeError bool = false

func runFile(path string) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var errorReporter glox.ErrorReporter = glox.NewConsoleErrorReporter()
	run(string(contents), errorReporter)
	if hadError {
		os.Exit(EXIT_ERROR)
	}
	if hadRuntimeError {
		os.Exit(EXIT_RUNTIME_ERROR)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	var errorReporter glox.ErrorReporter = glox.NewConsoleErrorReporter()
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		run(line, errorReporter)
		errorReporter.ClearError()
	}
}

func run(source string, errorReporter glox.ErrorReporter) {
	scanner := glox.NewScanner(source, errorReporter)
	if errorReporter.HasError() {
		hadError = true
		return
	}
	tokens := scanner.ScanTokens()
	for _, token := range tokens {
		fmt.Printf("Token: %v\n", token)
	}
	parser := glox.NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	astPrinter := glox.AstPrinter{}
	fmt.Println(astPrinter.Print(statements))
	interpreter := glox.NewInterpreter(errorReporter)
	interpreter.Interpret(statements)
	if errorReporter.HasError() {
		hadRuntimeError = true
		return
	}
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
		os.Exit(EXIT_BAD_ARGS)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}
