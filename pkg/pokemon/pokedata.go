package pokemon

import (
	"fmt"

	"github.com/Akatsuki-py/PokeTraveler/pkg/sound"
	"github.com/Akatsuki-py/PokeTraveler/pkg/util"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	iconPath  = "asset/pokemon/icon/%d/%d.png"
	imagePath = "asset/pokemon/image/%s.png"
)

// PokeData - ゲームのステートに左右されないデータ
type PokeData struct {
	ID    int              // ポケモンの図鑑番号
	icon  [2]*ebiten.Image // miniDexのアイコン
	image *ebiten.Image    // 対戦とかで出てくるグラフィック
}

var PokeDex map[int]*PokeData = map[int]*PokeData{}

func initPokeDex() {
	for i := 1; i < 251+1; i++ {
		PokeDex[i] = newPokeData(i)
	}
}

func newPokeData(ID int) *PokeData {
	pd := &PokeData{
		ID: ID,
	}

	for i := 0; i < 2; i++ {
		iconPath := fmt.Sprintf(iconPath, ID, i)
		icon, _, _ := ebitenutil.NewImageFromFile(iconPath, ebiten.FilterDefault)
		pd.icon[i] = icon
	}

	imagePath := fmt.Sprintf(imagePath, util.PaddingID(ID))
	img, _, _ := ebitenutil.NewImageFromFile(imagePath, ebiten.FilterDefault)
	pd.image = img

	return pd
}

// Icon - miniDexのアイコンを取得する フレームによって変わる
func (p *PokeData) Icon(frame int) *ebiten.Image {
	switch {
	case frame%32 < 16:
		return p.icon[0]
	default:
		return p.icon[1]
	}
}

// Image - ポケモンのイメージを取得する
func (p *PokeData) Image() *ebiten.Image {
	return p.image
}

// Cry - 鳴き声を出す
func (p *PokeData) Cry() {
	sound.Cry(p.ID)
}
