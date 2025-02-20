package library

import (
	"reflect"
	"testing"
)

func TestIntToBinaryArray(t *testing.T) {
	type args struct {
		number int
	}
	tests := []struct {
		name string
		args args
		want [5]uint8
	}{
		{
			name: "Test 1",
			args: args{number: 0},
			want: [5]uint8{0, 0, 0, 0, 0},
		},
		{
			name: "Test 2",
			args: args{number: 1},
			want: [5]uint8{1, 0, 0, 0, 0},
		},

		{
			name: "Test 3",
			args: args{number: 31},
			want: [5]uint8{1, 1, 1, 1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntToBinaryArray(tt.args.number); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntToBinaryArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
				L: IntToBinaryArray(252),
				M: IntToBinaryArray(92),
				R: IntToBinaryArray(242),
			},
			want: true,
		},
		{
			name: "Test 2",
			args: args{
				L: IntToBinaryArray(10),
				M: IntToBinaryArray(20),
				R: IntToBinaryArray(15),
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

func TestBinaryArrayToInt(t *testing.T) {
	type args struct {
		array ChordHash
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test 1",
			args: args{array: ChordHash{0, 0, 0, 0, 0, 0, 0, 0}},
			want: 0,
		},
		{
			name: "Test 2",
			args: args{array: ChordHash{1, 0, 0, 0, 0, 0, 0, 0}},
			want: 1,
		},
		{
			name: "Test 3",
			args: args{array: ChordHash{1, 1, 1, 1, 1, 1, 1, 1}},
			want: 255,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BinaryArrayToInt(tt.args.array); got != tt.want {
				t.Errorf("BinaryArrayToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
