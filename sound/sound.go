package sound

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/audio"
)

const (
	sampleRate = 44100
)

var (
	audioContext, _ = audio.NewContext(sampleRate)
)

func InitSE() {
	setSE("sound/asset/se/Select.wav", selectSE)
}

func InitBGM(bgmname string) {
	openBGM(fmt.Sprintf("sound/asset/bgm/%s", bgmname))

	go playBGM()
}

func Exit() {
	closeBGM()
	closeSE(selectSE)
}
