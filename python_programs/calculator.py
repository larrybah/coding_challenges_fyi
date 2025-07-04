#!/usr/bin/env python3

"""
A simple calculator program.
"""

first_number = int(input("Enter first number: "))

operator = input("Operator: +, -, *, /: ")

second_number = int(input("Enter second number: "))

result = 0

if operator == '+':
    result = first_number + second_number
    print(result)
elif operator == '-':
    result = first_number - second_number
    print(result)
elif operator == '*':
    result = first_number * second_number
    print(result)
elif operator == '/':
    if first_number == 0 or second_number == 0:
        print("Division by Zero is not allowed!")
    else:
        result = first_number / second_number
        print(result)
