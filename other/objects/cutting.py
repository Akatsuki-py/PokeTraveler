import cv2
import os

namelist = ["blue", "brown",  "gray", "green", "pink", "red"]

for name in namelist:
    os.mkdir(name)

img = cv2.imread("blue.png", cv2.IMREAD_UNCHANGED)

count = 0
for h in range(23):
    for w in range(25):
        tile = img[16*h:16*(h+1), 16*w:16*(w+1)]
        cv2.imwrite("./blue/{}.png".format(count), tile)
        count += 1

img = cv2.imread("brown.png", cv2.IMREAD_UNCHANGED)

count = 0
for h in range(23):
    for w in range(25):
        tile = img[16*h:16*(h+1), 16*w:16*(w+1)]
        cv2.imwrite("./brown/{}.png".format(count), tile)
        count += 1

img = cv2.imread("gray.png", cv2.IMREAD_UNCHANGED)

count = 0
for h in range(23):
    for w in range(25):
        tile = img[16*h:16*(h+1), 16*w:16*(w+1)]
        cv2.imwrite("./gray/{}.png".format(count), tile)
        count += 1

img = cv2.imread("green.png", cv2.IMREAD_UNCHANGED)

count = 0
for h in range(23):
    for w in range(25):
        tile = img[16*h:16*(h+1), 16*w:16*(w+1)]
        cv2.imwrite("./green/{}.png".format(count), tile)
        count += 1

img = cv2.imread("pink.png", cv2.IMREAD_UNCHANGED)

count = 0
for h in range(23):
    for w in range(25):
        tile = img[16*h:16*(h+1), 16*w:16*(w+1)]
        cv2.imwrite("./pink/{}.png".format(count), tile)
        count += 1

img = cv2.imread("red.png", cv2.IMREAD_UNCHANGED)

count = 0
for h in range(23):
    for w in range(25):
        tile = img[16*h:16*(h+1), 16*w:16*(w+1)]
        cv2.imwrite("./red/{}.png".format(count), tile)
        count += 1
