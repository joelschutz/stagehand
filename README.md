# Stagehand

![Codecov](https://img.shields.io/codecov/c/gh/joelschutz/stagehand?token=3BD1FTRKUC)
[![issues](https://img.shields.io/github/issues/joelschutz/stagehand)](https://github.com/joelschutz/stagehand/issues)
[![license](https://img.shields.io/github/license/joelschutz/stagehand)](https://github.com/joelschutz/stagehand/blob/main/LICENSE)

Stagehand is a scene manager library for [Ebitengine](https://ebitengine.org) that makes it easy to manage game scenes and state. With Stagehand, you can quickly create and transition between scenes, allowing you to focus on the content and gameplay of your game.


## Installation

To use Stagehand, simply import it into your project:

```
import "github.com/joelschutz/stagehand"
```

## Features

- Lightweight and easy-to-use scene management for Ebitengine
- Simple and flexible API for creating and switching between scenes
- Managed type-safe states with the power of generics
- Built-in support for transition effects between scenes
- Supports custom transition effects by implementing the `SceneTransition` interface

## Usage

To use Stagehand, you first need to create a struct that implements the `Scene` interface:

```go
type MyState struct {
    // your state data
}

type MyScene struct {
    // your scene fields
}

func (s *MyScene) Update() error {
    // your update code
}

func (s *MyScene) Draw(screen *ebiten.Image) {
    // your draw code
}

func (s *MyScene) Load(state MyState ,manager *stagehand.SceneManager) {
    // your load code
}

func (s *MyScene) Unload() MyState {
    // your unload code
}
```

Then, create an instance of the `SceneManager` passing initial scene and state.

```go
func main() {
    // ...
    scene1 := &MyScene{}
    state := MyState{}
    manager := stagehand.NewSceneManager[MyState](scene1, state)

    if err := ebiten.RunGame(sm); err != nil {
		log.Fatal(err)
	}
}
```

### Examples

We provide some example code so you can start fast:

- [Simple Example](https://github.com/stagehand/blob/master/examples/simple/main.go)
- [Timed Transition Example](https://github.com/stagehand/blob/master/examples/timed/main.go)
- [Transition Awareness Example](https://github.com/stagehand/blob/master/examples/aware/main.go)

## Transitions

You can switch scenes by calling `SwitchTo` method on the `SceneManager` giving the scene instance you wanna switch to.

```go
func (s *MyScene) Update() error {
    // ...
    scene2 := &OtherScene{}
    s.manager.SwitchTo(scene2)

    // ...
}
```

You can use the `SwitchWithTransition` method to switch between scenes with a transition effect. Stagehand provides two built-in transition effects: `FadeTransition` and `SlideTransition`.

### Fade Transition

The `FadeTransition` will fade out the current scene while fading in the new scene.

```go
func (s *MyScene) Update() error {
    // ...
    scene2 := &OtherScene{}
    s.manager.SwitchWithTransition(scene2. stagehand.NewFadeTransition(.05))

    // ...
}
```

In this example, the `FadeTransition` will fade 5% every frame. There is also the option for a timed transition using `NewTicksTimedFadeTransition`(for a ticks based timming) or `NewDurationTimedFadeTransition`(for a real-time based timming).

### Slide Transition

The `SlideTransition` will slide out the current scene and slide in the new scene.

```go
func (s *MyScene) Update() error {
    // ...
    scene2 := &OtherScene{}
    s.manager.SwitchWithTransition(scene2. stagehand.NewSlideTransition(stagehand.LeftToRight, .05))

    // ...
}
```

In this example, the `SlideTransition` will slide in the new scene from the left 5% every frame. There is also the option for a timed transition using `NewTicksTimedSlideTransition`(for a ticks based timming) or `NewDurationTimedSlideTransition`(for a real-time based timming).

### Custom Transitions

You can also define your own transition, simply implement the `SceneTransition` interface, we provide a helper `BaseTransition` that you can use like this:

```go
type MyTransition struct {
    stagehand.BaseTransition
	progress float64 // An example factor
}

func (t *MyTransition) Start(from, to stagehand.Scene[MyState], sm *SceneManager[MyState]) {
    // Start the transition from the "from" scene to the "to" scene here
    t.BaseTransition.Start(fromScene, toScene, sm)
    t.progress = 0
}

func (t *MyTransition) Update() error {
	// Update the progress of the transition
    t.progress += 0.01
	return t.BaseTransition.Update()
}

func (t *MyTransition) Draw(screen *ebiten.Image) {
	// Optionally you can use a helper function to render each scene frame
	toImg, fromImg := stagehand.PreDraw(screen.Bounds(), t.fromScene, t.toScene)

    // Draw transition effect here
}

```

### Transition Awareness

When a scene is transitioned, the `Load` and `Unload` methods are called **twice** for the destination and original scenes respectively. Once at the start and again at the end of the transition. This behavior can be changed for additional control by implementing the `TransitionAwareScene` interface.

```go
func (s *MyScene) PreTransition(destination Scene[MyState]) MyState  {
    // Runs before new scene is loaded
}

func (s *MyScene) PostTransition(lastState MyState, original Scene[MyState]) {
    // Runs when old scene is unloaded
}
```

With this you can insure that those methods are only called once on transitions and can control your scenes at each point of the transition. The execution order will be:

```shell
PreTransition Called on old scene
Load Called on new scene
Updated old scene
Updated new scene
...
Updated old scene
Updated new scene
Unload Called on old scene
PostTransition Called on new scene
```

## Contribution

Contributions are welcome! If you find a bug or have a feature request, please open an issue on GitHub. If you would like to contribute code, please fork the repository and submit a pull request.

Before submitting a pull request, please make sure to run the tests:

```
go test ./...
```

## License

Stagehand is released under the [MIT License](https://github.com/stagehand/blob/master/LICENSE).
