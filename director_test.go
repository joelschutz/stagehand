package stagehand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type DirectidScene struct {
	MockScene
}

func (m *DirectidScene) Load(state int, sm *SceneDirector[int]) {
	m.loadCalled = true
	m.unloadReturns = state
}

func TestSceneDirector_NewSceneDirector(t *testing.T) {
	mockScene := &DirectidScene{}
	ruleSet := make(map[Scene[int, *SceneDirector[int]]][]Directive[int])

	director := NewSceneDirector[int](mockScene, 1, ruleSet)

	assert.NotNil(t, director)
	assert.Equal(t, mockScene, director.current)
}

func TestSceneDirector_ProcessTrigger(t *testing.T) {
	mockScene := &DirectidScene{}
	mockScene2 := &DirectidScene{}
	ruleSet := make(map[Scene[int, *SceneDirector[int]]][]Directive[int])

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
