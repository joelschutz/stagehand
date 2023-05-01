package stagehand

import (
	"image"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

// MaxInt returns the maximum of two integers
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Pre-draw scenes
func PreDraw[T any](bounds image.Rectangle, fromScene, toScene Scene[T]) (*ebiten.Image, *ebiten.Image) {
	fromImg := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	fromScene.Draw(fromImg)

	toImg := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	toScene.Draw(toImg)

	return toImg, fromImg
}
