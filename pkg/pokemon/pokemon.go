package pokemon

import "github.com/hajimehoshi/ebiten"

// PokeData - ゲームのステートに左右されないデータ
type PokeData struct {
	ID    int              // ポケモンの図鑑番号
	Icon  [2]*ebiten.Image // miniDexのアイコン
	Image *ebiten.Image    // 対戦とかで出てくるグラフィック
	Text  []string         // 話しかけた時のテキスト
}

// Pokemon - ゲーム内のポケモンのデータ
type Pokemon struct {
	Data *PokeData
	X    int
	Y    int
}
