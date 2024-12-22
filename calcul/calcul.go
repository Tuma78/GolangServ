package calсul

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type Stack []string

func (s *Stack) push(v string) {
	*s = append(*s, v)
}


func (s *Stack) pop() (string, error) {
	if len(*s) == 0 {
		return "", errors.New("stack is empty")
	}
	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element, nil
}


func (s *Stack) stackEmpty() bool {
	return len(*s) == 0
}

func priority(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	case "(":
		return 0
	case ")":
		return 0
	default:
		return -1 
	}
}

func splitString(str string) []string {
	var tokens []string
	var currentToken strings.Builder

	for _, char := range str {
		if unicode.IsDigit(char) || char == '.' {
			currentToken.WriteRune(char)
		} else {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			if char != ' ' {
				tokens = append(tokens, string(char))
			}
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}

func ToRPN(str string) (string, error) {
	var current string
	var opStack Stack
	context := splitString(str)

	for _, element := range context {
		switch priority(element) {
		case -1: 
			current += element + " "
		case 1, 2: 
			for !opStack.stackEmpty() && priority(opStack[len(opStack)-1]) >= priority(element) {
				op, _ := opStack.pop()
				current += op + " "
			}
			opStack.push(element)
		case 0: 
			if element == "(" {
				opStack.push(element)
			} else if element == ")" {
				for !opStack.stackEmpty() && opStack[len(opStack)-1] != "(" {
					op, _ := opStack.pop()
					current += op + " "
				}
				opStack.pop() 
			}
		}
	}

	for !opStack.stackEmpty() {
		op, _ := opStack.pop()
		current += op + " "
	}

	return strings.TrimSpace(current), nil
}


func Calc(expression string) (float64, error) {
	rpn, err := ToRPN(expression)
	if err != nil {
		return 0, err
	}

	tokens := strings.Fields(rpn)
	var stack Stack

	for _, token := range tokens {
		if priority(token) == -1 { 
			stack.push(token)
		} else { 
			if len(stack) < 2 {
				return 0, errors.New("invalid expression")
			}
			operand2, _ := stack.pop()
			operand1, _ := stack.pop()

			num1, err := strconv.ParseFloat(operand1, 64)
			if err != nil {
				return 0, err
			}
			num2, err := strconv.ParseFloat(operand2, 64)
			if err != nil {
				return 0, err
			}

			var result float64
			switch token {
			case "+":
				result = num1 + num2
			case "-":
				result = num1 - num2
			case "*":
				result = num1 * num2
			case "/":
				if num2 == 0 {
					return 0, errors.New("division by zero")
				}
				result = num1 / num2
			default:
				return 0, errors.New("invalid operator")
			}

			stack.push(strconv.FormatFloat(result, 'f', -1, 64))
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid expression")
	}
	result, err := stack.pop()
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(result, 64)
}