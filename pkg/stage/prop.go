package stage

import (
	"encoding/json"
	"io/ioutil"
)

// Property - タイルのプロパティ
type Property struct {
	Block  int // 通行可能か
	Action int // このタイルに対して何らかのアクションが可能か？
	OneWay int // 通行可能な方向 0全方向可能 1下のみ可能 2右のみ 3左のみ
}

// GetProp - Get tile property
func (stage *Stage) GetProp(x, y int) (target *Property) {
	target = &Property{Block: 1}

	if x >= 0 && x/16 < stage.Width && y >= 0 && y/16 < stage.Height {
		index := (y/16)*stage.Width + (x / 16)
		tileIndex := stage.TileIndex[index]
		property, ok := stage.Properties[tileIndex]
		if ok {
			target = &property
		} else {
			target = &Property{}
		}
		return target
	}

	if warp := stage.GetWarp(x, y); warp != nil {
		return &Property{}
	}

	return target
}

func (stage *Stage) loadProps(firstGID int, filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	tileset := new(TileSet)
	if err := json.Unmarshal(file, tileset); err != nil {
		panic(err)
	}

	// 各タイルのプロパティをセットしていく
	for _, tile := range tileset.List {
		tileID := tile.ID + firstGID

		newProperty := Property{}
		for _, property := range tile.Properties {
			switch property.Name {
			case "block":
				newProperty.Block = property.Value
			case "action":
				newProperty.Action = property.Value
			case "oneway":
				newProperty.OneWay = property.Value
			}
		}
		stage.Properties[tileID] = newProperty
	}
}
