// Package main demonstrates a simple command-line calculator
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// add returns the sum of two numbers
func add(a, b float64) float64 {
	return a + b
}

// subtract returns the difference of two numbers (a - b)
func subtract(a, b float64) float64 {
	return a - b
}

// multiply returns the product of two numbers
func multiply(a, b float64) float64 {
	return a * b
}

// divide returns the quotient of two numbers (a / b)
// Returns 0 if b is zero to avoid division by zero errors
func divide(a, b float64) float64 {
	if b == 0 {
		fmt.Println("Error: Cannot divide by zero")
		return 0
	}
	return a / b
}

// performCalculation takes two numbers and an operator, then returns the result
func performCalculation(num1, num2 float64, operator string) float64 {
	switch operator {
	case "+":
		return add(num1, num2)
	case "-":
		return subtract(num1, num2)
	case "*":
		return multiply(num1, num2)
	case "/":
		return divide(num1, num2)
	default:
		fmt.Println("Error: Invalid operator. Use +, -, *, or /")
		return 0
	}
}

// main is the entry point of the program
func main() {
	// Create a scanner to read user input from standard input
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("===== Simple Calculator =====")
	fmt.Println("Supported operations: +, -, *, /")
	fmt.Println("Enter calculation in format: number operator number")
	fmt.Println("Example: 10 + 5")
	fmt.Println("Type 'quit' to exit")
	fmt.Println("=============================")

	// Main loop to continuously accept calculations
	for {
		fmt.Print("\nEnter calculation (or 'quit' to exit): ")

		// Read user input
		scanner.Scan()
		input := scanner.Text()

		// Check if user wants to exit
		if strings.ToLower(input) == "quit" {
			fmt.Println("Thank you for using the calculator. Goodbye!")
			break
		}

		// Split the input by spaces
		parts := strings.Fields(input)

		// Validate that we have exactly 3 parts: number, operator, number
		if len(parts) != 3 {
			fmt.Println("Error: Invalid format. Please use: number operator number")
			continue
		}

		// Parse the first number
		num1, err1 := strconv.ParseFloat(parts[0], 64)
		if err1 != nil {
			fmt.Printf("Error: '%s' is not a valid number\n", parts[0])
			continue
		}

		// Parse the second number
		num2, err2 := strconv.ParseFloat(parts[2], 64)
		if err2 != nil {
			fmt.Printf("Error: '%s' is not a valid number\n", parts[2])
			continue
		}

		// Extract the operator
		operator := parts[1]

		// Perform the calculation
		result := performCalculation(num1, num2, operator)

		// Display the result
		fmt.Printf("Result: %g %s %g = %g\n", num1, operator, num2, result)
	}
}
