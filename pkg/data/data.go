package data

import "fmt"

const MIN_LineLen = len("a;0.0")

type StationData struct {
	Min   int16
	Max   int16
	Count int32
	Sum   int64
}

func (d StationData) String() string {
	return fmt.Sprintf("%0.1f/%0.1f/%0.1f", (float64(d.Min) / 10.0), (float64(d.Sum) / float64(d.Count) / 10.0), float64(d.Max)/10.0)
}

//

func OffsetFromEndTo(buf []byte, b byte) int {
	endOffset := len(buf) - 1
	for i := endOffset; i >= 0; i-- {
		if buf[i] == b {
			return endOffset - i
		}
	}
	return -1
}

func ParseLine(line []byte) (string, int16) {
	const zero_char_value = int16('0')

	lineLen := len(line)
	offset := OffsetFromEndTo(line, ';')

	var temp int16 = 0
	var isPositive int16 = 1
	for _, v := range line[lineLen-offset:] {
		switch v {
		case '-':
			isPositive = -1
		case '.':
			continue
		default:
			var val int16 = int16(v) - zero_char_value

			temp = (temp * 10) + val
		}
	}

	temp = isPositive * temp

	return string(line[:lineLen-(offset+1)]), temp
}

func Merge(dest map[string]StationData, results map[string]StationData) {
	for k, result := range results {
		orig, exists := dest[k]
		if !exists {
			dest[k] = result
			continue
		}

		orig.Count += result.Count
		orig.Sum += result.Sum
		orig.Min = min(orig.Min, result.Min)
		orig.Max = max(orig.Max, result.Max)

		dest[k] = orig
	}
}
