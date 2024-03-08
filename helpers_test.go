package stagehand

import (
	"image"
	"testing"
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

func TestMaxInt(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"a bigger", args{2, 1}, 2},
		{"b bigger", args{1, 2}, 2},
		{"equal", args{1, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxInt(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("MaxInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPreDraw(t *testing.T) {
	from := &MockScene{}
	to := &MockScene{}

	toImg, fromImg := PreDraw[int, *SceneManager[int]](image.Rect(0, 0, 10, 10), from, to)

	assert.True(t, from.drawCalled)
	assert.True(t, to.drawCalled)
	assert.IsType(t, &ebiten.Image{}, fromImg)
	assert.IsType(t, &ebiten.Image{}, toImg)
}

func TestDurationToFactor(t *testing.T) {
	type args struct {
		frequency float64
		duration  time.Duration
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"1s 1hz", args{1, time.Second}, 1},
		{"1s 2hz", args{2, time.Second}, .5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DurationToFactor(tt.args.frequency, tt.args.duration); got != tt.want {
				t.Errorf("DurationToFactor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateProgress(t *testing.T) {
	type args struct {
		initialTime time.Time
		duration    time.Duration
	}
	tests := []struct {
		name     string
		args     args
		expected func(float64) bool
	}{
		{"1s - now", args{time.Now(), time.Second}, func(f float64) bool { return f >= 0 && f <= .1 }},
		{"1s - 1s ago", args{time.Now().Add(-time.Second), time.Second}, func(f float64) bool { return f >= 1 && f <= 1.1 }},
		{"1s - 2s ago", args{time.Now().Add(-2 * time.Second), time.Second}, func(f float64) bool { return f >= 2 && f <= 2.1 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateProgress(tt.args.initialTime, tt.args.duration); !tt.expected(got) {
				t.Errorf("CalculateProgress() = %v", got)
			}
		})
	}
}
