# GoLox

An implementation of Bob Nystrom's Lox programming language from [Crafting Interpreters](https://craftinginterpreters.com/).

Lox is a dynamic, high-level scripting language.

The [examples](/examples) directory has a few basic example programs in Lox with equivalent Python programs for comparison.

# Usage

Run `go build` to build `golox`.

Run `./golox` without any arguments to enter the REPL.

Run `./golox <filename>.lox` to run a lox file.

# Features

## REPL
 - [x] Pretty-print parsed program
 - [x] Pretty-print local variables after a command
 - [ ] Support raw keyboard mode 
    - [ ] up and down arrow keys to go to previous commands 
    - [ ] allow newlines for multiline REPL programs

### Interpreter

- Comments: `// this is a line comment`
- Data Types
    - [x] Numbers (represented by float64): `1.535`, `32`, etc.
    - [x] Booleans: `true`, `false`
    - [x] Strings: `"this is a string"`
    - [x] Nil: `nil`
    - [ ] *(Extension)* Lists
- Expressions
    - Arithmetic
        - [x] Addition: `18.9 + 16.3`
        - [x] Subtraction: `18.9 - 16.3`
        - [x] Multiplication: `18.9 * 16.3`
        - [x] Division: `18.9 / 16.3`
        - [x] Negation: `-18.9`, `-(10 * 2)`
    - Comparison & Equality
        - [x] Less Than: `18.9 < 16.3` is `false`
        - [x] Less Than or Equal: `18.9 <= 16.3` is `false`
        - [x] Greater Than: `18.9 > 16.3` is `true`
        - [x] Greater Than or Equal: `18.9 >= 16.3` is `true`
    - Logical Operators
        - [x] Not: `!false`
        - [ ] And: `true and false` is `false`
        - [ ] Or: `true or false` is `true`
    - [x] Precedence and Grouping: `(2 + 3 * 4) / 2` is `7`
    - [ ] String Concatenation: `"hey" + " " + "there"` is `"hey there"`
    - [ ] *(Extension)* Lists Concatenation
- Statements
    - [x] Print Statements: 
        ```
        print "hello";
        print 1.84; 
        print x;
        ```
    - [x] Expression Statements: 
        ```
        "hello";
        1.84;
        x;
        ```
    - [x] If/Else Statements:
        ```
        if (10 > 5) {
            x = 10;
        }

        if (x > 12) {
            x = 12;
        } else {
            x = x + 1;
        }
        ```
    - [x] While Loops: 
        ```
        while (x < 12) {
            x = x + 1;
        }
        ```
    - [x] Variable Declarations: 
        ```
        var x = 103;
        var foo123 = "hello";
        ```
    - [x] Variable Assignments: 
        ```
        x = 103;
        foo123 = "hello";
        ```
    - [x] Block Statements: 
        ```
        var y = 11;
        {
            var x = 10 + y; 
            print x;
        }
        // x doesn't exist outside of the block scope
        ```
    - [x] Function Declaration & Calls: 
        - [x] Close over outer-scoped variables
        - [x] Recursion
        - [x] Return Statements exit scope and return value 
        ```
        var foo = 15;
        fun myFunc(x, y, z) {
            if (myFunc(x+1, y-1, z) > 10) {
                return 10;
            }
            return foo + x * y * z;
        }
        print myFunc(1, 4, 2);
        ```
    - [ ] Classes
        - [ ] Class Declaration & Instantiation: 
        - [ ] Class Methods and Properties
        ```
        class BaseClass {
            sayHi() { print "Hi!"; }
        }
        class Foo < BaseClass { 
            init(meat) {
                this.meat = meat;
            }
            cook() { print "Eggs and " + this.meat + " cooking!"; } 
            serve(customer) {print "Here's your order, " + customer;} 
        }

        // Instantiate and call method
        var foo = Foo("bacon");
        foo.serve("Billy");

        // Set arbitrary property value
        foo.whateverProperty = "Look, a new property!";

        // Call inherited method from base class
        foo.sayHi();
        ```
### Extensions
- [ ]  Standard Library 
- [ ]  Lists
- [ ]  Custom Garbage Collector (currently piggybacking on Go's GC)
- [ ]  Compile to bytecode or machine code instead of interpreting AST
    - Nothing done yet, but I'm interested in learning to use LLVMjit to make the language fast and efficient.

## Useful Resources

1. [Crafting Interpreters](https://craftinginterpreters.com/)
2. [Writing An Interpreter In Go](https://interpreterbook.com/)
3. [LLVM Usage For Creating a Programming Language](https://mukulrathi.com/create-your-own-programming-language/llvm-ir-cpp-api-tutorial/)
4. [LLVM C API Tutorial](https://www.pauladamsmith.com/blog/2015/01/how-to-get-started-with-llvm-c-api.html)
5. [LLVM Library for Go](https://llir.github.io/document/user-guide/basic/)
