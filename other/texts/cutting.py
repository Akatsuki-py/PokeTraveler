import cv2
import os

os.mkdir("texts")
img = cv2.imread("texts.png")

count = 0
for h in range(10):
    for w in range(10):
        tile = img[8*h:8*(h+1), 8*w:8*(w+1)]
        cv2.imwrite("./texts/{}.png".format(count), tile)
        count += 1
