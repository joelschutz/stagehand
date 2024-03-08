package stagehand

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type ProtoScene[T any] interface {
	ebiten.Game
}

type SceneController[T any] interface {
	*SceneManager[T] | *SceneDirector[T]
}

type Scene[T any, M SceneController[T]] interface {
	ProtoScene[T]
	Load(T, M) // Runs when scene is first started, must keep state and SceneManager
	Unload() T // Runs when scene is discarted, must return last state
}

type TransitionAwareScene[T any, M SceneController[T]] interface {
	Scene[T, M]
	PreTransition(Scene[T, M]) T   // Runs before new scene is loaded, must return last state
	PostTransition(T, Scene[T, M]) // Runs when old scene is unloaded
}
