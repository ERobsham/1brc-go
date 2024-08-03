package genscript

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

func TestReadStations(t *testing.T) {
	stationsFile, _ := os.Open("./weather_stations.csv")

	scanner := bufio.NewScanner(stationsFile)

	names := map[string]struct{}{}

	for scanner.Scan() {
		line := scanner.Text()
		components := strings.Split(line, ";")
		names[components[0]] = struct{}{}
	}

	fmt.Println("Name count:", len(names))

	datas := map[[4]uint64]struct{}{}

	for name := range names {
		data := convertToFixedData(name)

		datas[data] = struct{}{}
	}

	fmt.Println("NameData count:", len(datas))

	if len(datas) != len(names) {
		t.Fail()
	}
}

func convertToFixedData(name string) [4]uint64 {
	const (
		// OFF_4 = 8 * 4 // 0b0010_1000 == 40
		OFF_3 = 8 * 4 // 0b0010_0000 == 32
		OFF_2 = 8 * 3 // 0b0001_1000 == 24
		OFF_1 = 8 * 2 // 0b0001_0000 == 16
		OFF_0 = 8 * 1 // 0b0000_1000 == 8
	)

	var data [4]uint64

	nLen := uint8(len(name))
	// nLen4 := min(nLen, OFF_4)
	nLen3 := min(nLen, OFF_3)
	nLen2 := min(nLen, OFF_2)
	nLen1 := min(nLen, OFF_1)
	nLen0 := min(nLen, OFF_0)

	// idx4 := ((nLen4 >> 5) & 0b1) >> 2 & ((nLen4 >> 3) & 0b001)
	idx3 := ((nLen3 >> 5) & 0b001)
	idx2 := (((nLen2 >> 3) & 0b010) >> 1 & ((nLen2 >> 3) & 0b001)) | idx3
	idx1 := ((nLen1>>3)&0b010)>>1 | idx2
	idx0 := ((nLen0 >> 3) & 0b001) | idx1

	strData := unsafe.StringData(name)
	ptr := unsafe.Pointer(strData)

	data[0] = (*(*uint64)(unsafe.Pointer(uintptr(ptr) + uintptr(0)))) & (1<<(nLen0*8) - 1)
	data[1] = (*(*uint64)(unsafe.Pointer(uintptr(ptr) + uintptr(idx0*OFF_0)))) & (1<<((nLen1-nLen0)*8) - 1)
	data[2] = (*(*uint64)(unsafe.Pointer(uintptr(ptr) + uintptr(idx1*OFF_1)))) & (1<<((nLen2-nLen1)*8) - 1)
	data[3] = (*(*uint64)(unsafe.Pointer(uintptr(ptr) + uintptr(idx2*OFF_2)))) & (1<<((nLen3-nLen2)*8) - 1)

	return data
}

func convertToUint64(str string) uint64 {
	sLen := min(uint8(len(str)), 8)
	strData := unsafe.StringData(str)
	strAsUintPtr := (*uint64)(unsafe.Pointer(strData))
	return *strAsUintPtr & (1<<(sLen*8) - 1)
}

func Test_convertToFixedData(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want [4]uint64
	}{
		{
			name: "partial chunk",
			args: args{name: "123"},
			want: [4]uint64{
				convertToUint64("123"),
			},
		},
		{
			name: "one and a partial chunk",
			args: args{name: "12345678123"},
			want: [4]uint64{
				convertToUint64("12345678"),
				convertToUint64("123"),
			},
		},
		{
			name: "over size",
			args: args{name: "1234567812345678123456781234567812345678"},
			want: [4]uint64{
				convertToUint64("12345678"),
				convertToUint64("12345678"),
				convertToUint64("12345678"),
				convertToUint64("12345678"),
			},
		},
		{
			name: "four chunks",
			args: args{name: "12345678123456781234567812345678"},
			want: [4]uint64{
				convertToUint64("12345678"),
				convertToUint64("12345678"),
				convertToUint64("12345678"),
				convertToUint64("12345678"),
			},
		},
		{
			name: "three chunks",
			args: args{name: "123456781234567812345678"},
			want: [4]uint64{
				convertToUint64("12345678"),
				convertToUint64("12345678"),
				convertToUint64("12345678"),
			},
		},
		{
			name: "two chunks",
			args: args{name: "1234567812345678"},
			want: [4]uint64{
				convertToUint64("12345678"),
				convertToUint64("12345678"),
			},
		},
		{
			name: "one chunk",
			args: args{name: "12345678"},
			want: [4]uint64{
				convertToUint64("12345678"),
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertToFixedData(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToFixedData() = %v, want %v", got, tt.want)
			}
		})
	}
}
