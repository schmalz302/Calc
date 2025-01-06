package math

import (
	"errors"
	"strconv"
	"unicode"
)

// Calc вычисляет значение математического выражения в строке.
func Calc(expression string) (float64, error) {
	if len(expression) == 0 {
		return 0, errors.New("invalid expression")
	}
	// Преобразуем выражение в постфиксную нотацию
	postfix, err := toPostfix(expression)
	if err != nil {
		return 0,  errors.New("invalid expression")
	}

	// Вычисляем значение выражения в постфиксной нотации
	result, err := evalPostfix(postfix)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// toPostfix преобразует инфиксное выражение в постфиксную (обратную польскую) нотацию.
func toPostfix(expression string) ([]string, error) {
	var output []string
	var stack []rune
	precedence := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])

		if unicode.IsDigit(char) {
			// Чтение числа
			num := string(char)
			for i+1 < len(expression) && unicode.IsDigit(rune(expression[i+1])) {
				i++
				num += string(expression[i])
			}
			output = append(output, num)
		} else if char == '(' {
			stack = append(stack, char)
		} else if char == ')' {
			// Выгружаем все операторы до открывающей скобки
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				output = append(output, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return nil, errors.New("mismatched parentheses")
			}
			stack = stack[:len(stack)-1] // удаляем '('
		} else if char == '+' || char == '-' || char == '*' || char == '/' {
			// Выгружаем операторы с более высоким приоритетом
			for len(stack) > 0 && stack[len(stack)-1] != '(' && precedence[stack[len(stack)-1]] >= precedence[char] {
				output = append(output, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, char)
		} else if !unicode.IsSpace(char) {
			return nil, errors.New("invalid character in expression")
		}
	}

	// Перенос оставшихся операторов в выходную очередь
	for len(stack) > 0 {
		if stack[len(stack)-1] == '(' {
			return nil, errors.New("mismatched parentheses")
		}
		output = append(output, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

// evalPostfix вычисляет значение выражения в постфиксной нотации.
func evalPostfix(postfix []string) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else {
			// Обработка операторов
			if len(stack) < 2 {
				return 0, errors.New("invalid expression")
			}

			// Извлекаем два операнда
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, errors.New("division by zero")
				}
				stack = append(stack, a/b)
			default:
				return 0, errors.New("invalid operator")
			}
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid expression")
	}

	return stack[0], nil
}
