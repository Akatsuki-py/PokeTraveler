package stage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Meta struct {
	Region  string  `json:"region"`
	Neutral Neutral `json:"neutral"`
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
