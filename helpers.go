package stagehand

import (
	"image"
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

// Clock helpers for mocking
type ClockInterface interface {
	Since(time.Time) time.Duration
	Now() time.Time
	Until(time.Time) time.Duration
	Sleep(time.Duration)
}

type RealClock struct{}

func (RealClock) Since(t time.Time) time.Duration { return time.Since(t) }
func (RealClock) Now() time.Time                  { return time.Now() }
func (RealClock) Until(t time.Time) time.Duration { return time.Until(t) }
func (RealClock) Sleep(d time.Duration)           { time.Sleep(d) }

var Clock ClockInterface = RealClock{}

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
	return float64(Clock.Since(initialTime)) / float64(duration)
}
