package stagehand

type SceneTransitionTrigger int

// A Directive is a struct that represents how a scene should be transitioned
type Directive[T any] struct {
	Dest       Scene[T]
	Transition SceneTransition[T]
	Trigger    SceneTransitionTrigger
}

// A SceneDirector is a struct that manages the transitions between scenes
type SceneDirector[T any] struct {
	SceneManager[T]
	RuleSet map[Scene[T]][]Directive[T]
}

func NewSceneDirector[T any](scene Scene[T], state T, RuleSet map[Scene[T]][]Directive[T]) *SceneDirector[T] {
	s := &SceneDirector[T]{RuleSet: RuleSet}
	s.current = scene
	scene.Load(state, s)
	return s
}

// ProcessTrigger finds if a transition should be triggered
func (d *SceneDirector[T]) ProcessTrigger(trigger SceneTransitionTrigger) {
	if prevTransition, ok := d.current.(SceneTransition[T]); ok {
		// previous transition is still running, if related to trigger end it to start the next transition
		isTransitionTrigger := false
		for _, directiveList := range d.RuleSet {
			for _, directive := range directiveList {
				if directive.Trigger == trigger && prevTransition == directive.Transition {
					isTransitionTrigger = true
					break
				}
			}
		}
		if isTransitionTrigger {
			prevTransition.End()
			return
		}
	}

	for _, directive := range d.RuleSet[d.current.(Scene[T])] {
		if directive.Trigger == trigger {
			if directive.Transition != nil {
				// With transition
				// Equivalent to SwitchWithTransition
				sc := d.current.(Scene[T])
				directive.Transition.Start(sc, directive.Dest, d)
				if c, ok := sc.(TransitionAwareScene[T]); ok {
					directive.Dest.Load(c.PreTransition(directive.Dest), d)
				} else {
					directive.Dest.Load(sc.Unload(), d)
				}
				d.current = directive.Transition
			} else {
				// No transition
				// Equivalent to SwitchTo
				if c, ok := d.current.(Scene[T]); ok {
					directive.Dest.Load(c.Unload(), d)
					d.current = directive.Dest
				}
			}

		}
	}
}

func (d *SceneDirector[T]) ReturnFromTransition(scene, origin Scene[T]) {
	if c, ok := scene.(TransitionAwareScene[T]); ok {
		c.PostTransition(origin.Unload(), origin)
	} else {
		scene.Load(origin.Unload(), d)
	}
	d.current = scene
}
