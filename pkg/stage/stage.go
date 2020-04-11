package stage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Akatsuki-py/PokeTraveler/pkg/char"
	"github.com/Akatsuki-py/PokeTraveler/pkg/object"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

/*
ゲームに必要なもの
各タイルの画像データ
各タイルのインデックス
各タイルのプロパティ一覧
*/

// Stage - マップのデータ
type Stage struct {
	name       string
	Index      int              // マップのIndex
	Width      int              // マップの横幅(タイル)
	Height     int              // マップの立幅(タイル)
	Image      *ebiten.Image    // マップ全体を画像データにしたもの
	TileIndex  []int            // len = Width*Height
	Properties map[int]Property // タイル番号 => プロパティ
	Actions    []*Action
	Objects    []*object.Object
	Warps      []*Warp
	BGM        *BGM
	Meta       *Meta
}

const (
	assetPath = "asset/map"
)

var (
	popupImage, _, _ = ebitenutil.NewImageFromFile("asset/map/popup.png", ebiten.FilterDefault)
)

// Load - マップを読み込む関数
func (stage *Stage) Load(stagename string, index int) {
	filename := fmt.Sprintf("asset/map/%s/map%d/stage.json", stagename, index)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	raw := new(rawStage)
	if err := json.Unmarshal(file, raw); err != nil {
		panic(err)
	}

	stage.Properties = map[int]Property{}
	stage.Width = raw.Width
	stage.Height = raw.Height

	stage.Image, _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("%s/%s/map%d/stage.png", assetPath, stagename, index), ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	stage.TileIndex = make([]int, stage.Height*stage.Width)
	copy(stage.TileIndex, raw.Layers[0].Data)
	for i, layer := range raw.Layers {
		if i == 0 {
			continue
		}
		for i, tile := range layer.Data {
			if tile != 0 {
				stage.TileIndex[i] = tile
			}
		}
	}

	// 各タイルセットについて
	for _, tileset := range raw.Tilesets {
		firstGID := tileset.FirstGID
		source := tileset.Source
		filename := "asset" + source[8:]
		stage.loadProps(firstGID, filename)
	}

	stage.name = stagename
	stage.Index = index
	stage.loadActions(fmt.Sprintf("%s/%s/map%d/actions.json", assetPath, stagename, index))
	stage.loadObjects(fmt.Sprintf("%s/%s/map%d/objects.json", assetPath, stagename, index))
	stage.loadWarps(fmt.Sprintf("%s/%s/map%d/warp.json", assetPath, stagename, index))
	stage.loadBGM(fmt.Sprintf("%s/%s/map%d/bgm.json", assetPath, stagename, index))

	// メタデータ
	stage.Meta = newMeta(stagename)
}

// Name - Get stage name
func (stage *Stage) Name() string {
	return stage.name
}

// Popup - マップのPopup画像を取得する
func (stage *Stage) Popup() (popup *ebiten.Image, ok bool) {
	ok = stage.Meta.Popup
	if ok {
		// ポップアップが有効な時
		popup, _ = ebiten.NewImageFromImage(popupImage, ebiten.FilterDefault)
		name := stage.Name()
		x := 72 - len(stage.Name())*8/2
		char.RenderString(popup, name, float64(x), 16)
	}

	return popup, ok
}
