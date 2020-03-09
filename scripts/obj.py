import sys

def calc_obj_id(color, obj_id):
    if color == "blue":
        return int(obj_id) + 18
    elif color == "brown":
        return int(obj_id) + 18 + (144*1)
    elif color == "gray":
        return int(obj_id) + 18 + (144*2)
    elif color == "green":
        return int(obj_id) + 18 + (144*3)
    elif color == "pink":
        return int(obj_id) + 18 + (144*4)
    elif color == "red":
        return int(obj_id) + 18 + (144*5)
    elif color == "special":
        return int(obj_id) + 18 + (144*6)
    elif color == "user":
        return int(obj_id) + 18 + (144*6) + 8
    else:
        print("color: blue brown green pink red special user")
        return -1

def main():
    args = sys.argv

    if len(args) < 3:
        print("please input color and id")
        return

    color = args[1]
    obj_id = args[2]
    print(calc_obj_id(color, obj_id))

if __name__ == "__main__":
    main()
