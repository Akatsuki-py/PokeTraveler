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
	l := audio.NewInfiniteLoop(bgm.stream, bgm.stream.Length())
	p, _ := audio.NewPlayer(audioContext, l)

	fadeCount := 40
	fadeVolume := 1. / float64(fadeCount)

	// フェードを用いるBGMはfadeCountフレームで完全に再生される
	fadeinCount := 0
	if fade {
		fadeinCount = fadeCount
		p.Play()
		p.SetVolume(float64(fadeCount-fadeinCount) * fadeVolume)
	} else {
		p.Play()
	}
	fadeoutCount := 0

loop:
	for range time.Tick(time.Millisecond * 10) {
		select {
		case <-done:
			done = make(chan interface{})
			if fade {
				// フェードアウトにはfadeCountフレーム要する
				fadeoutCount = fadeCount
			} else {
				p.Close()
				break loop
			}
		default:
			if fadeoutCount > 0 {
				p.SetVolume(float64(fadeoutCount) * fadeVolume)
				fadeoutCount--
				if fadeoutCount == 0 {
					p.Close()
					break loop
				}
			}
			if fadeinCount > 0 {
				fadeinCount--
				p.SetVolume(float64(fadeCount-fadeinCount) * fadeVolume)
			}
		}
	}
}

func closeBGM() {
	if isPlay {
		close(done)
	}
}
