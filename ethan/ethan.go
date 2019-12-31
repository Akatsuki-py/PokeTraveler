package ethan

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Ethan 主人公のデータ
type Ethan struct {
	Image     [10]*ebiten.Image
	X         int
	Y         int
	direction string
}

// Init ethan
func (ethan *Ethan) Init(x, y int) {
	ethan.Image[0], _, _ = ebitenutil.NewImageFromFile("ethan/ethan00.png", ebiten.FilterDefault)
	ethan.Image[1], _, _ = ebitenutil.NewImageFromFile("ethan/ethan01.png", ebiten.FilterDefault)
	ethan.Image[2], _, _ = ebitenutil.NewImageFromFile("ethan/ethan02.png", ebiten.FilterDefault)
	ethan.Image[3], _, _ = ebitenutil.NewImageFromFile("ethan/ethan03.png", ebiten.FilterDefault)
	ethan.Image[4], _, _ = ebitenutil.NewImageFromFile("ethan/ethan04.png", ebiten.FilterDefault)
	ethan.Image[5], _, _ = ebitenutil.NewImageFromFile("ethan/ethan05.png", ebiten.FilterDefault)
	ethan.Image[6], _, _ = ebitenutil.NewImageFromFile("ethan/ethan06.png", ebiten.FilterDefault)
	ethan.Image[7], _, _ = ebitenutil.NewImageFromFile("ethan/ethan07.png", ebiten.FilterDefault)
	ethan.Image[8], _, _ = ebitenutil.NewImageFromFile("ethan/ethan08.png", ebiten.FilterDefault)
	ethan.Image[9], _, _ = ebitenutil.NewImageFromFile("ethan/ethan09.png", ebiten.FilterDefault)

	ethan.X = x
	ethan.Y = y
	ethan.direction = "down"
}

// Avatar Ethan Avatar image
func (ethan *Ethan) Avatar() *ebiten.Image {
	switch ethan.direction {
	case "up":
		switch {
		case ethan.Y%16 == 0:
			return ethan.Image[1]
		case ethan.Y%16 > 8 && (ethan.Y/16)%2 == 0:
			return ethan.Image[4]
		case ethan.Y%16 > 8 && (ethan.Y/16)%2 == 1:
			return ethan.Image[8]
		default:
			return ethan.Image[1]
		}
	case "down":
		switch {
		case ethan.Y%16 == 0:
			return ethan.Image[0]
		case ethan.Y%16 < 8 && (ethan.Y/16)%2 == 0:
			return ethan.Image[3]
		case ethan.Y%16 < 8 && (ethan.Y/16)%2 == 1:
			return ethan.Image[7]
		default:
			return ethan.Image[0]
		}
	case "right":
		switch {
		case ethan.X%16 == 0:
			return ethan.Image[6]
		case ethan.X%16 < 8:
			return ethan.Image[9]
		default:
			return ethan.Image[6]
		}
	case "left":
		switch {
		case ethan.X%16 == 0:
			return ethan.Image[2]
		case ethan.X%16 < 8:
			return ethan.Image[5]
		default:
			return ethan.Image[2]
		}
	}
	return ethan.Image[0]
}

// Set Ethan position. If -1 is set, position is unchanged.
func (ethan *Ethan) Set(x, y int) {
	if x >= 0 {
		ethan.X = x
	}
	if y >= 0 {
		ethan.Y = y
	}
}

// SetDirection set ethan direction
func (ethan *Ethan) SetDirection(direction string) {
	switch direction {
	case "Up", "up":
		ethan.direction = "up"
	case "Down", "down":
		ethan.direction = "down"
	case "Right", "right":
		ethan.direction = "right"
	case "Left", "left":
		ethan.direction = "left"
	}
}

// Ahead 主人公の一マス前の座標を返す
func (ethan *Ethan) Ahead() (x, y int) {
	switch ethan.direction {
	case "up":
		return ethan.X, ethan.Y - 16
	case "down":
		return ethan.X, ethan.Y + 16
	case "right":
		return ethan.X + 16, ethan.Y
	case "left":
		return ethan.X - 16, ethan.Y
	}
	return -17, -17
}

// GoAhead 前に進む
func (ethan *Ethan) GoAhead() {
	switch ethan.direction {
	case "up":
		ethan.GoUp()
	case "down":
		ethan.GoDown()
	case "right":
		ethan.GoRight()
	case "left":
		ethan.GoLeft()
	}
}

// GoUp ethan move up
func (ethan *Ethan) GoUp() {
	ethan.direction = "up"
	ethan.Y--
}

// GoDown ethan move down
func (ethan *Ethan) GoDown() {
	ethan.direction = "down"
	ethan.Y++
}

// GoRight ethan move right
func (ethan *Ethan) GoRight() {
	ethan.direction = "right"
	ethan.X++
}

// GoLeft ethan move left
func (ethan *Ethan) GoLeft() {
	ethan.direction = "left"
	ethan.X--
}

// Moving if Ethan is moving?
func (ethan *Ethan) Moving() bool {
	return ethan.X%16 != 0 || ethan.Y%16 != 0
}

// Exist 指定した場所に主人公がいるか
func (ethan *Ethan) Exist(x, y int) bool {
	exist := false

	existX := false
	existY := false
	if ethan.Moving() {
		switch ethan.direction {
		case "up":
			existX = ethan.X/16 == x/16
			existY = (ethan.Y-15)/16 == y/16
		case "down":
			existX = ethan.X/16 == x/16
			existY = (ethan.Y+15)/16 == y/16
		case "right":
			existX = (ethan.X+15)/16 == x/16
			existY = ethan.Y/16 == y/16
		case "left":
			existX = (ethan.X-15)/16 == x/16
			existY = ethan.Y/16 == y/16
		}
	} else {
		existX = ethan.X/16 == x/16
		existY = ethan.Y/16 == y/16
	}

	exist = existX && existY
	return exist
}
