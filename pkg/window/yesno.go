package window

import (
	"github.com/Akatsuki-py/PokeTraveler/pkg/util"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	yesnoWindowImage, _, _ = ebitenutil.NewImageFromFile("asset/window/yesno.png", ebiten.FilterDefault)
)

const (
	yes = iota
	no
)

// YesNoWindow - はい・いいえの選択肢を管理するウィンドウ
type YesNoWindow struct {
	cursor int
	image  *ebiten.Image
}

// NewYesNoWindow - コンストラクタ
func NewYesNoWindow() *YesNoWindow {
	ynw := &YesNoWindow{}

	ynw.SetYes()
	return ynw
}

// Yes - 選択肢がYesか
func (ynw *YesNoWindow) Yes() bool {
	return ynw.cursor == yes
}

// SetYes - 選択肢をYesにする
func (ynw *YesNoWindow) SetYes() {
	ynw.cursor = yes

	target, _ := ebiten.NewImageFromImage(yesnoWindowImage, ebiten.FilterDefault)
	ynw.image = util.SetCursor(target, yes, 8)
}

// SetNo - 選択肢をNoにする
func (ynw *YesNoWindow) SetNo() {
	ynw.cursor = no

	target, _ := ebiten.NewImageFromImage(yesnoWindowImage, ebiten.FilterDefault)
	ynw.image = util.SetCursor(target, no, 8)
}

// Image - getter of image
func (ynw *YesNoWindow) Image() *ebiten.Image {
	return ynw.image
}
