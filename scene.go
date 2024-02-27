package stagehand

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type ProtoScene[T any] interface {
	ebiten.Game
}

type Scene[T any] interface {
	ProtoScene[T]
	Load(T, *SceneManager[T]) // Runs when scene is first started, must keep state and SceneManager
	Unload() T                // Runs when scene is discarted, must return last state
}
type TransitionAwareScene[T any] interface {
	Scene[T]
	PreTransition(Scene[T]) T   // Runs before Load, must return last state
	PostTransition(T, Scene[T]) // Runs when scene is fully loaded
}

type SceneManager[T any] struct {
	current ProtoScene[T]
}

func NewSceneManager[T any](scene Scene[T], state T) *SceneManager[T] {
	s := &SceneManager[T]{current: scene}
	scene.Load(state, s)
	return s
}

// Scene Switching
func (s *SceneManager[T]) SwitchTo(scene Scene[T]) {
	if c, ok := s.current.(Scene[T]); ok {
		scene.Load(c.Unload(), s)
		s.current = scene
	}
}

// Ebiten Interface
func (s *SceneManager[T]) Update() error {
	return s.current.Update()
}

func (s *SceneManager[T]) Draw(screen *ebiten.Image) {
	s.current.Draw(screen)
}

func (s *SceneManager[T]) Layout(w, h int) (int, int) {
	return s.current.Layout(w, h)
}
