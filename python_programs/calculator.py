#!/usr/bin/env python3

"""
A Simple Calculator program that handles basic arithmetic.
"""

first_number = int(input("Enter first number: "))
operator = input("Enter the Operator (+, -, /, *): ")
second_number = int(input("Enter second number: "))

if operator == '+':
    result = first_number + second_number
    print(f"Result: {result}")
elif operator == '-':
    result = first_number - second_number
    print(f"Result: {result}")
elif operator == '/':
    result = first_number / second_number
    print(f"Result: {result}")
elif operator == '*':
    result = first_number * second_number
    print(f"Result: {result}")
