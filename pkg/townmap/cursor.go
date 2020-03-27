package townmap

const (
	minX = 0
	minY = 0
	maxX = 140
	maxY = 102
)

const (
	Up = iota
	Down
	Right
	Left
)

type Cursor struct {
	x         int // カーソルのX位置(16単位ではなくピクセル単位)
	y         int // カーソルのY位置(16単位ではなくピクセル単位)
	direction int
	Valid     bool // カーソルが有効化されているか
}

// GetXY - x,yを返す
func (cs *Cursor) GetXY() (x, y int) {
	return cs.x, cs.y
}

// GetTileXY - カーソルの移動方向を加味したカーソルのタイル位置を返す
func (cs *Cursor) GetTileXY() (tileX, tileY int) {
	if cs.Moving() {
		switch cs.direction {
		case Up:
			return cs.x / 16, (cs.y - 15) / 16
		case Down:
			return cs.x / 16, (cs.y + 15) / 16
		case Right:
			return (cs.x + 15) / 16, cs.y / 16
		case Left:
			return (cs.x - 15) / 16, cs.y / 16
		}
	}
	return cs.x / 16, cs.y / 16
}

func (cs *Cursor) SetXY(x, y int) {
	cs.x, cs.y = x, y
	cs.Valid = true
}

func (cs *Cursor) Moving() bool {
	return cs.x%16 != 0 || cs.y%16 != 0
}

// GoAhead - カーソルを前に進ませる
func (cs *Cursor) GoAhead() {
	switch cs.direction {
	case Up:
		cs.GoUp()
	case Down:
		cs.GoDown()
	case Right:
		cs.GoRight()
	case Left:
		cs.GoLeft()
	}
}

// GoUp - カーソルを上に動かす
func (cs *Cursor) GoUp() {
	cs.direction = Up
	if cs.y == minY {
		cs.y = maxY
	} else {
		cs.y -= 2
	}
}

// GoDown - カーソルを下に動かす
func (cs *Cursor) GoDown() {
	cs.direction = Down
	if cs.y == maxY {
		cs.y = minY
	} else {
		cs.y += 2
	}
}

// GoRight - カーソルを右に動かす
func (cs *Cursor) GoRight() {
	cs.direction = Right
	if cs.x == maxX {
		cs.x = minX
	} else {
		cs.x += 2
	}
}

// GoLeft - カーソルを左に動かす
func (cs *Cursor) GoLeft() {
	cs.direction = Left
	if cs.x == minX {
		cs.x = maxX
	} else {
		cs.x -= 2
	}
}
