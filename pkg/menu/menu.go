package menu

import (
	"github.com/Akatsuki-py/PokeTraveler/pkg/char"
	"github.com/Akatsuki-py/PokeTraveler/pkg/util"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	Map  = "Map"
	Save = "Save"
	Exit = "Exit"
)

var (
	menuImage, _, _ = ebitenutil.NewImageFromFile("asset/menu/menu.png", ebiten.FilterDefault)
)

type Menu struct {
	cursor int
	list   []string
	image  *ebiten.Image
}

func New() *Menu {
	m := &Menu{
		list: []string{Map, Save, Exit},
	}

	for i, item := range m.list {
		x, y := float64(16), float64(16*i+16)
		char.RenderString(menuImage, item, x, y)
	}

	target, _ := ebiten.NewImageFromImage(menuImage, ebiten.FilterDefault)
	m.image = util.SetCursor(target, 0, 16)
	return m
}

// Increment cursor
func (m *Menu) Increment() {
	m.cursor++
	if m.cursor == len(m.list) {
		m.cursor = 0
	}
	target, _ := ebiten.NewImageFromImage(menuImage, ebiten.FilterDefault)
	m.image = util.SetCursor(target, m.cursor, 16)
}

// Decrement cursor
func (m *Menu) Decrement() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.list) - 1
	}
	target, _ := ebiten.NewImageFromImage(menuImage, ebiten.FilterDefault)
	m.image = util.SetCursor(target, m.cursor, 16)
}

// Current - get current menu
func (m *Menu) Current() string {
	return m.list[m.cursor]
}

// Image - get current image
func (m *Menu) Image() *ebiten.Image {
	return m.image
}

// Exit - メニューを閉じるときの処理
func (m *Menu) Exit() {
	m.cursor = 0
	target, _ := ebiten.NewImageFromImage(menuImage, ebiten.FilterDefault)
	m.image = util.SetCursor(target, 0, 16)
}
