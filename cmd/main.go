package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/Akatsuki-py/PokeTraveler/pkg/char"
	"github.com/Akatsuki-py/PokeTraveler/pkg/ethan"
	"github.com/Akatsuki-py/PokeTraveler/pkg/menu"
	"github.com/Akatsuki-py/PokeTraveler/pkg/object"
	"github.com/Akatsuki-py/PokeTraveler/pkg/save"
	"github.com/Akatsuki-py/PokeTraveler/pkg/sound"
	"github.com/Akatsuki-py/PokeTraveler/pkg/stage"
	"github.com/Akatsuki-py/PokeTraveler/pkg/townmap"
	"github.com/Akatsuki-py/PokeTraveler/pkg/util"
	"github.com/Akatsuki-py/PokeTraveler/pkg/window"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	modeStage  = iota
	modeOneWay // 段差
	modeWindow
	modeWarp
	modeMenu
	modeTownMap
	modeIntroduction
	modeSave
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
	SaveData *save.Data
	YesNo    *window.YesNoWindow
}

var game Game
var win *window.Window
var lastAction int

var (
	creditImage, _, _ = ebitenutil.NewImageFromFile("asset/intro/credit.png", ebiten.FilterDefault)
	titleImage, _, _  = ebitenutil.NewImageFromFile("asset/intro/title.png", ebiten.FilterDefault)
)

func initGame(game *Game) {
	game.SaveData = save.New("./savedata.json")
	char.Init()
	sound.InitSE()
	game.Mode = modeIntroduction
}

func initStage(game *Game) {

	avatarID := game.SaveData.Avatar.ID
	pointX, pointY := game.SaveData.Point.X, game.SaveData.Point.Y

	game.Ethan = *ethan.New(avatarID, pointX*16, pointY*16)
	game.Mode = modeStage
	game.TownMap = *townmap.New()
	game.Menu = *menu.New()
	game.YesNo = window.NewYesNoWindow()

	stageName, stageIndex := game.SaveData.Point.Stage, game.SaveData.Point.Index
	game.Stage.Load(stageName, stageIndex)
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

	// イントロダクションモード(クレジットやタイトルの描画モード)
	if game.Mode == modeIntroduction {
		renderIntroduction(screen)
		return nil
	}

	// それ以外は基本的にステージを描画
	renderStage(screen)

	// オブジェクトの動作
	if game.Mode == modeStage {
		moveObject()
	}

	// オブジェクトの描画
	renderObject(screen)

	switch game.Mode {
	case modeStage:
		// ステージ画面
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
		// メッセージウィンドウ描画モード
		win.Render(screen)
		if ebiten.IsKeyPressed(ebiten.KeyS) && isActionOK() && win.ThisPageEnd() {
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
		// ステージ間の移動中
		screen.Fill(color.NRGBA{0xff, 0xff, 0xdd, 0xff})
		if game.coolTime == 0 {
			game.Mode = modeStage
		}

	case modeOneWay:
		// 段差移動中
		if game.coolTime > 0 {
			game.Ethan.GoAhead()
		} else {
			game.Mode = modeStage
		}
		renderEthan(screen)

	case modeMenu:
		// メニュー画面を開いている
		switch {
		case util.BtnA() && isActionOK():
			sound.Select()
			game.coolTime = 17

			switch game.Menu.Current() {
			case "Map":
				game.TownMap.Cursor.Valid = false
				game.Mode = modeTownMap
			case "Save":
				game.Mode = modeSave
				win = window.New(save.Message("ethan"))
				win.Render(screen)
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
		// タウンマップを開いている
		if game.TownMap.Cursor.Moving() {
			game.TownMap.Cursor.GoAhead()
		} else {
			if (util.BtnStart() || util.BtnB()) && isActionOK() {
				game.Mode = modeMenu
				game.coolTime = 17
			} else if isActionOK() {
				switch {
				case util.KeyUp():
					game.TownMap.Cursor.GoUp()
				case util.KeyDown():
					game.TownMap.Cursor.GoDown()
				case util.KeyRight():
					game.TownMap.Cursor.GoRight()
				case util.KeyLeft():
					game.TownMap.Cursor.GoLeft()
				}
			}
		}
		renderTownMap(screen)

	case modeSave:
		// セーブ画面
		win.Render(screen)
		renderYesNo(screen)
		if isActionOK() && win.ThisPageEnd() {
			switch {
			case util.BtnA() && game.YesNo.Yes():
				sound.Select()
				game.coolTime = 17

				if win.IsEnd() {
					// 現在の状態をセーブ
					game.SaveData.Point.Stage = game.Stage.Name()
					game.SaveData.Point.Index = game.Stage.Index
					game.SaveData.Point.X, game.SaveData.Point.Y = game.Ethan.X/16, game.Ethan.Y/16
					save.Write(game.SaveData)
					game.Mode = modeStage
				} else {
					win.NextPage()
				}
			case util.BtnA() && !game.YesNo.Yes():
				sound.Select()
				game.Mode = modeMenu
				game.YesNo.SetYes()
				game.coolTime = 17
			case util.BtnB():
				game.Mode = modeMenu
				game.YesNo.SetYes()
				game.coolTime = 17
			case util.KeyUp() && !game.YesNo.Yes():
				sound.Select()
				game.YesNo.SetYes()
				game.coolTime = 17
			case util.KeyDown() && game.YesNo.Yes():
				sound.Select()
				game.YesNo.SetNo()
				game.coolTime = 17
			}
		}
	}

	return nil
}

func renderYesNo(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(0), float64(56))
	screen.DrawImage(game.YesNo.Image(), op)
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

func renderIntroduction(screen *ebiten.Image) {
	switch {
	case game.Count < 150:
		screen.DrawImage(creditImage, nil)
	case game.Count < 210:
		screen.Fill(color.NRGBA{0x00, 0x00, 0x00, 0xff})
	case game.Count == 210:
		screen.Fill(color.NRGBA{0x00, 0x00, 0x00, 0xff})
		sound.ExitBGM()
		sound.InitBGM("7.mp3", false)
	case !isActionOK():
		screen.Fill(color.NRGBA{0xff, 0xff, 0xdd, 0xff})
		if game.coolTime == 1 {
			initStage(&game)
		}
	default:
		screen.DrawImage(titleImage, nil)

		if util.BtnStart() && isActionOK() {
			game.coolTime = 20
		}
	}
}
