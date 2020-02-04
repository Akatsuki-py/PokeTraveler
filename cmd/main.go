package main

import (
	"demo/pkg/char"
	"demo/pkg/ethan"
	"demo/pkg/object"
	"demo/pkg/sound"
	"demo/pkg/stage"
	"demo/pkg/window"
	"fmt"

	"github.com/hajimehoshi/ebiten"
)

const (
	modeStage = iota
	modeWindow
)

// Game ゲーム情報を管理する
type Game struct {
	Count int
	Stage stage.Stage
	Ethan ethan.Ethan
	Mode  int
}

var game Game
var win *window.Window
var lastAction int

func initGame(game *Game) {
	game.Count = 0
	game.Ethan = *ethan.New(64, 64)
	game.Mode = modeStage

	char.Init()
	sound.InitBGM("1.mp3")
	sound.InitSE()
}

func render(screen *ebiten.Image) error {
	defer func() {
		game.Count++
		if game.Count%2 == 0 && win != nil {
			win.IncrementCharPointer()
		}
	}()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	renderStage(screen)

	// オブジェクトの動作
	if game.Mode == modeStage {
		moveObject()
	}

	renderObject(screen)

	switch game.Mode {
	case modeStage:
		if game.Ethan.Moving() {
			property := game.Stage.GetProperty(game.Ethan.X, game.Ethan.Y)
			object := game.Stage.GetObject(game.Ethan.X, game.Ethan.Y)
			if property.Block == 0 && object == nil {
				game.Ethan.GoAhead()
			} else {
				game.Ethan.Collision()
			}

			if warp := game.Stage.GetWarp(game.Ethan.X, game.Ethan.Y); warp != nil {
				doWarp(warp)
			}
		} else {
			game.Ethan.Move()
			goAhead := false
			switch {
			case ebiten.IsKeyPressed(ebiten.KeyUp) && isActionOK():
				game.Ethan.SetDirection(object.Up)
				goAhead = true
			case ebiten.IsKeyPressed(ebiten.KeyDown) && isActionOK():
				game.Ethan.SetDirection(object.Down)
				goAhead = true
			case ebiten.IsKeyPressed(ebiten.KeyRight) && isActionOK():
				game.Ethan.SetDirection(object.Right)
				goAhead = true
			case ebiten.IsKeyPressed(ebiten.KeyLeft) && isActionOK():
				game.Ethan.SetDirection(object.Left)
				goAhead = true
			case btnA() && isActionOK():
				propety := game.Stage.GetProperty(game.Ethan.Ahead())
				object := game.Stage.GetObject(game.Ethan.Ahead())
				if propety.Action == 1 {
					action := game.Stage.GetAction(game.Ethan.Ahead())
					if action != nil {
						fmt.Println(action.Value)
					}
				} else if object != nil {
					game.Mode = modeWindow
					object.SetDirectionByPoint(game.Ethan.X, game.Ethan.Y)
					win = window.New(object.Text)
					win.Render(screen)
				}
			}

			if goAhead {
				property := game.Stage.GetProperty(game.Ethan.Ahead())
				object := game.Stage.GetObject(game.Ethan.Ahead())
				if property.Block == 0 && object == nil {
					game.Ethan.GoAhead()
				} else if object == nil {
					game.Ethan.Collision()
				}
			}
		}
	case modeWindow:
		win.Render(screen)
		if ebiten.IsKeyPressed(ebiten.KeyS) && isActionOK() {
			if win.IsEnd() {
				game.Mode = modeStage
			} else {
				win.NextPage()
				win.Render(screen)
			}
		}
	}

	renderEthan(screen)
	return nil
}

func renderStage(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(64-game.Ethan.X), float64(64-game.Ethan.Y))
	screen.DrawImage(game.Stage.Image, op)
}

func renderObject(screen *ebiten.Image) {
	for _, obj := range game.Stage.Objects {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(64-game.Ethan.X+obj.X), float64(64-game.Ethan.Y+obj.Y))
		screen.DrawImage(obj.Avatar(), op)
	}
}

func moveObject() {
	for _, obj := range game.Stage.Objects {
		if obj.Moving() {
			obj.GoAhead()
		}

		if game.Count%120 == 0 {
			direction := object.RandamDirection()
			aheadX, aheadY := obj.Ahead(direction)
			property := game.Stage.GetProperty(aheadX, aheadY)
			object := game.Stage.GetObject(aheadX, aheadY)
			enable := obj.AheadOK(direction)
			if property.Block == 0 && object == nil && enable {
				obj.SetDirection(direction)
				if !game.Ethan.Exist(aheadX, aheadY) {
					obj.GoAhead()
				}
			}
		}
	}
}

func renderEthan(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(64), float64(64))
	screen.DrawImage(game.Ethan.Avatar(), op)
}

func doWarp(warp *stage.Warp) {
	game.Stage.Load(warp.Dst)
	game.Ethan.Set(warp.Pos[0]*16, warp.Pos[1]*16)
}

func main() {
	initGame(&game)
	game.Stage.Load("Zero Town")

	if err := ebiten.Run(render, 160, 144, 2, "demo"); err != nil {
		panic(err)
	}
}

func btnA() bool {
	return ebiten.IsKeyPressed(ebiten.KeyS)
}

func isActionOK() bool {
	delta := game.Count - lastAction
	coolTime := 17 // 17フレーム
	if delta > coolTime {
		lastAction = game.Count
		return true
	}

	return false
}