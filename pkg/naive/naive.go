package naive

import (
	"bufio"
	"gobrc/pkg/data"
	"os"
	"strings"
)

const (
	zero_char_value = int16('0')
)

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

func ParseFileInto(path string, results map[string]data.StationData) uint64 {

	inputFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)

	lineCount := uint64(0)

	for inputScanner.Scan() {
		line := inputScanner.Text()

		name, temp := ParseLine(line)

		currData, found := results[name]
		if !found {
			currData.Min = temp
			currData.Max = temp
		} else {
			currData.Min = min(currData.Min, temp)
			currData.Max = max(currData.Max, temp)
		}

		currData.Count += 1
		currData.Sum += int64(temp)
		results[name] = currData

		lineCount += 1
	}

	return lineCount
}
