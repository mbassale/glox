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
	var interpreter glox.Interpreter = glox.NewInterpreter(errorReporter)
	run(string(contents), &interpreter, errorReporter)
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
	var interpreter glox.Interpreter = glox.NewInterpreter(errorReporter)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		run(line, &interpreter, errorReporter)
		errorReporter.ClearError()
	}
}

func run(source string, interpreter *glox.Interpreter, errorReporter glox.ErrorReporter) {
	scanner := glox.NewScanner(source, errorReporter)
	if (errorReporter).HasError() {
		hadError = true
		return
	}
	tokens := scanner.ScanTokens()
	for _, token := range tokens {
		fmt.Printf("Token: %v\n", token)
	}
	parser := glox.NewParser(tokens, errorReporter)
	statements := parser.Parse()
	if errorReporter.HasError() {
		hadError = true
		return
	}
	astPrinter := glox.AstPrinter{}
	fmt.Println(astPrinter.Print(statements))
	lastValue, _ := interpreter.Interpret(statements)
	fmt.Printf("=%v\n", lastValue)
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
