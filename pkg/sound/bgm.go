package sound

import (
	"os"
	"time"

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
	done   = make(chan interface{})
)

func openBGM(path string) {
	bgm.file, _ = os.Open(path)
	bgm.stream, _ = mp3.Decode(audioContext, bgm.file)
	isPlay = true
}

func playBGM() {
	p, _ := audio.NewPlayer(audioContext, bgm.stream)
	p.Play()

loop:
	for range time.Tick(time.Millisecond * 100) {
		select {
		case <-done:
			p.Close()
			done = make(chan interface{})
			break loop
		default:
			if !p.IsPlaying() {
				p.Rewind()
				p.Play()
			}
		}
	}
}

func closeBGM() {
	if isPlay {
		close(done)
	}
}
