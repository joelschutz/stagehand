package stagehand

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

type MockScene struct {
	loadCalled    bool
	unloadCalled  bool
	updateCalled  bool
	drawCalled    bool
	layoutCalled  bool
	unloadReturns int
}

func (m *MockScene) Load(state int, sm *SceneManager[int]) {
	m.loadCalled = true
	m.unloadReturns = state
}

func (m *MockScene) Unload() int {
	m.unloadCalled = true
	return m.unloadReturns
}

func (m *MockScene) Update() error {
	m.updateCalled = true
	return nil
}

func (m *MockScene) Draw(screen *ebiten.Image) {
	m.drawCalled = true
}

func (m *MockScene) Layout(w, h int) (int, int) {
	m.layoutCalled = true
	return w, h
}

type MockTransition[T any] struct {
	fromScene   Scene[T]
	toScene     Scene[T]
	startCalled bool
	state       T
}

func NewMockTransition[T any]() *MockTransition[T] { return &MockTransition[T]{} }

func (t *MockTransition[T]) Start(fromScene, toScene Scene[T]) {
	t.fromScene = fromScene
	t.toScene = toScene
	t.startCalled = true
}

func (t *MockTransition[T]) End() {}

func (t *MockTransition[T]) Update() error { return nil }

func (t *MockTransition[T]) Draw(screen *ebiten.Image) {}

func (t *MockTransition[T]) Load(state T, sm *SceneManager[T]) { t.state = state }

func (t *MockTransition[T]) Unload() T { return t.state }

func (t *MockTransition[T]) Layout(w, h int) (int, int) { return w, h }

func TestSceneManager_SwitchTo(t *testing.T) {
	sm := NewSceneManager[int](&MockScene{}, 0)
	mockScene := &MockScene{}
	sm.SwitchTo(mockScene)
	assert.True(t, mockScene.loadCalled)
	assert.True(t, sm.current == mockScene)
}

func TestSceneManager_SwitchWithTransition(t *testing.T) {
	sm := NewSceneManager[int](&MockScene{}, 0)
	mockScene := &MockScene{}
	mockTransition := &MockTransition[int]{}
	sm.SwitchWithTransition(mockScene, mockTransition)
	assert.True(t, mockTransition.startCalled)
	assert.True(t, sm.current == mockTransition)
}

func TestSceneManager_Update(t *testing.T) {
	sm := NewSceneManager[int](&MockScene{}, 0)
	mockScene := &MockScene{}
	sm.SwitchTo(mockScene)
	sm.Update()
	assert.True(t, mockScene.updateCalled)
}

func TestSceneManager_Draw(t *testing.T) {
	sm := NewSceneManager[int](&MockScene{}, 0)
	mockScene := &MockScene{}
	screen := &ebiten.Image{}
	sm.SwitchTo(mockScene)
	sm.Draw(screen)
	assert.True(t, mockScene.drawCalled)
}

func TestSceneManager_Layout(t *testing.T) {
	sm := NewSceneManager[int](&MockScene{}, 0)
	mockScene := &MockScene{}
	w, h := 800, 600
	sm.SwitchTo(mockScene)
	rw, rh := sm.Layout(w, h)
	assert.Equal(t, w, rw)
	assert.Equal(t, h, rh)
	assert.True(t, mockScene.layoutCalled)
}

func TestSceneManager_Load_Unload(t *testing.T) {
	sm := NewSceneManager[int](&MockScene{}, 42)
	mockScene := &MockScene{}
	sm.SwitchTo(mockScene)
	unloaded := sm.current.Unload()
	assert.True(t, mockScene.unloadCalled)
	assert.Equal(t, 42, unloaded)
}
