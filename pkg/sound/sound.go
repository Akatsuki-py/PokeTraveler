package sound

import (
	"fmt"

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
}

func InitBGM(bgmname string) {
	closeBGM()
	openBGM(fmt.Sprintf(assetPath+"/bgm/%s", bgmname))
	go playBGM()
}

func Exit() {
	closeBGM()
	closeSE(selectSE)
}
