package menu

import (
	"demo/pkg/char"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	cursorImage, _, _ = ebitenutil.NewImageFromFile("asset/menu/cursor.png", ebiten.FilterDefault)
	menuImage, _, _   = ebitenutil.NewImageFromFile("asset/menu/menu.png", ebiten.FilterDefault)
)

type Menu struct {
	cursor int
	list   []string
	image  *ebiten.Image
}

func New() *Menu {
	m := &Menu{
		list: []string{"Map", "Save", "Exit"},
	}

	for i, item := range m.list {
		x, y := float64(16), float64(16*i+16)
		char.RenderString(menuImage, item, x, y)
	}

	m.image = m.setCursor(0)
	return m
}

// Increment cursor
func (m *Menu) Increment() {
	m.cursor++
	if m.cursor == len(m.list) {
		m.cursor = 0
	}
	m.image = m.setCursor(m.cursor)
}

// Decrement cursor
func (m *Menu) Decrement() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.list) - 1
	}
	m.image = m.setCursor(m.cursor)
}

// Current - get current menu
func (m *Menu) Current() string {
	return m.list[m.cursor]
}

// Image - get current image
func (m *Menu) Image() *ebiten.Image {
	return m.image
}

func (m *Menu) setCursor(cursor int) *ebiten.Image {
	result, _ := ebiten.NewImageFromImage(menuImage, ebiten.FilterDefault)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(8), float64(16*cursor+16))
	result.DrawImage(cursorImage, op)
	return result
}

// Exit - メニューを閉じるときの処理
func (m *Menu) Exit() {
	m.cursor = 0
	m.image = m.setCursor(0)
}
