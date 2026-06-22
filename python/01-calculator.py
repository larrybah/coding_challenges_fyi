#!/usr/bin/env python3

def calculator(num1, operator, num2):
    """Perform calculation based on operator"""
    if operator == '+':
        return num1 + num2
    elif operator == '-':
        return num1 - num2
    elif operator == '*':
        return num1 * num2
    elif operator == '/':
        if num2 == 0:
            return "Error: Division by zero"
        return num1 / num2
    else:
        return "Error: Invalid operator"

def main():
    print("=== CLI Calculator ===")

    try:
        """Get first number"""
        num1 = float(input("Enter first number: "))

        """ Get Operator """
        operator = input("Enter operator (+ - * /): ").strip()

        """ Validate operator """

        if operator not in ['+', '-', '*', '/']:
            print ("Error: Invalid Operator. Please use +, -, *, /")
            return

        """ Get number 2 """
        num2 =  float(input("Enter second number: "))

        result = calculator (num1, operator, num2)
        print(f"\nResult: {num1} {operator} {num2} = {result}")

    except ValueError:
        print("Error: Please enter a valid numbers")
    except keyboardInterrupt:
        print("\n\nCalculator closed.")

if __name__ == "__main__":
    main()
