package main

import (
	"strings"
)

const SPLITTING_CHARS = " ;(){}'\":+*-/,"

type Tokens struct {
	tokens          []string
	index           int
	collectionIndex int
}

func (self *Tokens) Consume() string {
	token := self.tokens[self.index]
	self.index += 1
	return token
}

func (self *Tokens) NextConsumed() string {
	return self.tokens[self.index]
}

func (self *Tokens) IsEmpty() bool {
	return len(self.tokens) == self.index
}

func Tkns(tokens []string) Tokens {
	return Tokens{tokens, 0, 0}
}

func (self *Tokens) StartCollecting() {
	self.collectionIndex = self.index
}

func (self *Tokens) Collected() []string {
	return self.tokens[self.collectionIndex:self.index]
}

func (self *Tokens) CollectedMinusOne() []string {
	return self.tokens[self.collectionIndex : self.index-1]
}

func (self *Tokens) ConsumeTuple() [][]string {
	var output [][]string

	depth := 0
	assertToken(self.Consume(), "(")
	self.StartCollecting()
	for {
		token := self.Consume()
		if token == "(" {
			depth += 1
		} else if token == ")" {
			if depth == 0 {
				output = append(output, self.CollectedMinusOne())
				break
			}
			depth -= 1
		} else if token == "," && depth == 0 {
			output = append(output, self.CollectedMinusOne())
			self.StartCollecting()
		}

	}
	return output
}

func (self *Tokens) ConsumeUntilEndBlock() []string {
	self.StartCollecting()
	depth := -1
	for {
		token := self.Consume()
		if token == "{" {
			depth += 1
		} else if token == "}" {
			if depth == 0 {
				break
			}
			depth -= 1
		}
	}
	return self.Collected()
}

func (self *Tokens) ConsumeCurlyBrackets() []string {
	assertToken(self.Consume(), "{")
	self.StartCollecting()
	depth := 0
	for {
		token := self.Consume()
		if token == "{" {
			depth += 1
		} else if token == "}" {
			if depth == 0 {
				break
			}
			depth -= 1
		}
	}
	return self.CollectedMinusOne()
}

func (self *Tokens) ConsumeLine() []string {
	self.StartCollecting()
	for v := self.Consume(); v != ";"; v = self.Consume() {
	}
	return self.Collected()
}

func Tokenise(s string) Tokens {
	tokens := []string{}
	var chars strings.Builder

	in_comment := false

	for _, char := range s {
		if char == '#' {
			in_comment = true
			continue
		} else if in_comment {
			if char == '\n' {
				in_comment = false
			}
		} else if strings.ContainsRune(SPLITTING_CHARS, char) {
			if chars.Len() > 0 {
				tokens = append(tokens, chars.String())
				chars.Reset()
			}
			if char != ' ' {
				tokens = append(tokens, string(char))
			}
		} else if char != '\n' {
			chars.WriteRune(char)
		}
	}

	if chars.Len() > 0 {
		tokens = append(tokens, chars.String())
	}

	return Tkns(tokens)
}
