def fib(n):
    pf = 0
    f = 1
    i = 0
    while i < n:
        temp = f
        f = f + pf
        pf = temp
        i = i + 1
    return f

i = 0
while i < 100:
    print(fib(i))
    i = i + 1
