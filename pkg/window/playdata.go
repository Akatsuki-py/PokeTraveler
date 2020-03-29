package window

import (
	"github.com/Akatsuki-py/PokeTraveler/pkg/char"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	playDataImage, _, _ = ebitenutil.NewImageFromFile("asset/window/playdata.png", ebiten.FilterDefault)
)

// PlayData - ゲームのプレイ状態を表すウィンドウ
type PlayData struct {
	image *ebiten.Image
}

// NewPlayData - コンストラクタ
func NewPlayData() *PlayData {
	pd := &PlayData{}
	return pd
}

// Image - get window image
func (pd *PlayData) Image() *ebiten.Image {
	return pd.image
}

// SetImage - set window image
func (pd *PlayData) SetImage(name string) {
	newImage, _ := ebiten.NewImageFromImage(playDataImage, ebiten.FilterDefault)
	time := "00:00" // TODO
	char.RenderString(newImage, name, 64, 16)
	char.RenderString(newImage, time, 80, 64)
	pd.image = newImage
}
