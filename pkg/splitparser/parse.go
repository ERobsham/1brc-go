package splitparser

import (
	"gobrc/pkg/data"
	"os"
	"runtime"
	"sync"
)

func ParseFileInto(path string, results map[string]data.StationData) uint64 {
	lineCount := uint64(0)

	inputFile, err := os.Open(path)
	info, _ := os.Stat(path)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	fileLen := int(info.Size())
	fileOffset := 0

	numChunks := runtime.NumCPU()
	chunkSize := fileLen / numChunks

	parsers := []*ChunkParser{}
	wg := sync.WaitGroup{}
	wg.Add(numChunks)

	var buffer [128]byte
	for i := 0; i < numChunks; i++ {
		startOffset := fileOffset
		length := (chunkSize - 1)

		if startOffset+length > fileLen {
			length = fileLen - (startOffset + 1)
		}

		bytesRead, _ := inputFile.ReadAt(buffer[:], int64(startOffset+length))
		for i := 0; i < bytesRead && buffer[i] != '\n'; i++ {
			length += 1
		}
		fileOffset += length + 1

		chunkLen := fileOffset - startOffset
		parser := ChunkParser{
			file:        inputFile,
			startOffset: startOffset,
			length:      chunkLen,

			results: map[string]data.StationData{},
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
	file        *os.File
	startOffset int
	length      int

	lineCount uint64
	results   map[string]data.StationData
}

func (p *ChunkParser) parse() {
	const BUFF_SIZE = 1 << 24
	fileOffset := p.startOffset
	endOffset := p.startOffset + p.length

	var buf [BUFF_SIZE]byte
	for fileOffset < endOffset {

		bytesRead, _ := p.file.ReadAt(buf[:], int64(fileOffset))
		offset := 0

	outer:
		for offset+data.MIN_LineLen < bytesRead && fileOffset < endOffset {
			lineStart := offset
			lineEnd := offset + data.MIN_LineLen
			for buf[lineEnd] != '\n' {
				lineEnd += 1

				if lineEnd >= bytesRead {
					break outer
				}
			}
			offset = lineEnd + 1
			fileOffset += offset - lineStart

			name, temp := data.ParseLine(buf[lineStart:lineEnd])

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
}
