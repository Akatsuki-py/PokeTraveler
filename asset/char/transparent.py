import PIL.Image as Image


def do(number):
    name = str(number) + ".png"
    org = Image.open(name)
    trans = Image.new('RGBA', org.size, (0, 0, 0, 0))

    width = org.size[0]
    height = org.size[1]

    for x in range(width):
        for y in range(height):
            pixel = org.getpixel((x, y))

            # 白なら処理しない
            if pixel[0] == 248 and pixel[1] == 248 and pixel[2] == 248:
                continue

            # 白以外なら、用意した画像にピクセルを書き込み
            trans.putpixel((x, y), pixel)
    
    trans.save(name)
    print("convert {}...".format(name))

for i in range(100):
    do(i)

for i in [1000, 1001, 1002]:
    do(i)


