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

type State struct {
	Count        int
	OnTransition bool
}

type BaseScene struct {
	bounds image.Rectangle
	count  State
	sm     *stagehand.SceneManager[State]
}

func (s *BaseScene) Layout(w, h int) (int, int) {
	s.bounds = image.Rect(0, 0, w, h)
	return w, h
}

func (s *BaseScene) Load(st State, sm *stagehand.SceneManager[State]) {
	s.count = st
	s.sm = sm
}

func (s *BaseScene) Unload() State {
	return s.count
}

func (s *BaseScene) PreTransition(toScene stagehand.Scene[State, *stagehand.SceneManager[State]]) State {
	s.count.OnTransition = true
	return s.count
}

func (s *BaseScene) PostTransition(state State, fromScene stagehand.Scene[State, *stagehand.SceneManager[State]]) {
	s.count.OnTransition = false
}

type FirstScene struct {
	BaseScene
}

func (s *FirstScene) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.count.Count++
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		s.sm.SwitchWithTransition(&SecondScene{}, stagehand.NewSlideTransition[State, *stagehand.SceneManager[State]](stagehand.TopToBottom, .05))
	}
	return nil
}

func (s *FirstScene) Draw(screen *ebiten.Image) {
	if s.count.OnTransition {
		screen.Fill(color.RGBA{0, 0, 0, 255}) // Fill Black
	} else {
		screen.Fill(color.RGBA{255, 0, 0, 255}) // Fill Red
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Count: %v, WindowSize: %s", s.count.Count, s.bounds.Max), s.bounds.Dx()/2, s.bounds.Dy()/2)
}

type SecondScene struct {
	BaseScene
}

func (s *SecondScene) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.count.Count--
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		s.sm.SwitchWithTransition(&FirstScene{}, stagehand.NewSlideTransition[State, *stagehand.SceneManager[State]](stagehand.BottomToTop, .05))
	}
	return nil
}

func (s *SecondScene) Draw(screen *ebiten.Image) {
	if s.count.OnTransition {
		screen.Fill(color.RGBA{255, 255, 255, 255}) // Fill White
	} else {
		screen.Fill(color.RGBA{0, 0, 255, 255}) // Fill Blue
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Count: %v, WindowSize: %s", s.count.Count, s.bounds.Max), s.bounds.Dx()/2, s.bounds.Dy()/2)
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("My Game")
	ebiten.SetWindowResizable(true)

	state := State{Count: 10}

	s := &FirstScene{}
	sm := stagehand.NewSceneManager[State](s, state)

	if err := ebiten.RunGame(sm); err != nil {
		log.Fatal(err)
	}
}
