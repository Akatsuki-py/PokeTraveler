package window

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Window - ウィンドウをつかさどる構造体
type Window struct {
	Text []string
	Page uint
}

// New - コンストラクタ
func New(text []string) *Window {
	return &Window{
		Text: text,
		Page: 0,
	}
}

func (win *Window) RenderText(screen *ebiten.Image) {
	tw, _, _ := ebitenutil.NewImageFromFile("window/text_window.png", ebiten.FilterDefault)
	text := win.Text[win.Page]
	err := ebitenutil.DebugPrint(tw, text)
	if err != nil {
		log.Fatalln(err)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(0), float64(+96))
	screen.DrawImage(tw, op)
}

func (win *Window) IsEnd() bool {
	return len(win.Text) == int(win.Page+1)
}
