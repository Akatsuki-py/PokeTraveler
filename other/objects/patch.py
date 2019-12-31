import os

for i in range(104, 137):
    os.rename(str(i), str(i-1))