package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/Akatsuki-py/PokeTraveler/pkg/char"
	"github.com/Akatsuki-py/PokeTraveler/pkg/ethan"
	"github.com/Akatsuki-py/PokeTraveler/pkg/menu"
	"github.com/Akatsuki-py/PokeTraveler/pkg/object"
	"github.com/Akatsuki-py/PokeTraveler/pkg/sound"
	"github.com/Akatsuki-py/PokeTraveler/pkg/stage"
	"github.com/Akatsuki-py/PokeTraveler/pkg/townmap"
	"github.com/Akatsuki-py/PokeTraveler/pkg/util"
	"github.com/Akatsuki-py/PokeTraveler/pkg/window"
	"github.com/hajimehoshi/ebiten"
)

const (
	modeStage  = iota
	modeOneWay // 段差
	modeWindow
	modeWarp
	modeMenu
	modeTownMap
)

// Game ゲーム情報を管理する
type Game struct {
	Count    int
	Stage    stage.Stage
	Ethan    ethan.Ethan
	Mode     int
	coolTime uint
	Menu     menu.Menu
	TownMap  townmap.TownMap
}

var game Game
var win *window.Window
var lastAction int

func initGame(game *Game) {
	char.Init()
	game.Ethan = *ethan.New(2, 37*16, 16*16)
	game.Mode = modeStage
	game.TownMap = *townmap.New()
	game.Menu = *menu.New()
	sound.InitSE()

	game.Stage.Load("Oxalis City", 1)
}

func render(screen *ebiten.Image) error {
	if game.Count == 0 {
		initGame(&game)
	}

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
				changeStage(screen, warp)
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

			case util.BtnA() && isActionOK():
				object := game.Stage.GetObject(game.Ethan.Ahead())
				if action := game.Stage.GetAction(game.Ethan.Ahead()); action != nil {
					// アクションがあるならそのアクションを取らせる
					if action.Type == "text" {
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

			case util.BtnStart() && isActionOK():
				// メニューを開く
				game.Mode = modeMenu
				sound.Menu()
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
					changeStage(screen, warp)
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

	case modeMenu:
		switch {
		case util.BtnA() && isActionOK():
			sound.Select()
			game.coolTime = 17

			switch game.Menu.Current() {
			case "Map":
				game.TownMap.Cursor.Valid = false
				game.Mode = modeTownMap
			case "Exit":
				game.Menu.Exit()
				game.Mode = modeStage
			}

		case (util.BtnStart() || util.BtnB()) && isActionOK():
			game.Menu.Exit()
			game.Mode = modeStage
			game.coolTime = 17

		case util.KeyUp() && isActionOK():
			sound.Select()
			game.Menu.Decrement()
			game.coolTime = 17

		case util.KeyDown() && isActionOK():
			sound.Select()
			game.Menu.Increment()
			game.coolTime = 17
		}

		renderEthan(screen)
		renderMenu(screen)

	case modeTownMap:
		if (util.BtnStart() || util.BtnB()) && isActionOK() {
			game.Mode = modeMenu
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

func changeStage(screen *ebiten.Image, warp *stage.Warp) {
	previous := game.Stage.Name()

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

	current := warp.Dst

	// ステージが変わった際はポップアップを出す
	if previous != current {
		if popup, ok := game.Stage.Popup(); ok {
			start := game.Count
			go func() {
				for game.Count-start < 120 {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(0), float64(112))
					screen.DrawImage(popup, op)
				}
			}()
		}
	}
}

func main() {
	os.Exit(run())
}

func run() int {
	if err := ebiten.Run(render, 160, 144, 2, "PokeTraveler"); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return 1
	}
	return 0
}

func isActionOK() bool {
	return game.coolTime == 0
}

func renderMenu(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(84), float64(0))
	screen.DrawImage(game.Menu.Image(), op)
}

func renderTownMap(screen *ebiten.Image) {
	stage := game.Stage.Name()
	avatar := game.Ethan.AvatarDown()
	tm := game.TownMap.Open(stage, avatar)
	screen.DrawImage(tm, nil)
}
