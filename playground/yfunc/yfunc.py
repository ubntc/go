import sys
from typing import Callable, Type

FactorialFuncType = Callable[[int], int]
YCombinatorType = Callable[[FactorialFuncType], FactorialFuncType]

# Implement Y
# ===========

# Explicit Y combinator with less lambdas
def Y(f: YCombinatorType) -> FactorialFuncType:
    def self_apply(combinator: YCombinatorType) -> FactorialFuncType:
        return combinator(combinator)

    def generate_recursive_func(recursive_func: YCombinatorType) -> FactorialFuncType:
        return f(lambda iteration: recursive_func(recursive_func)(iteration))

    return self_apply(generate_recursive_func)

# Recursion Use Case
# ==================

def factorial_generator(f: FactorialFuncType) -> FactorialFuncType:
    return lambda n: 1 if n == 0 else n * f(n - 1)

factorial = Y(factorial_generator)


# Other Implementations of Y
# ==========================

# Basic Y combinator with lambdas
def Y_shortest(f: YCombinatorType) -> FactorialFuncType:
    return (lambda x: x(x))(lambda y: f(lambda it: y(y)(it)))

# Basic Y combinator with lambdas and named vars
def Y_lambda(factorial_generator: YCombinatorType) -> FactorialFuncType:
    return (
        # 1. Define and call a function that takes another function `fn` as an argument and applies the
        #    given function to itself. This sets up the stage for the recursion in the following steps.
        lambda fn: fn(fn)
    )(
        # 2. Define the core recursive handler. Use `factorial_generator` and provide it with a function
        #    capable of recursion. This sets up the recursive function call structure.
        lambda recursive_handler: factorial_generator(
            # 3. Define a function that takes an argument `iteration`, representing the current value
            #    being processed. In this case `iteration` is an integer number.
            #    This function invokes `recursive_handler` with itself, creating the recursion needed to calculate the factorial.
            lambda iteration: recursive_handler(recursive_handler)(iteration)
        )
    )

class YCombinator:
    def __init__(self, f: YCombinatorType):
        self.f = f

    def self_apply(self, instance) -> FactorialFuncType:
        return instance.generate_recursive_func(instance)

    def generate_recursive_func(self, instance) -> FactorialFuncType:
        return self.f(lambda iteration: instance.generate_recursive_func(instance)(iteration))

    def Y(self) -> FactorialFuncType:
        return self.self_apply(self)

factorial_generator_instance = YCombinator(factorial_generator)
factorial = factorial_generator_instance.Y()

def run():
    for i in [5, 10, 100, 200, 500, 1000]:
        name = f"factorial({i})"
        try:
            result = factorial(i)
            size = len(f"{result}")
            if size > 100:
                result = "omitted"
            result = f"value:{result} (digits:{size})"
        except Exception as ex:
            result = f"error:{ex}"
            size = "N/A"

        print(f"{name} = {result}")

if __name__ == '__main__': run()
