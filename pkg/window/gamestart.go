package window

import (
	"github.com/Akatsuki-py/PokeTraveler/pkg/util"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	GameStartContinue = "CONTINUE"
	GameStartNewGame  = "NEW_GAME"
	GameStartOption   = "OPTION"
)

var (
	continueImage, _, _ = ebitenutil.NewImageFromFile("asset/window/continue_game.png", ebiten.FilterDefault)
	newGameImage, _, _  = ebitenutil.NewImageFromFile("asset/window/new_game.png", ebiten.FilterDefault)
)

// GameStartWindow - ゲーム開始時の選択肢を管理するウィンドウ
type GameStartWindow struct {
	cursor       int
	image        *ebiten.Image
	continueFlag bool
}

// NewGameStartWindow - コンストラクタ
func NewGameStartWindow(continueFlag bool) *GameStartWindow {
	gsw := &GameStartWindow{
		continueFlag: continueFlag,
	}

	gsw.image = gsw.setCursor(0)
	return gsw
}

// Mode - 選択されているモードを返す
func (gsw *GameStartWindow) Mode() string {
	if gsw.continueFlag {
		switch gsw.cursor {
		case 0:
			return GameStartContinue
		case 1:
			return GameStartNewGame
		case 2:
			return GameStartOption
		default:
			return GameStartNewGame
		}
	} else {
		switch gsw.cursor {
		case 0:
			return GameStartNewGame
		case 1:
			return GameStartOption
		default:
			return GameStartNewGame
		}
	}
}

// Increment cursor
func (gsw *GameStartWindow) Increment() {
	limit := 2
	if gsw.continueFlag {
		limit = 3
	}

	gsw.cursor++
	if gsw.cursor == limit {
		gsw.cursor = 0
	}

	gsw.setCursor(gsw.cursor)
}

// Decrement cursor
func (gsw *GameStartWindow) Decrement() {
	limit := 2
	if gsw.continueFlag {
		limit = 3
	}

	gsw.cursor--
	if gsw.cursor < 0 {
		gsw.cursor = limit - 1
	}

	gsw.setCursor(gsw.cursor)
}

func (gsw *GameStartWindow) setCursor(cursor int) *ebiten.Image {
	var target *ebiten.Image

	if gsw.continueFlag {
		target, _ = ebiten.NewImageFromImage(continueImage, ebiten.FilterDefault)
	} else {
		target, _ = ebiten.NewImageFromImage(newGameImage, ebiten.FilterDefault)
	}

	util.SetCursor(target, cursor, 16)
	return target
}

// Image - getter of image
func (gsw *GameStartWindow) Image() *ebiten.Image {
	return gsw.image
}
