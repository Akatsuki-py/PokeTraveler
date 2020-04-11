package stage

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Akatsuki-py/PokeTraveler/pkg/sound"
)

// BGM BGMの詳細
type BGM struct {
	Name string `json:"name"`
	Fade bool   `json:"fade"`
}

func (stage *Stage) loadBGM(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	bgm := new(BGM)
	if err := json.Unmarshal(file, bgm); err != nil {
		panic(err)
	}

	// BGMが変わる時だけを開始する(家などの出入りによってBGMが最初からになるのを避けている)
	if stage.BGM == nil || stage.BGM.Name != bgm.Name {
		stage.BGM = bgm
		sound.ExitBGM()
		go sound.InitBGM(bgm.Name, bgm.Fade)
	}
}
