package stagehand

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type SceneTransition[T any, M SceneController[T]] interface {
	ProtoScene[T]
	Start(fromScene, toScene Scene[T, M], sm *SceneManager[T])
	End()
}

type BaseTransition[T any, M SceneController[T]] struct {
	fromScene Scene[T, M]
	toScene   Scene[T, M]
	sm        *SceneManager[T]
}

func (t *BaseTransition[T, M]) Start(fromScene, toScene Scene[T, M], sm *SceneManager[T]) {
	t.fromScene = fromScene
	t.toScene = toScene
	t.sm = sm
}

// Update updates the transition state
func (t *BaseTransition[T, M]) Update() error {
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
func (t *BaseTransition[T, M]) Layout(outsideWidth, outsideHeight int) (int, int) {
	sw, sh := t.fromScene.Layout(outsideWidth, outsideHeight)
	tw, th := t.toScene.Layout(outsideWidth, outsideHeight)

	return MaxInt(sw, tw), MaxInt(sh, th)
}

func (s *SceneManager[T]) ReturnFromTransition(scene, orgin Scene[T, *SceneManager[T]]) {
	if c, ok := scene.(TransitionAwareScene[T, *SceneManager[T]]); ok {
		c.PostTransition(orgin.Unload(), orgin)
	} else {
		scene.Load(orgin.Unload(), s)
	}
	s.current = scene
}

func (s *SceneManager[T]) SwitchWithTransition(scene Scene[T, *SceneManager[T]], transition SceneTransition[T, *SceneManager[T]]) {
	sc := s.current.(Scene[T, *SceneManager[T]])
	transition.Start(sc, scene, s)
	if c, ok := sc.(TransitionAwareScene[T, *SceneManager[T]]); ok {
		scene.Load(c.PreTransition(scene), s)
	} else {
		scene.Load(sc.Unload(), s)
	}
	s.current = transition
}

// Ends transition to the next scene
func (t *BaseTransition[T, M]) End() {
	t.sm.ReturnFromTransition(t.toScene.(Scene[T, *SceneManager[T]]), t.fromScene.(Scene[T, *SceneManager[T]]))
}

type FadeTransition[T any, M SceneController[T]] struct {
	BaseTransition[T, M]
	factor       float32 // factor used for the fade-in/fade-out effect
	alpha        float32 // alpha value used for the fade-in/fade-out effect
	isFadingIn   bool    // whether the transition is currently fading in or out
	frameUpdated bool
}

func NewFadeTransition[T any, M SceneController[T]](factor float32) *FadeTransition[T, M] {
	return &FadeTransition[T, M]{
		factor: factor,
	}
}

// Start starts the transition from the given "from" scene to the given "to" scene
func (t *FadeTransition[T, M]) Start(fromScene, toScene Scene[T, M], sm *SceneManager[T]) {
	t.BaseTransition.Start(fromScene, toScene, sm)
	t.alpha = 0
	t.isFadingIn = true
}

// Update updates the transition state
func (t *FadeTransition[T, M]) Update() error {
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
func (t *FadeTransition[T, M]) Draw(screen *ebiten.Image) {
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

type SlideTransition[T any, M SceneController[T]] struct {
	BaseTransition[T, M]
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

func NewSlideTransition[T any, M SceneController[T]](direction SlideDirection, factor float64) *SlideTransition[T, M] {
	return &SlideTransition[T, M]{
		direction: direction,
		factor:    factor,
	}
}

// Start starts the transition from the given "from" scene to the given "to" scene
func (t *SlideTransition[T, M]) Start(fromScene Scene[T, M], toScene Scene[T, M], sm *SceneManager[T]) {
	t.BaseTransition.Start(fromScene, toScene, sm)
	t.offset = 0
}

// Update updates the transition state
func (t *SlideTransition[T, M]) Update() error {
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
func (t *SlideTransition[T, M]) Draw(screen *ebiten.Image) {
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

func NewTicksTimedFadeTransition[T any, M SceneController[T]](duration time.Duration) *FadeTransition[T, M] {
	return NewFadeTransition[T, M](float32(DurationToFactor(float64(ebiten.TPS()), duration)))
}

type TimedFadeTransition[T any, M SceneController[T]] struct {
	FadeTransition[T, M]
	initialTime time.Time
	duration    time.Duration
}

func NewDurationTimedFadeTransition[T any, M SceneController[T]](duration time.Duration) *TimedFadeTransition[T, M] {
	return &TimedFadeTransition[T, M]{
		duration:       duration,
		FadeTransition: *NewFadeTransition[T, M](0.),
	}
}

func (t *TimedFadeTransition[T, M]) Start(fromScene, toScene Scene[T, M], sm *SceneManager[T]) {
	t.FadeTransition.Start(fromScene, toScene, sm)
	t.initialTime = Clock.Now()
}

func (t *TimedFadeTransition[T, M]) Update() error {
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

func NewTicksTimedSlideTransition[T any, M SceneController[T]](direction SlideDirection, duration time.Duration) *SlideTransition[T, M] {
	return NewSlideTransition[T, M](direction, DurationToFactor(float64(ebiten.TPS()), duration))
}

type TimedSlideTransition[T any, M SceneController[T]] struct {
	SlideTransition[T, M]
	initialTime time.Time
	duration    time.Duration
}

func NewDurationTimedSlideTransition[T any, M SceneController[T]](direction SlideDirection, duration time.Duration) *TimedSlideTransition[T, M] {
	return &TimedSlideTransition[T, M]{
		duration:        duration,
		SlideTransition: *NewSlideTransition[T, M](direction, 0.),
	}
}

func (t *TimedSlideTransition[T, M]) Start(fromScene, toScene Scene[T, M], sm *SceneManager[T]) {
	t.SlideTransition.Start(fromScene, toScene, sm)
	t.initialTime = Clock.Now()
}

func (t *TimedSlideTransition[T, M]) Update() error {
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
