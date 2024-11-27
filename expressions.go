package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const OPERATIONS = "+-*/"

type stack[T any] []T

func (self *stack[T]) pop() T {
	lastIndex := len(*self) - 1
	item := (*self)[lastIndex]
	*self = (*self)[:lastIndex]
	return item
}

func (self *stack[T]) append(item T) {
	*self = append(*self, item)
}

func ParseExpression(tokens Tokens, varScope *varScope, funcScope *funcScope) float32 {
	var queue []any
	{
		var stack stack[string]

		for i, token := range tokens.tokens {
			asFloat, isntNum := strconv.ParseFloat(token, 32)
			switch {
			case token == "(":
				stack.append(token)
			case token == ")":
				l := len(stack)
				for l > 0 && stack[l-1] != "(" {
					queue = append(queue, stack.pop())
					l = len(stack)
				}
				stack.pop()
			case strings.Contains(OPERATIONS, token):
				l := len(stack)
				for l > 0 && (stack[l-1] == "*" || stack[l-1] == "/") {
					queue = append(queue, stack.pop())
					l = len(stack)
				}
				stack.append(token)
			case isntNum == nil:
				queue = append(queue, float32(asFloat))
			default:
				if value, ok := (*varScope)[token]; ok {
					queue = append(queue, value)
				} else if _, ok := (*funcScope)[token]; ok {
					queue = append(queue, FuncCall(tokens.Slice(i, -1), varScope, funcScope))
				}
			}
		}

		for _, item := range slices.Backward(stack) {
			queue = append(queue, item)
		}
	}

	var stack stack[float32]

	for _, item := range queue {
		switch i := item.(type) {
		case float32:
			stack = append(stack, i)
		case string:
			switch i {
			case "+":
				stack.append(stack.pop() + stack.pop())
			case "*":
				stack.append(stack.pop() * stack.pop())
			case "-":
				a := stack.pop()
				stack.append(stack.pop() - a)
			case "/":
				a := stack.pop()
				stack.append(stack.pop() / a)
			default:
				panic(fmt.Sprintf("Internal Error: Unexpected operation '%s'", i))
			}
		}
	}

	return stack[0]
}
