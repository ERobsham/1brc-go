package mmapparse

import (
	"gobrc/pkg/data"
)

const CHUNK_SIZE = 1 << 16

func ParseFileInto(path string, results map[string]data.StationData) uint64 {
	lineCount := uint64(0)

	inputFile, err := OpenFileAt(path)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	fileLen := inputFile.Len()
	fileOffset := 0
	data := inputFile.ReadChunk(fileOffset, fileLen)

	const minLineLen = len("a;0.0")

	for fileOffset < fileLen {

		lineStart := fileOffset
		lineEnd := fileOffset + minLineLen

		for data[lineEnd] != '\n' {
			lineEnd += 1
		}
		fileOffset = lineEnd + 1

		name, temp := parseLine(data[lineStart:lineEnd])

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

func offsetFromEndTo(buf []byte, b byte) int {
	endOffset := len(buf) - 1
	for i := endOffset; i >= 0; i-- {
		if buf[i] == b {
			return endOffset - i
		}
	}
	return -1
}

func parseLine(line []byte) (string, int16) {
	const zero_char_value = int16('0')

	lineLen := len(line)
	offset := offsetFromEndTo(line, ';')

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
