package sound

import (
	"fmt"

	"github.com/Akatsuki-py/PokeTraveler/pkg/util"
	"github.com/hajimehoshi/ebiten/audio"
)

const (
	sampleRate = 44100
	assetPath  = "asset/sound"
)

var (
	audioContext, _ = audio.NewContext(sampleRate)
)

func InitSE() {
	setSE(assetPath+"/se/Select.wav", selectSE)
	setSE(assetPath+"/se/Collision.wav", collisionSE)
	setSE(assetPath+"/se/GoInside.wav", GoInsideSE)
	setSE(assetPath+"/se/GoOutside.wav", GoOutsideSE)
	setSE(assetPath+"/se/Ledge.wav", LedgeSE)
	setSE(assetPath+"/se/Menu.wav", MenuSE)
	setSE(assetPath+"/se/Save.wav", SaveSE)

	for i := 1; i <= 251; i++ {
		cry := fmt.Sprintf(assetPath+"/cry/%sCry.wav", util.PaddingID(i))
		CrySE[i] = NewSE(cry)
	}
}

func InitBGM(bgmname string, fade bool) {
	openBGM(fmt.Sprintf(assetPath+"/bgm/%s", bgmname))
	go playBGM(fade)
}

func ExitBGM() {
	closeBGM()
}

func Exit() {
	closeBGM()
	closeSE(selectSE)
}
