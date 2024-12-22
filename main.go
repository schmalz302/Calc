package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

func Calc(expression string) (float64, error) {
	postfix, err := toPostfix(expression)
	if err != nil {
		return 0, err
	}

	result, err := evalPostfix(postfix)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func toPostfix(expression string) ([]string, error) {
	var output []string
	var stack []rune
	precedence := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	expression = strings.ReplaceAll(expression, " ", "")

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])

		if unicode.IsDigit(char) {
			num := string(char)
			for i+1 < len(expression) && unicode.IsDigit(rune(expression[i+1])) {
				i++
				num += string(expression[i])
			}
			output = append(output, num)
		} else if char == '(' {
			stack = append(stack, char)
		} else if char == ')' {
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				output = append(output, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return nil, errors.New("mismatched parentheses")
			}
			stack = stack[:len(stack)-1]
		} else if precedence[char] > 0 {
			for len(stack) > 0 && stack[len(stack)-1] != '(' && precedence[stack[len(stack)-1]] >= precedence[char] {
				output = append(output, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, char)
		} else {
			return nil, errors.New("invalid character in expression")
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == '(' {
			return nil, errors.New("mismatched parentheses")
		}
		output = append(output, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func evalPostfix(postfix []string) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, errors.New("invalid expression")
			}

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

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Expression is not valid", http.StatusUnprocessableEntity)
		return
	}

	result, err := Calc(req.Expression)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "mismatched") {
			writeErrorResponse(w, "Expression is not valid", http.StatusUnprocessableEntity)
		} else {
			writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := Response{Result: fmt.Sprintf("%v", result)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
