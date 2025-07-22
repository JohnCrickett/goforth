: fib over over + ;  ( generate the next number in the Fibonacci sequence )
: fibn 10 1 do fib loop ;
0 1 fibn ( generate the next 10 numbers )
