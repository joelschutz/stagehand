package stagehand

import ebiten "github.com/hajimehoshi/ebiten/v2"

type SceneManager[T any] struct {
	current ProtoScene[T]
}

func NewSceneManager[T any](scene Scene[T, *SceneManager[T]], state T) *SceneManager[T] {
	s := &SceneManager[T]{current: scene}
	scene.Load(state, s)
	return s
}

// Scene Switching
func (s *SceneManager[T]) SwitchTo(scene Scene[T, *SceneManager[T]]) {
	if c, ok := s.current.(Scene[T, *SceneManager[T]]); ok {
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