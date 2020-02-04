package object

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	Up = iota
	Down
	Right
	Left
)

// Objects - JSONから取得するための構造体
type Objects struct {
	List []*Object `json:"objects"`
}

// Object - オブジェクトデータ
type Object struct {
	ID        int    `json:"id"`    // object ID
	Start     [2]int `json:"start"` // start point
	X         int
	Y         int
	Territory [][2]int          `json:"territory"` // move range
	Direction int               `json:"direction"` // current direction
	Type      string            `json:"type"`      // object type e.g. npc,trainer,etc...
	Text      []string          `json:"text"`      // what object says
	Image     [10]*ebiten.Image // object avatar data
}

// Load - オブジェクトデータを読み込む
func Load(objfile string) []*Object {
	file, err := ioutil.ReadFile(objfile)
	if err != nil {
		panic(err)
	}

	objects := new(Objects)
	if err := json.Unmarshal(file, objects); err != nil {
		panic(err)
	}

	for i := range objects.List {
		objects.List[i].X = objects.List[i].Start[0] * 16
		objects.List[i].Y = objects.List[i].Start[1] * 16
		objects.List[i].loadImage(objects.List[i].ID)
	}

	return objects.List
}

// Avatar - 現在のイメージデータを返す
func (object *Object) Avatar() *ebiten.Image {
	switch object.Direction {
	case Up:
		switch {
		case object.Y%16 == 0:
			return object.Image[1]
		case object.Y%16 > 8 && (object.Y/16)%2 == 0:
			return object.Image[4]
		case object.Y%16 > 8 && (object.Y/16)%2 == 1:
			return object.Image[8]
		default:
			return object.Image[1]
		}
	case Down:
		switch {
		case object.Y%16 == 0:
			return object.Image[0]
		case object.Y%16 < 8 && (object.Y/16)%2 == 0:
			return object.Image[3]
		case object.Y%16 < 8 && (object.Y/16)%2 == 1:
			return object.Image[7]
		default:
			return object.Image[0]
		}
	case Right:
		switch {
		case object.X%16 == 0:
			return object.Image[6]
		case object.X%16 < 8:
			return object.Image[9]
		default:
			return object.Image[6]
		}
	case Left:
		switch {
		case object.X%16 == 0:
			return object.Image[2]
		case object.X%16 < 8:
			return object.Image[5]
		default:
			return object.Image[2]
		}
	}
	return object.Image[0]
}

// Set - object position. If -1 is set, position is unchanged.
func (object *Object) Set(x, y int) {
	if x >= 0 {
		object.X = x
	}
	if y >= 0 {
		object.Y = y
	}
}

// SetDirection - set object direction
func (object *Object) SetDirection(direction int) {
	switch direction {
	case Up, Down, Right, Left:
		object.Direction = direction
	}
}

// SetDirectionByPoint - set object direction by point
func (object *Object) SetDirectionByPoint(x, y int) {
	switch {
	case y*16 > object.Y*16:
		object.Direction = Down
	case y*16 < object.Y*16:
		object.Direction = Up
	case x*16 > object.X*16:
		object.Direction = Right
	case x*16 < object.X*16:
		object.Direction = Left
	}
}

// Ahead - オブジェクトの一マス前の座標を返す
func (object *Object) Ahead(direction int) (x, y int) {
	switch direction {
	case Up:
		return object.X, object.Y - 16
	case Down:
		return object.X, object.Y + 16
	case Right:
		return object.X + 16, object.Y
	case Left:
		return object.X - 16, object.Y
	default:
		switch object.Direction {
		case Up:
			return object.X, object.Y - 16
		case Down:
			return object.X, object.Y + 16
		case Right:
			return object.X + 16, object.Y
		case Left:
			return object.X - 16, object.Y
		default:
			return -17, -17
		}
	}
}

// GoAhead - 前に進む
func (object *Object) GoAhead() {
	switch object.Direction {
	case Up:
		object.GoUp()
	case Down:
		object.GoDown()
	case Right:
		object.GoRight()
	case Left:
		object.GoLeft()
	}
}

// GoUp - object move up
func (object *Object) GoUp() {
	object.Direction = Up
	object.Y--
}

// GoDown - object move down
func (object *Object) GoDown() {
	object.Direction = Down
	object.Y++
}

// GoRight - object move right
func (object *Object) GoRight() {
	object.Direction = Right
	object.X++
}

// GoLeft - object move left
func (object *Object) GoLeft() {
	object.Direction = Left
	object.X--
}

// Moving - オブジェクトが移動モーション中か
func (object *Object) Moving() bool {
	return object.X%16 != 0 || object.Y%16 != 0
}

// RandamDirection - オブジェクトの向きをランダムに決定して、かつその方向に進行可能か返す
func RandamDirection() (direction int) {
	// 次の向きをランダムに決定
	d := (time.Now().UnixNano() / 1000) % 4
	switch d {
	case 0:
		direction = Down
	case 1:
		direction = Up
	case 2:
		direction = Right
	case 3:
		direction = Left
	}

	return direction
}

// AheadOK - その方向に進行可能かどうか(ブロックは考慮しない)
func (object *Object) AheadOK(direction int) bool {
	x, y := object.Ahead(direction)
	enable := false
	for _, square := range object.Territory {
		if x == square[0]*16 && y == square[1]*16 {
			enable = true
			break
		}
	}
	return enable
}

// loadImage
func (object *Object) loadImage(objectID int) {
	if objectID < 0 {
		return
	}

	var dir string
	switch {
	case objectID < 18:
		// big
		dir = "big"
	case objectID < 18+144:
		// blue
		objectID -= 18
		dir = "blue"
	case objectID < 18+144*2:
		// brown
		objectID -= 18 + 144
		dir = "brown"
	case objectID < 18+144*3:
		// gray
		objectID -= 18 + 144*2
		dir = "gray"
	case objectID < 18+144*4:
		// green
		objectID -= 18 + 144*3
		dir = "green"
	case objectID < 18+144*5:
		// pink
		objectID -= 18 + 144*4
		dir = "pink"
	case objectID < 18+144*6:
		// red
		objectID -= 18 + 144*5
		dir = "red"
	case objectID < 18+144*6+8:
		// special
		objectID -= 18 + 144*6
		dir = "special"
	default:
		// user
		objectID -= 18 + 144*6 + 8
		dir = "user"
	}

	path := fmt.Sprintf("asset/object/%s/%d", dir, objectID)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for i, finfo := range files {
		imgPath := fmt.Sprintf("%s/%s", path, finfo.Name())
		object.Image[i], _, err = ebitenutil.NewImageFromFile(imgPath, ebiten.FilterDefault)
		if err != nil {
			panic(err)
		}
	}
}
