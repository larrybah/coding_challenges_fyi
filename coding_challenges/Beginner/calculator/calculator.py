#!/usr/bin/env python3

"""
Accept user input and calculate result base on the operator provided.
"""

first_number = int(input("Enter first number: "))
operator = input("Enter operator(+,-,/,*): ") 
second_number = int(input("Enter second number: "))

operators = ['+', '-', '*', '/']

for operator in operators:
    if operator == '+':
        result = first_number + second_number
        print(f"{first_number} + {second_number} = {result}")
    elif operator == '-':
        result = first_number - second_number
        print(f"{first_number} + {second_number} = {result}")
    elif operator == '/':
        result = first_number / second_number
        print(f"{first_number} / {second_number} = {result}")
    elif operator == '*':
        result = first_number * second_number
        print(f"{first_number} * {second_number} = {result}")
    else:
        print("Enter a valid Operator!")
