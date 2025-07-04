#!/usr/bin/env python3

"""
A simple calculator program.
"""

def get_numbers():
    """Get two numbers from the user."""
    while True:
        try:
            first_number = int(input("Enter first number: "))
            second_number = int(input("Enter second number: "))
            return first_number, second_number
        except ValueError:
            print("Invalid input. Please enter numbers only.")

def get_operator():
    """Get a valid operator from the user."""
    while True:
        operator = input("Operator: +, -, *, /: ")
        if operator in ['+', '-', '*', '/']:
            return operator
        else:
            print("Invalid operator. Try using (+, -, *, /)")

def calculate(first_number, operator, second_number):
    """Perform the calculation."""
    if operator == '+':
        return first_number + second_number
    elif operator == '-':
        return first_number - second_number
    elif operator == '*':
        return first_number * second_number
    elif operator == '/':
        if second_number == 0:
            return "Division by zero is not allowed!"
        else:
            return first_number / second_number

def main():
    """Main function to run the calculator."""
    first_number, second_number = get_numbers()
    operator = get_operator()
    result = calculate(first_number, operator, second_number)
    print(f"Result: {result}")

if __name__ == "__main__":
    main()
