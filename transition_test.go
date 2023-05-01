package stagehand

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

type baseTransitionImplementation struct {
	BaseTransition[int]
}

func (b *baseTransitionImplementation) Draw(screen *ebiten.Image) {}

func TestBaseTransition_Update(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	trans := &baseTransitionImplementation{}
	trans.Start(from, to)

	err := trans.Update()
	assert.NoError(t, err)
	assert.True(t, from.updateCalled)
	assert.True(t, to.updateCalled)
}

func TestBaseTransition_Layout(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	trans := &baseTransitionImplementation{}
	trans.Start(from, to)

	sw, sh := trans.Layout(100, 100)
	assert.Equal(t, 100, sw)
	assert.Equal(t, 100, sh)

	assert.True(t, from.layoutCalled)
	assert.True(t, to.layoutCalled)
}

func TestBaseTransition_Load(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	trans := &baseTransitionImplementation{}
	trans.Start(from, to)
	trans.Load(42, &SceneManager[int]{})

	assert.True(t, to.loadCalled)
	assert.False(t, from.loadCalled)

}

func TestBaseTransition_Unload(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	trans := &baseTransitionImplementation{}
	trans.Start(from, to)

	trans.Unload()
	assert.True(t, from.unloadCalled)
	assert.False(t, to.unloadCalled)
}

func TestBaseTransition_End(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	trans := &baseTransitionImplementation{}
	trans.Start(from, to)
	sm := NewSceneManager[int](trans, 0)

	trans.End()
	assert.Equal(t, to, sm.current)
}

func TestBaseTransition_Start(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	trans := &baseTransitionImplementation{}
	trans.Start(from, to)

	assert.Equal(t, from, trans.fromScene)
	assert.Equal(t, to, trans.toScene)
}
