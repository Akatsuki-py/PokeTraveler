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

func playBGM(fade bool) {
	p, _ := audio.NewPlayer(audioContext, bgm.stream)

	fadeinCount := 0
	if fade {
		fadeinCount = 25
	} else {
		p.Play()
	}
	fadeoutCount := 0

loop:
	for range time.Tick(time.Millisecond * 100) {
		select {
		case <-done:
			done = make(chan interface{})
			if fade {
				fadeoutCount = 20
			} else {
				p.Close()
				break loop
			}
		default:
			if fadeoutCount > 0 {
				p.SetVolume(float64(fadeoutCount) * 0.05)
				fadeoutCount--
				if fadeoutCount == 0 {
					p.Close()
					break loop
				}
			}
			if fadeinCount > 0 {
				fadeinCount--
				if fadeinCount == 0 {
					p.Play()
				}
			}
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
