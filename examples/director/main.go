package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joelschutz/stagehand"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type State int

const (
	Trigger stagehand.SceneTransitionTrigger = iota
)

type BaseScene struct {
	bounds image.Rectangle
	count  State
	active bool
	sm     *stagehand.SceneDirector[State]
}

func (s *BaseScene) Layout(w, h int) (int, int) {
	s.bounds = image.Rect(0, 0, w, h)
	return w, h
}

func (s *BaseScene) Load(st State, sm stagehand.SceneController[State]) {
	s.count = st
	s.active = true
	s.sm = sm.(*stagehand.SceneDirector[State])
}

func (s *BaseScene) Unload() State {
	s.active = false
	return s.count
}

type FirstScene struct {
	BaseScene
}

func (s *FirstScene) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.count++
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) && s.active {
		s.sm.ProcessTrigger(Trigger)
	}
	return nil
}

func (s *FirstScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 0, 0, 255}) // Fill Red
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Count: %v, WindowSize: %s\nActive? %v", s.count, s.bounds.Max, s.active), s.bounds.Dx()/2, s.bounds.Dy()/2)
}

type SecondScene struct {
	BaseScene
}

func (s *SecondScene) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.count--
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) && s.active {
		s.sm.ProcessTrigger(Trigger)
	}
	return nil
}

func (s *SecondScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 255, 255}) // Fill Blue
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Count: %v, WindowSize: %s\nActive? %v", s.count, s.bounds.Max, s.active), s.bounds.Dx()/2, s.bounds.Dy()/2)
}

type ThirdScene struct {
	BaseScene
}

func (s *ThirdScene) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.count++
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) && s.active {
		s.sm.ProcessTrigger(Trigger)
	}
	return nil
}

func (s *ThirdScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 255, 0, 255}) // Fill Green
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Count: %v, WindowSize: %s\nActive? %v", s.count, s.bounds.Max, s.active), s.bounds.Dx()/2, s.bounds.Dy()/2)
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("My Game")
	ebiten.SetWindowResizable(true)

	state := State(10)

	s1 := &FirstScene{}
	s2 := &SecondScene{}
	s3 := &ThirdScene{}
	trans := stagehand.NewSlideTransition[State](stagehand.BottomToTop, 0.02)
	trans2 := stagehand.NewSlideTransition[State](stagehand.TopToBottom, 0.02)
	trans3 := stagehand.NewSlideTransition[State](stagehand.LeftToRight, 0.02)
	rs := map[stagehand.Scene[State]][]stagehand.Directive[State]{
		s1: {
			stagehand.Directive[State]{
				Dest:       s2,
				Trigger:    Trigger,
				Transition: trans,
			},
		},
		s2: {
			stagehand.Directive[State]{
				Dest:       s3,
				Trigger:    Trigger,
				Transition: trans2,
			},
		},
		s3: {
			stagehand.Directive[State]{
				Dest:       s1,
				Trigger:    Trigger,
				Transition: trans3,
			},
		},
	}
	sm := stagehand.NewSceneDirector[State](s1, state, rs)

	if err := ebiten.RunGame(sm); err != nil {
		log.Fatal(err)
	}
}
