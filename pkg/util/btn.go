package util

import "github.com/hajimehoshi/ebiten"

func BtnA() bool {
	return ebiten.IsKeyPressed(ebiten.KeyS)
}

func BtnB() bool {
	return ebiten.IsKeyPressed(ebiten.KeyA)
}

func BtnStart() bool {
	return ebiten.IsKeyPressed(ebiten.KeyEnter)
}

func KeyUp() bool {
	return ebiten.IsKeyPressed(ebiten.KeyUp)
}

func KeyDown() bool {
	return ebiten.IsKeyPressed(ebiten.KeyDown)
}
