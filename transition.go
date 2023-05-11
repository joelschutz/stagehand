package stagehand

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type SceneTransition[T any] interface {
	Scene[T]
	Start(fromScene, toScene Scene[T])
	End()
}

type BaseTransition[T any] struct {
	fromScene Scene[T]
	toScene   Scene[T]
	sm        *SceneManager[T]
}

func (t *BaseTransition[T]) Start(fromScene, toScene Scene[T]) {
	t.fromScene = fromScene
	t.toScene = toScene
}

// Update updates the transition state
func (t *BaseTransition[T]) Update() error {
	// Update the scenes
	err := t.fromScene.Update()
	if err != nil {
		return err
	}

	err = t.toScene.Update()
	if err != nil {
		return err
	}

	return nil
}

// Layout updates the layout of the scenes
func (t *BaseTransition[T]) Layout(outsideWidth, outsideHeight int) (int, int) {
	sw, sh := t.fromScene.Layout(outsideWidth, outsideHeight)
	tw, th := t.toScene.Layout(outsideWidth, outsideHeight)

	return MaxInt(sw, tw), MaxInt(sh, th)
}

// Loads the next scene
func (t *BaseTransition[T]) Load(state T, manager *SceneManager[T]) {
	t.sm = manager
	t.toScene.Load(state, manager)
}

// Unloads the last scene
func (t *BaseTransition[T]) Unload() T {
	return t.fromScene.Unload()
}

// Ends transition to the next scene
func (t *BaseTransition[T]) End() {
	t.sm.SwitchTo(t.toScene)
}

type FadeTransition[T any] struct {
	BaseTransition[T]
	factor       float32 // factor used for the fade-in/fade-out effect
	alpha        float32 // alpha value used for the fade-in/fade-out effect
	isFadingIn   bool    // whether the transition is currently fading in or out
	frameUpdated bool
}

func NewFadeTransition[T any](factor float32) *FadeTransition[T] {
	return &FadeTransition[T]{
		factor: factor,
	}
}

// Start starts the transition from the given "from" scene to the given "to" scene
func (t *FadeTransition[T]) Start(fromScene, toScene Scene[T]) {
	t.BaseTransition.Start(fromScene, toScene)
	t.alpha = 0
	t.isFadingIn = true
}

// Update updates the transition state
func (t *FadeTransition[T]) Update() error {
	if !t.frameUpdated {
		// Update the alpha value based on the current state of the transition
		if t.isFadingIn {
			t.alpha += t.factor
			if t.alpha >= 1.0 {
				t.alpha = 1.0
				t.isFadingIn = false
			}
		} else {
			t.alpha -= t.factor
			if t.alpha <= 0.0 {
				t.alpha = 0.0
				t.End()
			}
		}
		t.frameUpdated = true
	}

	// Update the scenes
	return t.BaseTransition.Update()
}

// Draw draws the transition effect
func (t *FadeTransition[T]) Draw(screen *ebiten.Image) {
	toImg, fromImg := PreDraw(screen.Bounds(), t.fromScene, t.toScene)
	toOp, fromOp := &ebiten.DrawImageOptions{}, &ebiten.DrawImageOptions{}

	// Draw the scenes with the appropriate alpha value
	if t.isFadingIn {
		toOp.ColorScale.ScaleAlpha(t.alpha)
		screen.DrawImage(fromImg, toOp)

		fromOp.ColorScale.ScaleAlpha(1.0 - t.alpha)
		screen.DrawImage(fromImg, fromOp)
	} else {
		fromOp.ColorScale.ScaleAlpha(t.alpha)
		screen.DrawImage(fromImg, fromOp)

		toOp.ColorScale.ScaleAlpha(1.0 - t.alpha)
		screen.DrawImage(toImg, toOp)
	}
	t.frameUpdated = false
}

type SlideTransition[T any] struct {
	BaseTransition[T]
	factor       float64 // factor used for the slide-in/slide-out effect
	direction    SlideDirection
	offset       float64
	frameUpdated bool
}

type SlideDirection int

const (
	LeftToRight SlideDirection = iota
	RightToLeft
	TopToBottom
	BottomToTop
)

func NewSlideTransition[T any](direction SlideDirection, factor float64) *SlideTransition[T] {
	return &SlideTransition[T]{
		direction: direction,
		factor:    factor,
	}
}

// Start starts the transition from the given "from" scene to the given "to" scene
func (t *SlideTransition[T]) Start(fromScene Scene[T], toScene Scene[T]) {
	t.BaseTransition.Start(fromScene, toScene)
	t.offset = 0
}

// Update updates the transition state
func (t *SlideTransition[T]) Update() error {
	if !t.frameUpdated {
		// Update the offset value based on the current state of the transition
		if t.offset >= 1.0 {
			t.offset = 1.0
			t.End()
		} else {
			t.offset += t.factor
		}
		t.frameUpdated = true
	}

	// Update the scenes
	return t.BaseTransition.Update()
}

// Draw draws the transition effect
func (t *SlideTransition[T]) Draw(screen *ebiten.Image) {
	toImg, fromImg := PreDraw(screen.Bounds(), t.fromScene, t.toScene)
	toOp, fromOp := &ebiten.DrawImageOptions{}, &ebiten.DrawImageOptions{}

	w, h := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())

	var x, y float64

	switch t.direction {
	case LeftToRight:
		x = w * t.offset
		fromOp.GeoM.Translate(x, 0)
		toOp.GeoM.Translate(x-w, 0)
	case RightToLeft:
		x = w * (1 - t.offset)
		fromOp.GeoM.Translate(x-w, 0)
		toOp.GeoM.Translate(x, 0)
	case TopToBottom:
		y = h * t.offset
		fromOp.GeoM.Translate(0, y)
		toOp.GeoM.Translate(0, y-h)
	case BottomToTop:
		y = h * (1 - t.offset)
		fromOp.GeoM.Translate(0, y-h)
		toOp.GeoM.Translate(0, y)
	}
	screen.DrawImage(toImg, toOp)
	screen.DrawImage(fromImg, fromOp)
	t.frameUpdated = false
}

// Timed Variants of the transition

func NewTicksTimedFadeTransition[T any](duration time.Duration) *FadeTransition[T] {
	return NewFadeTransition[T](float32(DurationToFactor(float64(ebiten.TPS()), duration)))
}

type TimedFadeTransition[T any] struct {
	FadeTransition[T]
	initialTime time.Time
	duration    time.Duration
}

func NewDurationTimedFadeTransition[T any](duration time.Duration) *TimedFadeTransition[T] {
	return &TimedFadeTransition[T]{
		duration:       duration,
		FadeTransition: *NewFadeTransition[T](0.),
	}
}

func (t *TimedFadeTransition[T]) Start(fromScene, toScene Scene[T]) {
	t.FadeTransition.Start(fromScene, toScene)
	t.initialTime = Clock.Now()
}

func (t *TimedFadeTransition[T]) Update() error {
	if !t.frameUpdated {
		// Update the alpha value based on the current state of the transition
		if t.isFadingIn {
			t.alpha = float32(CalculateProgress(t.initialTime, t.duration/2))
			if t.alpha >= 1.0 {
				t.alpha = 1.0
				t.isFadingIn = false
			}
		} else {
			t.alpha = 1 - float32(CalculateProgress(t.initialTime.Add(t.duration/2), t.duration/2))
			if t.alpha <= 0.0 {
				t.alpha = 0.0
				t.End()
			}
		}
		t.frameUpdated = true
	}

	// Update the scenes
	return t.BaseTransition.Update()

}

func NewTicksTimedSlideTransition[T any](direction SlideDirection, duration time.Duration) *SlideTransition[T] {
	return NewSlideTransition[T](direction, DurationToFactor(float64(ebiten.TPS()), duration))
}

type TimedSlideTransition[T any] struct {
	SlideTransition[T]
	initialTime time.Time
	duration    time.Duration
}

func NewDurationTimedSlideTransition[T any](direction SlideDirection, duration time.Duration) *TimedSlideTransition[T] {
	return &TimedSlideTransition[T]{
		duration:        duration,
		SlideTransition: *NewSlideTransition[T](direction, 0.),
	}
}

func (t *TimedSlideTransition[T]) Start(fromScene, toScene Scene[T]) {
	t.SlideTransition.Start(fromScene, toScene)
	t.initialTime = Clock.Now()
}

func (t *TimedSlideTransition[T]) Update() error {
	if !t.frameUpdated {
		// Update the offset value based on the current state of the transition
		if t.offset >= 1.0 {
			t.offset = 1.0
			t.End()
		} else {
			t.offset = CalculateProgress(t.initialTime, t.duration)
		}
		t.frameUpdated = true
	}

	// Update the scenes
	return t.BaseTransition.Update()

}
