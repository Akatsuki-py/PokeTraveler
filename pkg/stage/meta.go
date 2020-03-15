package stage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Meta struct {
	Region  string  `json:"region"`  // 地方名
	Neutral Neutral `json:"neutral"` // 空を飛ぶなどでこのステージに来た時の開始地点
	Popup   bool    `json:"popup"`   // このステージに来た時にポップアップを行うか
}

type Neutral struct {
	ID int `json:"id"`
	X  int `json:"x"`
	Y  int `json:"y"`
}

func newMeta(stagename string) *Meta {
	filename := fmt.Sprintf("asset/map/%s/meta.json", stagename)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	meta := new(Meta)
	if err := json.Unmarshal(file, meta); err != nil {
		panic(err)
	}

	return meta
}
