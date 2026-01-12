#!/usr/bin/env python3

import ast
import operator

"""
Defines a mapping of AST node types to their corresponding  operator functions
This allows us to safely evaluate mathematical expressions without using eval()
"""

OPERATORS = {
        ast.Add: operator.add,      # Addition (+)
        ast.Sub: operator.sub,      # Subtraction (-)
        ast.Mult: operator.mul,     # Multiplication (*)
        ast.Div: operator.truediv,  # Division (/)
        ast.Pow: operator.pow,      # Exponentiation (**)
        ast.USub: operator.neg,     # Unary negation (-)

    }

def safe_eval(expression):
    """
    Safely evaluate a mathematical expression using AST parsing.

    This function parses the expression into an Abstract Syntax Tree (AST)
    and recursively evaluates it. This is much safer than using eval()
    because we explicitly control which operations are allowed.

    Args:
        expression: String containing the mathematical expression

    Returns:
        ValueError: If the expression contains unsupported operations
        SyntaxError: If the expression has invalid sytax
    """

    try:
        # Parse the expression into an AST
        # mode='eval' tells the parser we expecting a single expression
        
        node = ast.parse(expression, mode='eval').body
        return _eval_node(node)
    except SyntaxError as e:
        raise ValueError(f"Invalid Syntax: {e}")


