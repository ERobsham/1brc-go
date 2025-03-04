package mmapparse

import (
	"gobrc/pkg/data"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseFileInto(t *testing.T) {
	const test_data = "WS1;10.0\nWS2;20.0\nWS3;-10.0\nWS4;-20.0\nWS1;-10.0\nWS2;-20.0\nWS3;10.0\nWS4;20.0\nWS1;10.0\nWS2;20.0\nWS3;-10.0\nWS4;-20.0\nWS1;0.0\nWS2;0.0\nWS3;0.0\nWS4;0.0\n"
	var test_path = filepath.Join(os.TempDir(), "test-1brc-data.txt")
	os.WriteFile(test_path, []byte(test_data), os.ModePerm)
	defer os.Remove(test_path)

	var wantCount uint64 = 16
	var wantResults map[string]data.StationData = map[string]data.StationData{
		"WS1": {Min: -100, Max: 100, Sum: 100, Count: 4},
		"WS2": {Min: -200, Max: 200, Sum: 200, Count: 4},
		"WS3": {Min: -100, Max: 100, Sum: -100, Count: 4},
		"WS4": {Min: -200, Max: 200, Sum: -200, Count: 4},
	}

	var gotCount uint64
	var gotResults = map[string]data.StationData{}

	gotCount = ParseFileInto(test_path, gotResults)
	if gotCount != wantCount {
		t.Errorf("ParseFileInto() = %v, want %v", gotCount, wantCount)
	}
	if !reflect.DeepEqual(gotResults, wantResults) {
		t.Errorf("ParseFileInto() = %v, want %v", gotResults, wantResults)
	}
}
