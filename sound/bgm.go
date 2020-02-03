package sound

import (
	"os"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
)

type MP3 struct {
	file   *os.File
	stream *mp3.Stream
}

var (
	bgm    MP3
	isPlay bool
)

func openBGM(path string) {
	bgm.file, _ = os.Open(path)
	bgm.stream, _ = mp3.Decode(audioContext, bgm.file)
	isPlay = true
}

func playBGM() {
	p, _ := audio.NewPlayer(audioContext, bgm.stream)
	p.Play()
}

func closeBGM() {
	if isPlay {
		bgm.file.Close()
		bgm.stream.Close()
	}
}
