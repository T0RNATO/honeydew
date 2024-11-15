package main

import "fmt"

func InitSTD(funcScope *funcScope) {
	(*funcScope)["print"] = &STDFunction{nil,
		func(args ...any) any {
			l, _ := fmt.Println(args...)
			return l
		},
	}
}
