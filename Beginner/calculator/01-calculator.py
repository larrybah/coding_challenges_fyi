#!/usr/bin/env python3

"""
First Attempt:
Accept user input and calculate result base on the operator provided.
"""

first_number = int(input("Enter first number: "))
operator = input("Enter operator(+,-,/,*): ") 
second_number = int(input("Enter second number: "))

if operator == '+':
    result = first_number + second_number
    print(f"{first_number} + {second_number} = {result}")
elif operator == '-':
    result = first_number - second_number
    print(f"{first_number} + {second_number} = {result}")
elif operator == '/':
    if second_number == 0:
        print("Cannot Divide by Zero!")
    else:
        result = first_number / second_number
        print(f"{first_number} / {second_number} = {result}")
elif operator == '*':
    result = first_number * second_number
    print(f"{first_number} * {second_number} = {result}")
else:
    print("Enter a valid Operator(+,-,*,/)!")
