package ruleparser

import "go-backend/api/server/tools/util"

func TransformCondition(infix []Token) []Token {
	var suffix []Token
	var optSt util.Stack
	var postSt util.Stack
	for _, token := range infix {
		t := token.TokenType
		v := token.TokenValue
		if t == NumTokenType || t == ValTokenType {
			postSt.Push(token)
		} else if t == PairTokenType {
			if v == "(" {
				optSt.Push(token)
			} else {
				found := false
				for !optSt.IsEmpty() {
					if optSt.Top().(Token).TokenValue != "(" {
						postSt.Push(optSt.Pop())
					} else {
						optSt.Pop()
						found = true
					}
				}
				if !found {
					panic("Syntax Error: Pair not matched!")
				}
			}
		} else if t == OptTokenType {
			if optSt.IsEmpty() {
				optSt.Push(token)
			} else if optSt.Top().(Token).TokenType == OptTokenType {
				if IsHigherPriority(token.TokenValue, optSt.Top().(Token).TokenValue) {
					postSt.Push(token)
				} else {
					postSt.Push(optSt.Pop())
					optSt.Push(token)
				}
			} else {
				optSt.Push(token)
			}
		}
	}
	for !optSt.IsEmpty() {
		postSt.Push(optSt.Pop())
	}
	for !postSt.IsEmpty() {
		suffix = append([]Token{postSt.Pop().(Token)}, suffix...)
	}
	return suffix
}

func IsHigherPriority(optA string, optB string) bool {
	return getPriorityNum(optA) < getPriorityNum(optB)
}

func getPriorityNum(opt string) int {
	switch opt {
	case "*":
	case "/":
		return 1
	case "+":
	case "-":
		return 2
	case ">":
	case "<":
	case "!=":
	case "==":
		return 3
	}
	return 4
}
