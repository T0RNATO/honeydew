package main

import (
	"fmt"
	"io"
	"os"
)

var (
	globalVars  = make(varScope)
	globalFuncs = make(funcScope)
)

var fileName string
var fileTokens Tokens

func main() {
	fileName = os.Args[1]
	file, err := os.Open(fileName)
	e(err)

	contents, err := io.ReadAll(file)
	e(err)

	fileTokens = Tokenise(string(contents))

	RunTokens(fileTokens, &globalVars, &globalFuncs)
}

func RunTokens(tokens Tokens, varScope *varScope, funcScope *funcScope) any {
	InitSTD(funcScope)

	for !tokens.IsEmpty() {
		// Look at first token of the line to decide what to do
		token := tokens.NextConsumed()
		switch token {
		case "var":
			VarDefinition(tokens.ConsumeLine(), varScope, funcScope)
		case "function":
			FuncDefinition(tokens.ConsumeUntilEndBlock(), varScope, funcScope)
		case "return":
			if varScope == &globalVars {
				throw("Cannot use return statement at top level", tokens)
			}
			return Return(tokens.ConsumeLine(), varScope, funcScope)
		case "if":
			IfBlock(tokens.ConsumeUntilEndBlock(), varScope, funcScope)
		default:
			tokens.StartCollecting()
			name := token
			if _, exists := (*varScope)[name]; exists {
				VarUpdate(tokens.ConsumeLine(), varScope, funcScope)
			} else if _, exists := (*funcScope)[name]; exists {
				FuncCall(tokens.ConsumeLine(), varScope, funcScope)
			} else {
				throw(fmt.Sprintf("Unknown symbol '%s'", token), tokens)
			}
		}
	}
	return nil
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}
