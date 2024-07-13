package naive

import "testing"

func TestParseLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 int16
	}{
		{
			name:  "max value",
			args:  args{line: "test;99.9"},
			want:  "test",
			want1: 999,
		},
		{
			name:  "min value",
			args:  args{line: "test;-99.9"},
			want:  "test",
			want1: -999,
		},
		{
			name:  "zero value",
			args:  args{line: "test;0.0"},
			want:  "test",
			want1: 0,
		},
		{
			name:  "1.0",
			args:  args{line: "test;1.0"},
			want:  "test",
			want1: 10,
		},
		{
			name:  "-1.1",
			args:  args{line: "test;-1.1"},
			want:  "test",
			want1: -11,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ParseLine(tt.args.line)
			if got != tt.want {
				t.Errorf("ParseLine() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseLine() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestParseLinePanicsWithBadValues(t *testing.T) {
	tests := []struct {
		name string
		line string
	}{
		{
			name: "empty line",
			line: "",
		},
		{
			name: "no temp",
			line: "test;",
		},
		{
			name: "only temp",
			line: ";-11.1",
		},
		{
			name: "only delimiter",
			line: ";",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Errorf("ParseLine() didn't panic for a bad value!")
				}
			}()
			ParseLine(tt.line)
		})
	}
}
