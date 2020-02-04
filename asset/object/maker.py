import os
import cv2

for i in range(144):
    for color in ["blue", "brown", "gray", "green", "pink", "red"]:
        files = os.listdir("{}/{}".format(color, i))
        count = len(files)

        if count == 6:
            # 2.png => 6.png
            # 3.png => 7.png
            # 4.png => 8.png
            # 5.png => 9.png
            for j in [2, 3, 4, 5]:
                src = "{}/{}/{}.png".format(color, i, j)
                dst = "{}/{}/{}.png".format(color, i, j+4)
                img = cv2.imread(src, cv2.IMREAD_UNCHANGED)
                inversed_img = cv2.flip(img, 1)
                cv2.imwrite(dst, inversed_img)
        elif count == 3:
            # 2.png => 3.png
            src = "{}/{}/2.png".format(color, i)
            dst = "{}/{}/3.png".format(color, i)
            img = cv2.imread(src, cv2.IMREAD_UNCHANGED)
            inversed_img = cv2.flip(img, 1)
            cv2.imwrite(dst, inversed_img)
