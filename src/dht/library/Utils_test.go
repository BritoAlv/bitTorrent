package library

import (
	"testing"
)

func TestBetween(t *testing.T) {
	type args struct {
		L ChordHash
		M ChordHash
		R ChordHash
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases .
		{
			name: "Test 1",
			args: args{
				L: 252,
				M: 92,
				R: 242,
			},
			want: true,
		},
		{
			name: "Test 2",
			args: args{
				L: 10,
				M: 20,
				R: 15,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Between(tt.args.L, tt.args.M, tt.args.R); got != tt.want {
				t.Errorf("Between() = %v, want %v", got, tt.want)
			}
		})
	}
}