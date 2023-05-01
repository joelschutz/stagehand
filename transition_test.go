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

func TestFadeTransition_UpdateOncePerFrame(t *testing.T) {
	var value float32 = .6
	from := &MockScene{}
	to := &MockScene{}
	trans := NewFadeTransition[int](value)
	trans.Start(from, to)

	err := trans.Update()
	assert.NoError(t, err)
	assert.Equal(t, value, trans.alpha)
	assert.True(t, trans.frameUpdated)

	err = trans.Update()
	assert.NoError(t, err)
	assert.Equal(t, value, trans.alpha)
}

func TestFadeTransition_Update(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	trans := NewFadeTransition[int](.5)
	trans.Start(from, to)
	sm := NewSceneManager[int](trans, 0)

	err := sm.Update()
	assert.NoError(t, err)
	assert.Equal(t, float32(.5), trans.alpha)
	assert.True(t, trans.isFadingIn)

	trans.frameUpdated = false

	err = sm.Update()
	assert.NoError(t, err)
	assert.Equal(t, float32(1.0), trans.alpha)
	assert.False(t, trans.isFadingIn)

	trans.frameUpdated = false

	err = sm.Update()
	assert.NoError(t, err)
	assert.Equal(t, float32(.5), trans.alpha)
	assert.False(t, trans.isFadingIn)

	trans.frameUpdated = false

	err = sm.Update()
	assert.NoError(t, err)
	assert.Equal(t, float32(0), trans.alpha)
	assert.Equal(t, to, sm.current)

}

func TestFadeTransition_Start(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	trans := NewFadeTransition[int](.5)
	trans.Start(from, to)

	assert.Equal(t, from, trans.fromScene)
	assert.Equal(t, to, trans.toScene)
	assert.Equal(t, float32(.5), trans.factor)
	assert.True(t, trans.isFadingIn)
}

func TestFadeTransition_Draw(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	trans := NewFadeTransition[int](.5)
	trans.Start(from, to)

	trans.Update()
	trans.Draw(ebiten.NewImage(100, 100))
	assert.False(t, trans.frameUpdated)

	trans.isFadingIn = false
	trans.Draw(ebiten.NewImage(100, 100))
	assert.False(t, trans.frameUpdated)
}

func TestSlideTransition_UpdateOncePerFrame(t *testing.T) {
	var value float64 = .6
	from := &MockScene{}
	to := &MockScene{}
	trans := NewSlideTransition[int](RightToLeft, value)
	trans.Start(from, to)

	err := trans.Update()
	assert.NoError(t, err)
	assert.Equal(t, value, trans.offset)
	assert.True(t, trans.frameUpdated)

	err = trans.Update()
	assert.NoError(t, err)
	assert.Equal(t, value, trans.offset)
}

func TestSlideTransition_Update(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	variations := []SlideDirection{
		LeftToRight, RightToLeft, TopToBottom, BottomToTop,
	}

	for _, direction := range variations {
		trans := NewSlideTransition[int](direction, .5)
		trans.Start(from, to)
		sm := NewSceneManager[int](trans, 0)

		err := sm.Update()
		assert.NoError(t, err)
		assert.Equal(t, .5, trans.offset)

		trans.frameUpdated = false

		err = sm.Update()
		assert.NoError(t, err)
		assert.Equal(t, 1.0, trans.offset)

		trans.frameUpdated = false

		err = sm.Update()
		assert.NoError(t, err)
		assert.Equal(t, 1.0, trans.offset)
		assert.Equal(t, to, sm.current)
	}
}

func TestSlideTransition_Start(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	trans := NewSlideTransition[int](TopToBottom, .5)
	trans.Start(from, to)

	assert.Equal(t, from, trans.fromScene)
	assert.Equal(t, to, trans.toScene)
	assert.Equal(t, .5, trans.factor)
	assert.Equal(t, TopToBottom, trans.direction)
}

func TestSlideTransition_Draw(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	variations := []SlideDirection{
		LeftToRight, RightToLeft, TopToBottom, BottomToTop,
	}

	for _, direction := range variations {
		trans := NewSlideTransition[int](direction, .5)
		trans.Start(from, to)

		trans.Update()
		trans.Draw(ebiten.NewImage(100, 100))
		assert.False(t, trans.frameUpdated)
	}
}
