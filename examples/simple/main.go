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

type BaseScene struct {
	bounds image.Rectangle
	count  State
	sm     *stagehand.SceneManager[State]
}

func (s *BaseScene) Layout(w, h int) (int, int) {
	s.bounds = image.Rect(0, 0, w, h)
	return w, h
}

func (s *BaseScene) Load(st State, sm stagehand.SceneController[State]) {
	s.count = st
	s.sm = sm.(*stagehand.SceneManager[State])
}

func (s *BaseScene) Unload() State {
	return s.count
}

type FirstScene struct {
	BaseScene
}

func (s *FirstScene) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.count++
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		s.sm.SwitchTo(&SecondScene{})
	}
	return nil
}

func (s *FirstScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 0, 0, 255}) // Fill Red
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Count: %v, WindowSize: %s", s.count, s.bounds.Max), s.bounds.Dx()/2, s.bounds.Dy()/2)
}

type SecondScene struct {
	BaseScene
}

func (s *SecondScene) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.count--
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		s.sm.SwitchWithTransition(&ThirdScene{}, stagehand.NewFadeTransition[State](.05))
	}
	return nil
}

func (s *SecondScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 255, 255}) // Fill Blue
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Count: %v, WindowSize: %s", s.count, s.bounds.Max), s.bounds.Dx()/2, s.bounds.Dy()/2)
}

type ThirdScene struct {
	BaseScene
}

func (s *ThirdScene) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.count *= 2
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		s.sm.SwitchWithTransition(&FirstScene{}, stagehand.NewSlideTransition[State](stagehand.RightToLeft, .05))
	}
	return nil
}

func (s *ThirdScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 255, 0, 255}) // Fill Green
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Count: %v, WindowSize: %s", s.count, s.bounds.Max), s.bounds.Dx()/2, s.bounds.Dy()/2)
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("My Game")
	ebiten.SetWindowResizable(true)

	state := State(10)

	s := &FirstScene{}
	sm := stagehand.NewSceneManager[State](s, state)

	if err := ebiten.RunGame(sm); err != nil {
		log.Fatal(err)
	}
}
