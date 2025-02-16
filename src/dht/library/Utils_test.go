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
