package window

import (
	"fmt"

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
func (pd *PlayData) SetImage(name string, minutes uint) {
	newImage, _ := ebiten.NewImageFromImage(playDataImage, ebiten.FilterDefault)
	time := getPlayTime(minutes)
	char.RenderString(newImage, name, 64, 16)

	start := 96 - 8*(len(time)-3) // プレイ時間によって"xx:yy"の長さが変わるので調節してやる必要がある
	char.RenderString(newImage, time, float64(start), 64)
	pd.image = newImage
}

func getPlayTime(minutes uint) string {
	var hours uint
	hours, minutes = minutes/60, minutes%60

	if minutes < 10 {
		return fmt.Sprintf("%d:0%d", hours, minutes)
	}
	return fmt.Sprintf("%d:%d", hours, minutes)
}
