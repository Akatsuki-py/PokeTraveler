package main

import (
	"demo/ethan"
	"demo/object"
	"demo/stage"
	"demo/window"
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten"
)

// Game ゲーム情報を管理する
type Game struct {
	Count int
	Stage stage.Stage
	Ethan ethan.Ethan
	Mode  string
}

var game Game
var win *window.Window
var lastAction int64

func initGame(game *Game) {
	game.Count = 0
	game.Ethan.Init(64, 64)
	game.Mode = "stage"
}

func render(screen *ebiten.Image) error {
	defer func() {
		game.Count++
	}()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	renderStage(screen)

	// オブジェクトの動作
	if game.Mode == "stage" {
		moveObject()
	}

	renderObject(screen)

	switch game.Mode {
	case "stage":
		if game.Ethan.Moving() {
			game.Ethan.GoAhead()

			if warp := game.Stage.GetWarp(game.Ethan.X, game.Ethan.Y); warp != nil {
				doWarp(warp)
			}
		} else {
			goAhead := false
			switch {
			case ebiten.IsKeyPressed(ebiten.KeyUp):
				game.Ethan.SetDirection("up")
				goAhead = true
			case ebiten.IsKeyPressed(ebiten.KeyDown):
				game.Ethan.SetDirection("down")
				goAhead = true
			case ebiten.IsKeyPressed(ebiten.KeyRight):
				game.Ethan.SetDirection("right")
				goAhead = true
			case ebiten.IsKeyPressed(ebiten.KeyLeft):
				game.Ethan.SetDirection("left")
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
					game.Mode = "window"
					object.SetDirectionByPoint(game.Ethan.X, game.Ethan.Y)
					win = window.NewWindow(object.Text)
					win.RenderText(screen)
				}
			}

			if goAhead {
				property := game.Stage.GetProperty(game.Ethan.Ahead())
				object := game.Stage.GetObject(game.Ethan.Ahead())
				if property.Block == 0 && object == nil {
					game.Ethan.GoAhead()
				}
			}
		}
	case "window":
		win.RenderText(screen)
		if ebiten.IsKeyPressed(ebiten.KeyS) && isActionOK() {
			if win.IsEnd() {
				game.Mode = "stage"
			} else {
				win.Page++
				win.RenderText(screen)
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
			aheadX, aheadY := (obj).Ahead(direction)
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
	now := time.Now().UnixNano()
	delta := now - lastAction
	coolTime := int64(1000 * 100000 * 3)
	if delta > coolTime {
		lastAction = now
		return true
	}

	return false
}
