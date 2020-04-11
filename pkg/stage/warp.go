package stage

import (
	"encoding/json"
	"io/ioutil"
)

// Warps ワープポイントのリスト
type Warps struct {
	List []*Warp `json:"warps"`
}

// Warp ワープポイント
type Warp struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Dst   string `json:"dst"`
	DstID int    `json:"dstid"`
	Pos   [2]int `json:"pos"`
	InOut string `json:"inout"`
}

// GetWarp - Get warp point
func (stage *Stage) GetWarp(x, y int) (target *Warp) {
	for _, warp := range stage.Warps {
		if warp.X*16 == x && warp.Y*16 == y {
			target = warp
			break
		}
	}
	return target
}

func (stage *Stage) loadWarps(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	warps := new(Warps)
	if err := json.Unmarshal(file, warps); err != nil {
		panic(err)
	}
	stage.Warps = warps.List
}
