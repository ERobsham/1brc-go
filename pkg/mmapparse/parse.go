package mmapparse

import (
	"gobrc/pkg/data"
	"runtime"
	"sync"
)

const CHUNK_SIZE = 1 << 16

const MIN_LineLen = len("a;0.0")

func ParseFileInto(path string, results map[string]data.StationData) uint64 {
	lineCount := uint64(0)

	inputFile, err := OpenFileAt(path)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	fileLen := inputFile.Len()
	fileOffset := 0

	numChunks := runtime.NumCPU()
	chunkSize := fileLen / numChunks

	parsers := []*ChunkParser{}
	wg := sync.WaitGroup{}
	wg.Add(numChunks)

	for i := 0; i < numChunks; i++ {
		startOffset := fileOffset
		endOffset := startOffset + (chunkSize - 1)

		for inputFile.Read(endOffset) != '\n' && endOffset < fileLen {
			endOffset += 1
		}
		fileOffset = endOffset + 1

		chunkLen := fileOffset - startOffset
		parser := ChunkParser{
			results: map[string]data.StationData{},
			data:    inputFile.ReadChunk(startOffset, chunkLen),
		}
		parsers = append(parsers, &parser)

		go func(p *ChunkParser, g *sync.WaitGroup) {
			p.parse()
			wg.Done()
		}(&parser, &wg)
	}
	wg.Wait()

	for _, parser := range parsers {
		lineCount += parser.lineCount
		merge(results, parser.results)
	}

	return lineCount
}

type ChunkParser struct {
	lineCount uint64
	results   map[string]data.StationData
	data      []byte
}

func (p *ChunkParser) parse() {
	dataLen := len(p.data)
	offset := 0

	for offset < dataLen {

		lineStart := offset
		lineEnd := offset + MIN_LineLen

		for p.data[lineEnd] != '\n' {
			lineEnd += 1
		}
		offset = lineEnd + 1

		name, temp := parseLine(p.data[lineStart:lineEnd])

		currData, found := p.results[name]
		if !found {
			currData.Min = temp
			currData.Max = temp
		} else {
			currData.Min = min(currData.Min, temp)
			currData.Max = max(currData.Max, temp)
		}

		currData.Count += 1
		currData.Sum += int64(temp)
		p.results[name] = currData

		p.lineCount += 1
	}
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

func merge(dest map[string]data.StationData, results map[string]data.StationData) {
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
