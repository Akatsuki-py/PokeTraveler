package sound

import (
	"os"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
)

type WAV struct {
	file   *os.File
	stream *wav.Stream
	player *audio.Player
}

var (
	selectSE    = &WAV{}
	collisionSE = &WAV{}
	GoInsideSE  = &WAV{}
	GoOutsideSE = &WAV{}
	LedgeSE     = &WAV{}
	MenuSE      = &WAV{}
)

func setSE(path string, se *WAV) {
	var err error
	se.file, err = os.Open(path)
	if err != nil {
		panic(err)
	}

	se.stream, _ = wav.Decode(audioContext, se.file)
	se.player, _ = audio.NewPlayer(audioContext, se.stream)
}

// NewSE - create WAV instance
func NewSE(path string) *WAV {
	se := &WAV{}
	setSE(path, se)
	return se
}

func PlaySE(se *WAV) {
	if se.player.IsPlaying() {
		se.player.Seek(0)
	} else {
		se.player.Seek(0)
		se.player.Play()
	}
}

func closeSE(se *WAV) {
	se.file.Close()
	se.stream.Close()
}

func Select() {
	PlaySE(selectSE)
}

func Collision() {
	PlaySE(collisionSE)
}

func GoInside() {
	PlaySE(GoInsideSE)
}

func GoOutside() {
	PlaySE(GoOutsideSE)
}

func Ledge() {
	PlaySE(LedgeSE)
}

func Menu() {
	PlaySE(MenuSE)
}
