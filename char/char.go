package char

import (
	"fmt"

	"demo/util"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	assetPath = "char/asset"
)

var charset map[string]*ebiten.Image = map[string]*ebiten.Image{}

// Init - 文字の初期化処理を行う
func Init() {

	largeAlphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for _, r := range largeAlphabet {
		char := string(r)
		charcode := charcodes[char]
		path := fmt.Sprintf("%s/%d.png", assetPath, charcode)
		image, _, _ := ebitenutil.NewImageFromFile(path, ebiten.FilterDefault)
		charset[char] = image
	}

	smallAlphabet := "abcdefghijklmnopqrstuvwxyz"
	for _, r := range smallAlphabet {
		char := string(r)
		charcode := charcodes[char]
		path := fmt.Sprintf("%s/%d.png", assetPath, charcode)
		image, _, _ := ebitenutil.NewImageFromFile(path, ebiten.FilterDefault)
		charset[char] = image
	}

	otherSymbol := " -!?.♂♀0123456789_@:;[]/"
	for _, r := range otherSymbol {
		char := string(r)
		charcode := charcodes[char]
		path := fmt.Sprintf("%s/%d.png", assetPath, charcode)
		image, _, _ := ebitenutil.NewImageFromFile(path, ebiten.FilterDefault)
		charset[char] = image
	}

	phrases := []string{"'m", "'r", "'s", "'t", "'v", ":L"}
	for _, phrase := range phrases {
		charcode := charcodes[phrase]
		path := fmt.Sprintf("%s/%d.png", assetPath, charcode)
		image, _, _ := ebitenutil.NewImageFromFile(path, ebiten.FilterDefault)
		charset[phrase] = image
	}
}

// RenderChar - 1文字の描画を行う
func RenderChar(target *ebiten.Image, char string, x, y float64) {
	charImage, ok := charset[char]
	if !ok {
		panic(fmt.Sprintf("invalid char [%s]\n", char))
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	target.DrawImage(charImage, op)
}

// RenderString - 文字列の描画を行う
func RenderString(target *ebiten.Image, str string, x, y float64) {
	startX := x

	var skip bool
	var col, row uint // 2文字で1文字とかあるためrange strのインデックスとは独自にカウントする必要あり
	for i, r := range str {

		if skip {
			skip = false
			continue
		}

		char := string(r)
		switch char {
		case "'":
			if i+1 < len(str) && util.Contains([]string{"m", "r", "s", "t", "v"}, string(str[i+1])) {
				char += string(str[i+1])
				skip = true
			}
		case ":":
			if i+1 < len(str) && util.Contains([]string{"L"}, string(str[i+1])) {
				char += string(str[i+1])
				skip = true
			}
		case "\r":
			continue
		case "\n":
			col = 0
			row++
			continue
		}
		x := startX + float64(col*8)
		y += float64(row * 8)
		RenderChar(target, char, x, y)

		col++
	}
}
