import os

for i in range(144):
    files = os.listdir("blue/" +str(i))
    count = len(files)
    if count != 6:
        print(i, count)
