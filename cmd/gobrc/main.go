package main

import (
	"bufio"
	"fmt"
	"gobrc/pkg/naive"
	"os"
	"path/filepath"
	"runtime/pprof"
	"slices"
)

const (
	data_dir = "data"
)

var (
	FLAG_DebugLogs = ""
	FLAG_CPUProf   = ""
	FLAG_Output    = ""
	FLAG_DataFile  = "measurements-100m.txt"
)

var (
	ENABLE_DebugLogs = (FLAG_DebugLogs != "")
	ENABLE_CPUProf   = (FLAG_CPUProf != "")
	ENABLE_Output    = (FLAG_Output != "")
)

func main() {
	if ENABLE_CPUProf {
		pOutFile, err := os.Create("CPUProf.out")
		if err != nil {
			panic(err)
		}
		defer pOutFile.Close()

		pprof.StartCPUProfile(pOutFile)
		defer pprof.StopCPUProfile()
	}

	dir, _ := os.Getwd()
	path := filepath.Join(dir, data_dir, FLAG_DataFile)

	if ENABLE_DebugLogs {
		fmt.Println("parsing file: ", path)
	}

	inputFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)

	results := map[string]naive.StationData{}
	lineCount := uint64(0)

	for inputScanner.Scan() {
		line := inputScanner.Text()

		name, temp := naive.ParseLine(line)

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

	keys := make([]string, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	if ENABLE_DebugLogs {
		fmt.Printf("read %d lines\n", lineCount)
	}

	if ENABLE_Output {
		writer := bufio.NewWriter(os.Stdout)
		isFirst := true
		separator := []byte(", ")

		writer.WriteRune('{')
		for _, k := range keys {
			v := results[k]
			if isFirst {
				isFirst = false
			} else {
				writer.Write(separator)
			}
			writer.WriteString(k)
			writer.WriteRune('=')
			writer.WriteString(v.String())
		}
		writer.WriteRune('}')

		writer.Flush()
	}
}
