package tests

import (
	"github.com/Nikola-Milovic/tog-plugin/math"
	"testing"
)

func TestVector_AngleTo(t *testing.T) {
	type fields struct {
		X float32
		Y float32
	}
	type args struct {
		ov math.Vector
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float32
	}{
		{ "1", fields{3,4}, args{math.Vector{4,3}}, 0.96},
	 { "2", fields{3,2}, args{math.Vector{18,6}}, 0.9647},
	 { "2", fields{-8,2}, args{math.Vector{18,6}}, -	0.8436},
	 { "2", fields{18,6}, args{math.Vector{18,6}.PerpendicularClockwise()}, 1},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := math.Vector{
				X: tt.fields.X,
				Y: tt.fields.Y,
			}
			if got := v.AngleTo(tt.args.ov); math.AlmostEqual(got, tt.want){
				t.Errorf("AngleTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
