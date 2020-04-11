package pokemon

import (
	"fmt"

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

	imagePath := fmt.Sprintf(imagePath, toString(ID))
	img, _, _ := ebitenutil.NewImageFromFile(imagePath, ebiten.FilterDefault)
	pd.image = img

	return pd
}

func toString(ID int) string {
	switch {
	case ID < 10:
		return fmt.Sprintf("00%d", ID)
	case ID < 100:
		return fmt.Sprintf("0%d", ID)
	default:
		return fmt.Sprintf("%d", ID)
	}
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
