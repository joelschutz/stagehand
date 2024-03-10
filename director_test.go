package stagehand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSceneDirector_NewSceneDirector(t *testing.T) {
	mockScene := &MockScene{}
	ruleSet := make(map[Scene[int]][]Directive[int])

	director := NewSceneDirector[int](mockScene, 1, ruleSet)

	assert.NotNil(t, director)
	assert.Equal(t, mockScene, director.current)
}

func TestSceneDirector_ProcessTrigger(t *testing.T) {
	mockScene := &MockScene{}
	mockScene2 := &MockScene{}
	ruleSet := make(map[Scene[int]][]Directive[int])

	director := NewSceneDirector[int](mockScene, 1, ruleSet)

	rule := Directive[int]{Dest: mockScene2, Trigger: 2}
	ruleSet[mockScene] = []Directive[int]{rule}

	// Call the ProcessTrigger method with wrong trigger
	director.ProcessTrigger(1)
	assert.NotEqual(t, rule.Dest, director.current)

	// Call the ProcessTrigger method with correct trigger
	director.ProcessTrigger(2)
	assert.Equal(t, rule.Dest, director.current)
}

func TestSceneDirector_ProcessTriggerWithTransition(t *testing.T) {
	mockScene := &MockScene{}
	mockTransition := &baseTransitionImplementation{}
	ruleSet := make(map[Scene[int]][]Directive[int])

	director := NewSceneDirector[int](mockScene, 1, ruleSet)

	rule := Directive[int]{Dest: &MockScene{}, Trigger: 2, Transition: mockTransition}
	ruleSet[mockScene] = []Directive[int]{rule}

	// Call the ProcessTrigger method with wrong trigger
	director.ProcessTrigger(1)
	assert.NotEqual(t, rule.Transition, director.current)

	// Call the ProcessTrigger method with correct trigger
	director.ProcessTrigger(2)
	assert.Equal(t, rule.Transition, director.current)

	rule.Transition.End()
	assert.Equal(t, rule.Dest, director.current)
}

func TestSceneDirector_ProcessTriggerWithTransitionAwareness(t *testing.T) {
	mockScene := &MockTransitionAwareScene{}
	mockTransition := &baseTransitionImplementation{}
	ruleSet := make(map[Scene[int]][]Directive[int])

	director := NewSceneDirector[int](mockScene, 1, ruleSet)

	rule := Directive[int]{Dest: &MockTransitionAwareScene{}, Trigger: 2, Transition: mockTransition}
	ruleSet[mockScene] = []Directive[int]{rule}

	// Call the ProcessTrigger method with wrong trigger
	director.ProcessTrigger(1)
	assert.NotEqual(t, rule.Transition, director.current)

	// Call the ProcessTrigger method with correct trigger
	director.ProcessTrigger(2)
	assert.Equal(t, rule.Transition, director.current)

	rule.Transition.End()
	assert.Equal(t, rule.Dest, director.current)
}
