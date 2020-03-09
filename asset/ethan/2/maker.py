import os
import cv2

# 2.png => 6.png
# 3.png => 7.png
# 4.png => 8.png
# 5.png => 9.png
for i in [2, 3, 4, 5]:
    src = "{}.png".format(i)
    dst = "{}.png".format(i+4)
    img = cv2.imread(src, cv2.IMREAD_UNCHANGED)
    inversed_img = cv2.flip(img, 1)
    cv2.imwrite(dst, inversed_img)
