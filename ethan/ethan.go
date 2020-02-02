package ethan

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	Up = iota
	Down
	Right
	Left
)

// Ethan 主人公のデータ
type Ethan struct {
	Image     [10]*ebiten.Image
	X         int
	Y         int
	direction int
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
	ethan.direction = Down
}

// Avatar Ethan Avatar image
func (ethan *Ethan) Avatar() *ebiten.Image {
	switch ethan.direction {
	case Up:
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
	case Down:
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
	case Right:
		switch {
		case ethan.X%16 == 0:
			return ethan.Image[6]
		case ethan.X%16 < 8:
			return ethan.Image[9]
		default:
			return ethan.Image[6]
		}
	case Left:
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
		ethan.direction = Up
	case "Down", "down":
		ethan.direction = Down
	case "Right", "right":
		ethan.direction = Right
	case "Left", "left":
		ethan.direction = Left
	}
}

// Ahead 主人公の一マス前の座標を返す
func (ethan *Ethan) Ahead() (x, y int) {
	switch ethan.direction {
	case Up:
		return ethan.X, ethan.Y - 16
	case Down:
		return ethan.X, ethan.Y + 16
	case Right:
		return ethan.X + 16, ethan.Y
	case Left:
		return ethan.X - 16, ethan.Y
	}
	return -17, -17
}

// GoAhead 前に進む
func (ethan *Ethan) GoAhead() {
	switch ethan.direction {
	case Up:
		ethan.GoUp()
	case Down:
		ethan.GoDown()
	case Right:
		ethan.GoRight()
	case Left:
		ethan.GoLeft()
	}
}

// GoUp ethan move up
func (ethan *Ethan) GoUp() {
	ethan.direction = Up
	ethan.Y--
}

// GoDown ethan move down
func (ethan *Ethan) GoDown() {
	ethan.direction = Down
	ethan.Y++
}

// GoRight ethan move right
func (ethan *Ethan) GoRight() {
	ethan.direction = Right
	ethan.X++
}

// GoLeft ethan move left
func (ethan *Ethan) GoLeft() {
	ethan.direction = Left
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
		case Up:
			existX = ethan.X/16 == x/16
			existY = (ethan.Y-15)/16 == y/16
		case Down:
			existX = ethan.X/16 == x/16
			existY = (ethan.Y+15)/16 == y/16
		case Right:
			existX = (ethan.X+15)/16 == x/16
			existY = ethan.Y/16 == y/16
		case Left:
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
