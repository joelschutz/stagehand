package stagehand

type SceneTransitionTrigger int

// A Directive is a struct that represents how a scene should be transitioned
type Directive[T any] struct {
	Dest       Scene[T, *SceneDirector[T]]
	Transition SceneTransition[T, *SceneDirector[T]]
	Trigger    SceneTransitionTrigger
}

// A SceneDirector is a struct that manages the transitions between scenes
type SceneDirector[T any] struct {
	SceneManager[T]
	RuleSet map[Scene[T, *SceneDirector[T]]][]Directive[T]
}

func NewSceneDirector[T any](scene Scene[T, *SceneDirector[T]], state T, RuleSet map[Scene[T, *SceneDirector[T]]][]Directive[T]) *SceneDirector[T] {
	s := &SceneDirector[T]{RuleSet: RuleSet}
	s.current = scene
	scene.Load(state, s)
	return s
}

// ProcessTrigger finds if a transition should be triggered
func (d *SceneDirector[T]) ProcessTrigger(trigger SceneTransitionTrigger) {
	for _, directive := range d.RuleSet[d.current.(Scene[T, *SceneDirector[T]])] {
		if directive.Trigger == trigger {
			if directive.Transition != nil {
				// With transition
				// Equivalent to SwitchWithTransition
				sc := d.current.(Scene[T, *SceneDirector[T]])
				directive.Transition.Start(sc, directive.Dest, &d.SceneManager)
				if c, ok := sc.(TransitionAwareScene[T, *SceneDirector[T]]); ok {
					directive.Dest.Load(c.PreTransition(directive.Dest), d)
				} else {
					directive.Dest.Load(sc.Unload(), d)
				}
				d.current = directive.Transition
			} else {
				// No transition
				// Equivalent to SwitchTo
				if c, ok := d.current.(Scene[T, *SceneDirector[T]]); ok {
					directive.Dest.Load(c.Unload(), d)
					d.current = directive.Dest
				}
			}

		}
	}
}
