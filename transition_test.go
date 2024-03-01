package stagehand

import (
	"fmt"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

type MockClock struct {
	currentTime time.Time
}

func (m *MockClock) Now() time.Time                  { return m.currentTime }
func (m *MockClock) Sleep(d time.Duration)           { m.currentTime = m.currentTime.Add(d) }
func (m *MockClock) Since(t time.Time) time.Duration { return m.currentTime.Sub(t) }
func (m *MockClock) Until(t time.Time) time.Duration { return t.Sub(m.currentTime) }

type baseTransitionImplementation struct {
	BaseTransition[int]
}

func (b *baseTransitionImplementation) Draw(screen *ebiten.Image) {}

type MockTransitionAwareScene struct {
	MockScene
	preTransitionCalled  bool
	postTransitionCalled bool
}

func (m *MockTransitionAwareScene) PreTransition(fromScene Scene[int]) int {
	m.preTransitionCalled = true
	return 0
}

func (m *MockTransitionAwareScene) PostTransition(state int, toScene Scene[int]) {
	m.postTransitionCalled = true
}

func TestBaseTransition_Update(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	trans := &baseTransitionImplementation{}
	trans.Start(from, to, nil)

	err := trans.Update()
	assert.NoError(t, err)
	assert.True(t, from.updateCalled)
	assert.True(t, to.updateCalled)
}

func TestBaseTransition_Layout(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	trans := &baseTransitionImplementation{}
	trans.Start(from, to, nil)

	sw, sh := trans.Layout(100, 100)
	assert.Equal(t, 100, sw)
	assert.Equal(t, 100, sh)

	assert.True(t, from.layoutCalled)
	assert.True(t, to.layoutCalled)
}

func TestBaseTransition_End(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	trans := &baseTransitionImplementation{}
	sm := NewSceneManager[int](from, 0)
	sm.SwitchWithTransition(to, trans)
	trans.End()

	fmt.Println(sm.current.(Scene[int]), to)
	assert.Equal(t, to, sm.current)
}

func TestBaseTransition_Start(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	trans := &baseTransitionImplementation{}
	trans.Start(from, to, nil)

	assert.Equal(t, from, trans.fromScene)
	assert.Equal(t, to, trans.toScene)
}

func TestBaseTransition_Awareness(t *testing.T) {
	from := &MockTransitionAwareScene{}
	to := &MockTransitionAwareScene{}
	sm := NewSceneManager[int](from, 0)
	trans := &baseTransitionImplementation{}
	sm.SwitchWithTransition(to, trans)

	assert.True(t, from.preTransitionCalled)
	assert.True(t, to.loadCalled)
	assert.False(t, from.unloadCalled)
	assert.False(t, to.postTransitionCalled)

	trans.End()
	assert.True(t, from.unloadCalled)
	assert.True(t, to.postTransitionCalled)
}

func TestFadeTransition_UpdateOncePerFrame(t *testing.T) {
	var value float32 = .6
	from := &MockScene{}
	to := &MockScene{}
	trans := NewFadeTransition[int](value)
	trans.Start(from, to, nil)

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
	sm := NewSceneManager[int](from, 0)
	sm.SwitchWithTransition(to, trans)

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
	trans.Start(from, to, nil)

	assert.Equal(t, from, trans.fromScene)
	assert.Equal(t, to, trans.toScene)
	assert.Equal(t, float32(.5), trans.factor)
	assert.True(t, trans.isFadingIn)
}

func TestFadeTransition_Draw(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	trans := NewFadeTransition[int](.5)
	trans.Start(from, to, nil)

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
	trans.Start(from, to, nil)

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
		sm := NewSceneManager[int](from, 0)
		sm.SwitchWithTransition(to, trans)

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
	trans.Start(from, to, nil)

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
		trans.Start(from, to, nil)

		trans.Update()
		trans.Draw(ebiten.NewImage(100, 100))
		assert.False(t, trans.frameUpdated)
	}
}

func TestTimedFadeTransition_Update(t *testing.T) {
	now := time.Now()
	Clock = &MockClock{currentTime: now}
	from := &MockScene{}
	to := &MockScene{}
	trans := NewDurationTimedFadeTransition[int](time.Second)
	sm := NewSceneManager[int](from, 0)
	sm.SwitchWithTransition(to, trans)

	// Should not update if no time passed
	err := sm.Update()
	assert.NoError(t, err)
	assert.Equal(t, float32(.0), trans.alpha)
	assert.True(t, trans.isFadingIn)

	trans.frameUpdated = false

	Clock.Sleep(time.Second / 4)
	err = sm.Update()
	assert.NoError(t, err)
	assert.Equal(t, float32(.5), trans.alpha)
	assert.True(t, trans.isFadingIn)

	trans.frameUpdated = false

	Clock.Sleep(time.Second / 4)
	err = sm.Update()
	assert.NoError(t, err)
	assert.Equal(t, float32(1.0), trans.alpha)
	assert.False(t, trans.isFadingIn)

	trans.frameUpdated = false

	Clock.Sleep(time.Second / 4)
	err = sm.Update()
	assert.NoError(t, err)
	assert.Equal(t, float32(.5), trans.alpha)
	assert.False(t, trans.isFadingIn)

	trans.frameUpdated = false

	Clock.Sleep(time.Second / 4)
	err = sm.Update()
	assert.NoError(t, err)
	assert.Equal(t, float32(0), trans.alpha)
	assert.Equal(t, to, sm.current)

}

func TestTimedFadeTransition_Start(t *testing.T) {
	now := time.Now()
	Clock = &MockClock{currentTime: now}
	from := &MockScene{}
	to := &MockScene{}
	trans := NewDurationTimedFadeTransition[int](time.Second)
	trans.Start(from, to, nil)

	assert.Equal(t, from, trans.fromScene)
	assert.Equal(t, to, trans.toScene)
	assert.Equal(t, now, trans.initialTime)
	assert.True(t, trans.isFadingIn)
}

func TestTimedSlideTransition_Update(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}
	variations := []SlideDirection{
		LeftToRight, RightToLeft, TopToBottom, BottomToTop,
	}

	for _, direction := range variations {
		Clock = &MockClock{currentTime: time.Now()}
		trans := NewDurationTimedSlideTransition[int](direction, time.Second)
		sm := NewSceneManager[int](from, 0)
		sm.SwitchWithTransition(to, trans)

		// Should not update if no time passed
		err := sm.Update()
		assert.NoError(t, err)
		assert.Equal(t, .0, trans.offset)

		trans.frameUpdated = false

		Clock.Sleep(time.Second / 2)
		err = sm.Update()
		assert.NoError(t, err)
		assert.Equal(t, .5, trans.offset)

		trans.frameUpdated = false

		Clock.Sleep(time.Second / 2)
		err = sm.Update()
		assert.NoError(t, err)
		assert.Equal(t, 1.0, trans.offset)

		trans.frameUpdated = false

		Clock.Sleep(time.Second / 2)
		err = sm.Update()
		assert.NoError(t, err)
		assert.Equal(t, 1.0, trans.offset)
		assert.Equal(t, to, sm.current)
	}
}

func TestTimedSlideTransition_Start(t *testing.T) {
	now := time.Now()
	Clock = &MockClock{currentTime: now}
	from := &MockScene{}
	to := &MockScene{}
	trans := NewDurationTimedSlideTransition[int](TopToBottom, time.Second)
	trans.Start(from, to, nil)

	assert.Equal(t, from, trans.fromScene)
	assert.Equal(t, to, trans.toScene)
	assert.Equal(t, now, trans.initialTime)
	assert.Equal(t, TopToBottom, trans.direction)
}
