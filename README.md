# GoForth
A Forth-like interpreter written in Go for [Coding Challenges](https://codingchallenges.fyi/challenges/intro/).

## Forth

Forth is a stack-based programming language that uses space-separated commands, known as words. Words 
manipulate the stack, which stores integers.

### Example Programs

A simple calculation 2 + 2:

```forth
2 2 +
```

Comments are any content between an open bracket and a close bracket: `( this is a comment )`.

```forth
: fib over over + ;  ( generate the next number in the Fibonacci sequence ) 
1 1 fib fib fib fib fib fib fib fib fib ( generate the first few numbers the Fibonacci sequence )
( TODO put it in a loop make another word to do fibN )
```

### Language Support

Most Forth words affect the stack in some way. Some take values off the stack, some leave new values on the stack, 
and some do a mixture of both. These “stack effects” are commonly represented using comments of the form 
( before -- after ). For example, + is ( n1 n2 -- sum ) where n1 and n2 are the top two numbers on the stack which are 
popped off, and the result of adding them together, the sum, is the value left on the stack.

For example if the stack contains: `1 2 3` where 1 is the bottom of the stack and 3 is the top of the stack, the word 
`+` pops n1 (2) and n2 (3) off the stack and then adds them, pushing the sum onto the top of the stack. 

This interpreter supports the following built-in words:

| Word | Stack Effects            | Description                                                                                             |
|------|--------------------------|---------------------------------------------------------------------------------------------------------|
| bye  | no effect                | Exits the interpreter                                                                                   |
| +    | ( n1 n2 -- sum )         | Pops the top two elements on the stack, pushes the sum on to the top of the stack                       |
| -    | ( n1 n2 -- diff )        | Pops the top two elements on the stack, substracts n2 from n1 stores the result on the top of the stack |
| *    | ( n1 n2 -- multiplied )  | Pops the top two elements on the stack, pushes the product on to the top of the stack                   |
| /    | ( n1 n2 -- divided )     | Pops the top two elements on the stack, pushes the result of n2 / n1 on to the top of the stack         |
| mod  | ( n1 n2 -- modulus )     | Pops the top two elements on the stack, pushes the remainder of n2 / n1 on to the top of the stack      |
| swap | ( n1 n2 -- n2 n1 )       | Swaps the top two elements on the stack                                                                 |
| dup  | ( n -- n n )             | Duplicates the top element on the stack                                                                 |
| over | ( n1 n2 -- n1 n2 n1 )    | Duplicates the second from top element and pushes it on to the top of the stack                         |
| rot  | ( n1 n2 n3 -- n2 n3 n1 ) | Rotates the top three elements on the stack                                                             |
| drop | ( n1 -- )                | Pops the top element off the stack                                                                      |
| .    | ( n1 -- )                | Prints and pops the top of the stack                                                                    |
| emit | ( n1 -- )                | Prints the top of the stack as n ASCII character and pops the top of the stack                          |
| cr   | ( -- )                   | Prints a newline                                                                                        |
| ."   | ( -- )                   | Prints the string from after the space to the ending quote, i.e. ." hello" prints "hello"               |
| .S   | ( -- )                   | Prints the stack size and values on the stack from bottom to top                                        |
| :    | ( -- )                   | Starts the definition of a word                                                                         |
| ;    | ( -- )                   | Ends the definition of a word                                                                           | 

