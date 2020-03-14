package window

import (
	"demo/pkg/char"
	"demo/pkg/sound"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	messageWindow, _, _ = ebitenutil.NewImageFromFile("asset/window/message.png", ebiten.FilterDefault)
)

// Window - ウィンドウをつかさどる構造体
type Window struct {
	Text    []string
	page    uint
	charPtr uint
}

// New - コンストラクタ
func New(text []string) *Window {
	return &Window{
		Text:    text,
		page:    0,
		charPtr: 0,
	}
}

func (win *Window) Render(screen *ebiten.Image) {
	tw, _ := ebiten.NewImageFromImage(messageWindow, ebiten.FilterDefault)
	text := win.Text[win.page]
	win.renderText(tw, text)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(0), float64(96))
	screen.DrawImage(tw, op)
}

func (win *Window) IsEnd() bool {
	return win.Text[win.page+1] == eventEND
}

func (win *Window) NextPage() {
	if win.ThisPageEnd() {
		win.page++
		win.charPtr = 0
		sound.Select()
	}
}

func (win *Window) renderText(image *ebiten.Image, str string) {
	var col, row uint
	for i, r := range str {

		if win.charPtr < uint(i)+1 {
			break
		}

		c := string(r)
		switch c {
		case "'":
			continue
		case "\r":
			continue
		case "\n":
			col = 0
			row++
		case "m", "r", "s", "t", "v":
			if i > 0 && string(str[i-1]) == "'" {
				c = "'" + c
			}
			char.RenderChar(image, c, float64((col+1)*8+1), float64((row+1)*16))
			col++
		default:
			char.RenderChar(image, c, float64((col+1)*8+1), float64((row+1)*16))
			col++
		}
	}
}

func (win *Window) IncrementCharPointer() {
	length := uint(len(win.Text[win.page]))
	if win.charPtr >= length {
		return
	}
	win.charPtr++
}

func (win *Window) ThisPageEnd() bool {
	if win.page >= uint(len(win.Text)) {
		return false
	}
	return win.charPtr+1 >= uint(len(win.Text[win.page]))
}
