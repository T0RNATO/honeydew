package main

import (
	"fmt"
	"runtime"
	"strings"
)

var reset = "\033[0m"

var red = "\033[31m"
var blue = "\033[34m"
var cyan = "\033[0;36m"

var grey = "\033[30m"

func logS(s string) {
	fmt.Print(blue + "\"")
	fmt.Print(s)
	fmt.Print("\"" + reset)
}

func logNum(i any) {
	fmt.Print(red)
	fmt.Print(i)
	fmt.Print(reset)
}

func logB(b bool) {
	fmt.Print(cyan)
	fmt.Print(b)
	fmt.Print(reset)
}

func logList[T any](s []T, f func(a T)) {
	fmt.Print("[")
	for _, item := range s {
		f(item)
		fmt.Print(", ")
	}
	fmt.Print("]")
}

func logItem(a any) {
	switch i := a.(type) {
	case []string:
		logList(i, logS)
	case [][]string:
		logList(i, func(a []string) {
			logList(a, logS)
		})
	case []int, []float32, []float64:
		logList(i.([]any), logNum)
	case []any:
		logList(i, logItem)
	case string:
		logS(i)
	case int, float32, float64:
		logNum(i)
	case bool:
		logB(i)
	case Tokens:
		logList(i.tokens, logS)
		fmt.Printf(" @ nextread=%d", i.index)
	}
	fmt.Print(" ")
}

func Log[T any](a T, b ...any) T {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		folders := strings.Split(file, "/")
		fmt.Print(grey)
		fmt.Printf("[%s:%d] ", folders[len(folders)-1], line)
		fmt.Print(reset)
	}
	logItem(a)
	for _, item := range b {
		logItem(item)
	}
	fmt.Print("\n")

	return a
}
