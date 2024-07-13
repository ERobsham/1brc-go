package naive

import (
	"fmt"
	"strings"
)

const (
	zero_char_value = int16('0')
)

type StationData struct {
	Min   int16
	Max   int16
	Count int32
	Sum   int64
}

func (d StationData) String() string {
	return fmt.Sprintf("%0.1f/%0.1f/%0.1f", (float64(d.Min) / 10.0), (float64(d.Sum) / float64(d.Count) / 10.0), float64(d.Max)/10.0)
}

func ParseLine(line string) (string, int16) {
	pre, post, found := strings.Cut(line, ";")
	if !found || pre == "" || post == "" {
		panic("invalid ParseLine() data")
	}

	var temp int16 = 0
	var isPositive int16 = 1
	for _, v := range post {
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

	return pre, temp
}
