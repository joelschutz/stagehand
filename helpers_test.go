package stagehand

import (
	"image"
	"testing"

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

	toImg, fromImg := PreDraw[int](image.Rect(0, 0, 10, 10), from, to)

	assert.True(t, from.drawCalled)
	assert.True(t, to.drawCalled)
	assert.IsType(t, &ebiten.Image{}, fromImg)
	assert.IsType(t, &ebiten.Image{}, toImg)
}
