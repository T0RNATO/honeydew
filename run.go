package main

import "fmt"

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
	body []string
}

func (self *UserFunction) run(args []any, vScope *varScope, fScope *funcScope) any {
	var blockVarScope = make(varScope)
	var blockFuncScope = make(funcScope)

	if len(args) != len(self.args) {
		panic("Incorrect number of arguments: " + string(len(args))) // todo
	}

	for k, v := range *vScope {
		blockVarScope[k] = v
	}

	for i, arg := range args {
		blockVarScope[self.args[i].Name] = arg.(float32) // todo, temp
	}

	for k, v := range *fScope {
		blockFuncScope[k] = v
	}

	return RunTokens(Tkns(self.body), &blockVarScope, &blockFuncScope)
}

type STDFunction struct {
	args []Arg
	// returns string // todo
	Func func(args ...any) any
}

func (self *STDFunction) run(args []any, varScope *varScope, funcScope *funcScope) any {
	return self.Func(args...)
}

func assertToken(t string, expected string) {
	if t != expected {
		panic(fmt.Sprintf("Got '%s', expected '%s'", t, expected)) // todo
	}
}

func VarDefinition(tokens []string, varScope *varScope, funcScope *funcScope) {
	assertToken(tokens[2], "=")
	(*varScope)[tokens[1]] = ParseExpression(tokens[3:], varScope, funcScope)
}
func FuncDefinition(tkns []string, varScope *varScope, funcScope *funcScope) {
	tokens := Tkns(tkns[1:])
	name := tokens.Consume()
	assertToken(tokens.tokens[1], "(")

	allArgTokens := tokens.ConsumeTuple()
	var args []Arg

	for _, argTokens := range allArgTokens {
		argName := argTokens[0]
		assertToken(argTokens[1], ":")
		args = append(args, Arg{argName, argTokens[2]})
	}

	funcBody := tokens.ConsumeCurlyBrackets()

	(*funcScope)[name] = &UserFunction{args, funcBody}
}
func FuncCall(tkns []string, varScope *varScope, funcScope *funcScope) any {
	tokens := Tkns(tkns)
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
		panic(fmt.Sprintf("Unknown symbol '%s'", name)) // todo
	}
}
func VarUpdate(tokens []string, varScope *varScope, funcScope *funcScope) {
	assertToken(tokens[1], "=")
	(*varScope)[tokens[0]] = ParseExpression(tokens[2:], varScope, funcScope)
}
func Return(tokens []string, varScope *varScope, funcScope *funcScope) float32 {
	return ParseExpression(tokens[1:], varScope, funcScope)
}
