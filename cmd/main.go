package main

import (
	"demo/pkg/char"
	"demo/pkg/ethan"
	"demo/pkg/object"
	"demo/pkg/sound"
	"demo/pkg/stage"
	"demo/pkg/townmap"
	"demo/pkg/window"
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten"
)

const (
	modeStage  = iota
	modeOneWay // 段差
	modeWindow
	modeWarp
	modeTownMap
)

// Game ゲーム情報を管理する
type Game struct {
	Count    int
	Stage    stage.Stage
	Ethan    ethan.Ethan
	Mode     int
	coolTime uint
	TownMap  townmap.TownMap
}

var game Game
var win *window.Window
var lastAction int

func initGame(game *Game) {
	game.Count = 0
	game.Ethan = *ethan.New(1, 37*16, 16*16)
	game.Mode = modeStage
	game.TownMap = *townmap.New()

	char.Init()
	sound.InitSE()
}

func render(screen *ebiten.Image) error {
	defer func() {
		game.Count++
		if game.Count%2 == 0 && win != nil {
			win.IncrementCharPointer()
		}
		if game.coolTime > 0 {
			game.coolTime--
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
		// 主人公がマス目の間を移動中
		if game.Ethan.Moving() {
			property := game.Stage.GetProp(game.Ethan.X, game.Ethan.Y)
			object := game.Stage.GetObject(game.Ethan.X, game.Ethan.Y)
			if property.Block == 0 && object == nil {
				game.Ethan.GoAhead()
			} else {
				game.Ethan.Collision()
			}

			if warp := game.Stage.GetWarp(game.Ethan.X, game.Ethan.Y); warp != nil && (warp.InOut == "in" || warp.InOut == "none") {
				doWarp(warp)
			}
		} else {
			// 主人公がマス目にいるときはアクションを受け付ける
			game.Ethan.Move()
			goAhead := false
			switch {
			case ebiten.IsKeyPressed(ebiten.KeyUp) && isActionOK():
				if game.Ethan.IsOriented(object.Up) {
					goAhead = true
					game.coolTime = 17
				} else {
					game.Ethan.SetDirection(object.Up)
					game.coolTime = 5
				}
			case ebiten.IsKeyPressed(ebiten.KeyDown) && isActionOK():
				if game.Ethan.IsOriented(object.Down) {
					goAhead = true
					game.coolTime = 17
				} else {
					game.Ethan.SetDirection(object.Down)
					game.coolTime = 5
				}
			case ebiten.IsKeyPressed(ebiten.KeyRight) && isActionOK():
				if game.Ethan.IsOriented(object.Right) {
					goAhead = true
					game.coolTime = 17
				} else {
					game.Ethan.SetDirection(object.Right)
					game.coolTime = 5
				}
			case ebiten.IsKeyPressed(ebiten.KeyLeft) && isActionOK():
				if game.Ethan.IsOriented(object.Left) {
					goAhead = true
					game.coolTime = 17
				} else {
					game.Ethan.SetDirection(object.Left)
					game.coolTime = 5
				}
			case btnA() && isActionOK():
				propety := game.Stage.GetProp(game.Ethan.Ahead())
				object := game.Stage.GetObject(game.Ethan.Ahead())
				if propety.Action == 1 {
					action := game.Stage.GetAction(game.Ethan.Ahead())
					// アクションがあるならそのアクションを取らせる
					if action != nil && action.Type == "text" {
						game.Mode = modeWindow
						win = window.New(action.Value)
						win.Render(screen)
					}
				} else if object != nil {
					// 人との会話
					game.Mode = modeWindow
					object.SetDirectionByPoint(game.Ethan.X, game.Ethan.Y)
					win = window.New(object.Text)
					win.Render(screen)
				}
				game.coolTime = 17
			case btnStart() && isActionOK():
				game.Mode = modeTownMap
				game.coolTime = 17
			}

			// 障害物や段差を考慮して前に進ませる
			if goAhead {
				property := game.Stage.GetProp(game.Ethan.Ahead()) // 障害物、段差
				object := game.Stage.GetObject(game.Ethan.Ahead()) // 人
				action := game.Stage.GetAction(game.Ethan.Ahead()) // アクションオブジェクト

				if property.Action == 1 && action != nil && action.Type == "text" {
					// 移動先にテキストブロックがある
					game.Ethan.Collision()
				} else if oneway := property.OneWay; oneway > 0 {
					// 移動先に段差がある
					if game.Ethan.IsOriented(oneway) {
						game.coolTime = 32
						sound.Ledge()
						game.Mode = modeOneWay
					} else {
						game.Ethan.Collision()
					}
				} else if warp := game.Stage.GetWarp(game.Ethan.Ahead()); warp != nil && warp.InOut == "out" {
					// 移動先にワープブロックがある
					doWarp(warp)
				} else if property.Block == 0 && object == nil {
					// 移動先に何もない
					game.Ethan.GoAhead()
				} else if object == nil {
					// 移動先にblock属性を持ったタイルがある object==nilで正しい
					game.Ethan.Collision()
				}
			}
		}
		renderEthan(screen)
	case modeWindow:
		win.Render(screen)
		if ebiten.IsKeyPressed(ebiten.KeyS) && isActionOK() {
			if win.IsEnd() {
				game.Mode = modeStage
			} else {
				win.NextPage()
				win.Render(screen)
			}
			game.coolTime = 17
		}
		renderEthan(screen)
	case modeWarp:
		screen.Fill(color.NRGBA{0xff, 0xff, 0xdd, 0xff})
		if game.coolTime == 0 {
			game.Mode = modeStage
		}
	case modeOneWay:
		if game.coolTime > 0 {
			game.Ethan.GoAhead()
		} else {
			game.Mode = modeStage
		}
		renderEthan(screen)
	case modeTownMap:
		if btnStart() && isActionOK() {
			game.Mode = modeStage
			game.coolTime = 17
		}
		renderTownMap(screen)
	}

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
		op.GeoM.Translate(float64(64-game.Ethan.X+obj.X), float64(64-game.Ethan.Y+obj.Y-4))
		screen.DrawImage(obj.Avatar(), op)
	}
}

func moveObject() {
	for _, obj := range game.Stage.Objects {
		if obj.Moving() {
			obj.GoAhead()
		}

		if game.Count%100 == 0 {
			direction := object.RandamDirection()
			aheadX, aheadY := obj.Ahead(direction)
			property := game.Stage.GetProp(aheadX, aheadY)
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
	switch game.Mode {
	case modeOneWay:
		// 段差を飛ぶときは主人公のレンダリングに特殊処理
		switch {
		case game.coolTime > 16:
			op.GeoM.Translate(float64(64), float64(64-(4+(32-game.coolTime))))
			screen.DrawImage(game.Ethan.Avatar(), op)
		default:
			op.GeoM.Translate(float64(64), float64(64-(4+game.coolTime)))
			screen.DrawImage(game.Ethan.Avatar(), op)
		}

		hopOp := &ebiten.DrawImageOptions{}
		hopOp.GeoM.Translate(float64(64), float64(64+8-4))
		screen.DrawImage(game.Ethan.HopImage, hopOp)
	default:
		op.GeoM.Translate(float64(64), float64(64-4))
		screen.DrawImage(game.Ethan.Avatar(), op)
	}
}

func doWarp(warp *stage.Warp) {
	if warp.InOut == "in" {
		sound.GoInside()
		game.Mode = modeWarp
		game.coolTime = 20
	} else if warp.InOut == "out" {
		sound.GoOutside()
		game.Mode = modeWarp
		game.coolTime = 20
	}
	game.Stage.Load(warp.Dst, warp.DstID)
	game.Ethan.Set(warp.Pos[0]*16, warp.Pos[1]*16)
}

func main() {
	os.Exit(run())
}

func run() int {
	initGame(&game)
	game.Stage.Load("Oxalis City", 1)

	if err := ebiten.Run(render, 160, 144, 2, "demo"); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return 1
	}
	return 0
}

func btnA() bool {
	return ebiten.IsKeyPressed(ebiten.KeyS)
}

func btnStart() bool {
	return ebiten.IsKeyPressed(ebiten.KeyEnter)
}

func isActionOK() bool {
	return game.coolTime == 0
}

func renderTownMap(screen *ebiten.Image) {
	region := "naljo"
	tm := game.TownMap.Regions[region].Image
	screen.DrawImage(tm, nil)
}
