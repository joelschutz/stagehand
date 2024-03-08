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

type MockTransition[T any, M SceneController[T]] struct {
	fromScene   Scene[T, M]
	toScene     Scene[T, M]
	startCalled bool
}

func NewMockTransition[T any, M SceneController[T]]() *MockTransition[T, M] {
	return &MockTransition[T, M]{}
}

func (t *MockTransition[T, M]) Start(fromScene, toScene Scene[T, M], sm *SceneManager[T]) {
	t.fromScene = fromScene
	t.toScene = toScene
	t.startCalled = true
}

func (t *MockTransition[T, M]) End() {}

func (t *MockTransition[T, M]) Update() error { return nil }

func (t *MockTransition[T, M]) Draw(screen *ebiten.Image) {}

func (t *MockTransition[T, M]) Layout(w, h int) (int, int) { return w, h }

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
	mockTransition := &MockTransition[int, *SceneManager[int]]{}
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
	from := &MockScene{}
	to := &MockScene{}
	sm := NewSceneManager[int](from, 42)
	sm.SwitchTo(to)

	assert.True(t, to.loadCalled)
	assert.True(t, from.unloadCalled)
	assert.Equal(t, 42, sm.current.(Scene[int, *SceneManager[int]]).Unload())
}
