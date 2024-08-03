package data

import "testing"

func Test_offsetFromEndTo(t *testing.T) {
	type args struct {
		buf []byte
		b   byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "newline at end",
			args: args{buf: []byte("987654321\n"), b: '\n'},
			want: 0,
		},
		{
			name: "newline at end-1",
			args: args{buf: []byte("98765432\n1"), b: '\n'},
			want: 1,
		},
		{
			name: "newline at start",
			args: args{buf: []byte("\n987654321"), b: '\n'},
			want: 9,
		},
		{
			name: "no newline",
			args: args{buf: []byte("987654321"), b: '\n'},
			want: -1,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OffsetFromEndTo(tt.args.buf, tt.args.b); got != tt.want {
				t.Errorf("offsetFromEndTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
