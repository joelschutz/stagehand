package stagehand

import (
	"image"
	"time"

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

// Converts a frequency(cycle/second) to a factor(change/cycle) for a given duration
func DurationToFactor(frequency float64, duration time.Duration) float64 {
	return (1 / frequency) / duration.Seconds()
}

// Calculates the fraction of the duration that has passed since the initial time
func CalculateProgress(initialTime time.Time, duration time.Duration) float64 {
	return float64(time.Since(initialTime)) / float64(duration)
}
