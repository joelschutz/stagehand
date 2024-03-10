package stagehand

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type ProtoScene[T any] interface {
	ebiten.Game
}

type SceneController[T any] interface {
	// *SceneManager[T] | *SceneDirector[T]
	ReturnFromTransition(scene, orgin Scene[T])
}

type Scene[T any] interface {
	ProtoScene[T]
	Load(T, SceneController[T]) // Runs when scene is first started, must keep state and SceneManager
	Unload() T                  // Runs when scene is discarted, must return last state
}

type TransitionAwareScene[T any] interface {
	Scene[T]
	PreTransition(Scene[T]) T   // Runs before new scene is loaded, must return last state
	PostTransition(T, Scene[T]) // Runs when old scene is unloaded
}
