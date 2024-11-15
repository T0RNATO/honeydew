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

func main() {
	file, err := os.Open(os.Args[1])
	e(err)

	contents, err := io.ReadAll(file)
	e(err)

	Run(string(contents), &globalVars, &globalFuncs)
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
				panic("Cannot use return statement at top level") // todo
			}
			return Return(tokens.ConsumeLine(), varScope, funcScope)
		default:
			tokens.StartCollecting()
			name := token
			if _, exists := (*varScope)[name]; exists {
				VarUpdate(tokens.ConsumeLine(), varScope, funcScope)
			} else if _, exists := (*funcScope)[name]; exists {
				FuncCall(tokens.ConsumeLine(), varScope, funcScope)
			} else {
				panic(fmt.Sprintf("Unknown symbol '%s'", token)) // todo
			}
		}
	}
	return nil
}

func Run(s string, varScope *varScope, funcScope *funcScope) {
	RunTokens(Tokenise(s), varScope, funcScope)
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}
