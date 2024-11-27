package main

import (
	"fmt"
	"os"
	"strings"
)

type varScope = map[string]float32
type funcScope = map[string]Function

type Function interface {
	run(args []any, varScope *varScope, funcScope *funcScope) any
}

type Arg struct {
	Name string
	Type string
}

type UserFunction struct {
	args []Arg
	// returns string // todo
	body Tokens
}

func newScopeExtends(parentVars *varScope, parentFuncs *funcScope) (*varScope, *funcScope) {
	var blockVarScope = make(varScope)
	var blockFuncScope = make(funcScope)

	for k, v := range *parentVars {
		blockVarScope[k] = v
	}

	for k, v := range *parentFuncs {
		blockFuncScope[k] = v
	}

	return &blockVarScope, &blockFuncScope
}

func (self *UserFunction) run(args []any, vScope *varScope, fScope *funcScope) any {
	var blockVarScope, blockFuncScope = newScopeExtends(vScope, fScope)

	if len(args) != len(self.args) {
		panic("Incorrect number of arguments: " + string(len(args))) // todo, migrate to `throw`n errors
	}

	for i, arg := range args {
		(*blockVarScope)[self.args[i].Name] = arg.(float32) // todo, temp
	}

	return RunTokens(self.body, blockVarScope, blockFuncScope)
}

type STDFunction struct {
	args []Arg
	// returns string // todo
	Func func(args ...any) any
}

func (self *STDFunction) run(args []any, varScope *varScope, funcScope *funcScope) any {
	return self.Func(args...)
}

func assertToken(t Tokens, i int, expected string) {
	if t.tokens[i] == expected {
		return
	}

	errorLine := getFileLine(t)
	tokensStart := fileLines[errorLine-1]
	tokensEnd := fileLines[errorLine]
	line := fileTokens.tokens[tokensStart:tokensEnd]

	message := "Got '%s', expected '%s'\n  %s\n  %s%s"

	throw(
		fmt.Sprintf(message,
			t.tokens[i],
			expected,
			magenta+strings.Join(line, " "),
			strings.Repeat(" ", len(strings.Join(line[:i], " "))+1),
			strings.Repeat("^", len(t.tokens[i]))+reset,
		), t, errorLine)
}

func getFileLine(t Tokens) int {
	for i, line := range fileLines {
		if line > t.fileLineIndex+t.index {
			return i
		}
	}
	return len(fileLines)
}

func throw(message string, t Tokens, line ...int) {
	var l int
	if len(line) > 0 {
		l = line[0]
	} else {
		l = getFileLine(t)
	}
	fmt.Print(grey)
	fmt.Printf("[%s:%d]", fileName, l)
	fmt.Println(red, message, reset)
	os.Exit(1)
}

func VarDefinition(tokens Tokens, varScope *varScope, funcScope *funcScope) {
	tokens.Consume() // "var"
	name := tokens.Consume()
	tokens.AssertConsumption("=")
	(*varScope)[name] = ParseExpression(tokens, varScope, funcScope)
}
func FuncDefinition(tokens Tokens, varScope *varScope, funcScope *funcScope) {
	tokens.Consume() // "function"
	name := tokens.Consume()

	allArgTokens := tokens.ConsumeTuple()
	var args []Arg

	for _, argTokens := range allArgTokens {
		if argTokens.IsEmpty() {
			break
		}
		argName := argTokens.Consume()
		argTokens.AssertConsumption(":")
		args = append(args, Arg{argName, argTokens.Consume()})
	}

	funcBody := tokens.ConsumeCurlyBrackets()

	(*funcScope)[name] = &UserFunction{args, funcBody}
}
func FuncCall(tokens Tokens, varScope *varScope, funcScope *funcScope) any {
	name := tokens.Consume()
	args := tokens.ConsumeTuple()

	function, ok := (*funcScope)[name]

	if ok {
		var parsedArgs []any
		for _, arg := range args {
			parsedArgs = append(parsedArgs, ParseExpression(arg, varScope, funcScope))
		}
		return function.run(parsedArgs, varScope, funcScope)
	} else {
		panic("unreachable")
	}
}
func VarUpdate(tokens Tokens, varScope *varScope, funcScope *funcScope) {
	name := tokens.Consume()
	tokens.AssertConsumption("=")
	(*varScope)[name] = ParseExpression(tokens, varScope, funcScope)
}
func Return(tokens Tokens, varScope *varScope, funcScope *funcScope) float32 {
	return ParseExpression(tokens.Slice(1, -1), varScope, funcScope)
}
func IfBlock(tokens Tokens, varScope *varScope, funcScope *funcScope) {
	tokens.Consume()
	condition := tokens.ConsumeTuple()

	if len(condition) > 1 {
		throw("Expected only one condition in if condition", tokens)
	} else if len(condition) == 0 {
		throw("Expected a condition in if block", tokens)
	}
}
