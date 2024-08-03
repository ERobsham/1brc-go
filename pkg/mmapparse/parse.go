package mmapparse

import (
	"gobrc/pkg/data"
	"runtime"
	"sync"
)

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
		data.Merge(results, parser.results)
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
		lineEnd := offset + data.MIN_LineLen

		for p.data[lineEnd] != '\n' {
			lineEnd += 1
		}
		offset = lineEnd + 1

		name, temp := data.ParseLine(p.data[lineStart:lineEnd])

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
