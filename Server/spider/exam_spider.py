from numba import jit
import time
@jit
def foo(x,y):
        tt = time.time()
        s = 0
        for i in range(x,y):
                s += i
        print('Time used: {} sec'.format(time.time()-tt))
        return s
print(foo(1,100000000))