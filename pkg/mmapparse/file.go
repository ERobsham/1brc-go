package mmapparse

import (
	"errors"
	"os"
	"runtime"
	"sync"
	"syscall"
)

// provides _unsafe_ access to a memory mapped file

type MappedFile struct {
	isClosed  bool
	closeSync sync.Once

	data []byte
}

func OpenFileAt(path string) (*MappedFile, error) {
	const (
		OFFSET    = 0
		PROT_FLAG = syscall.PROT_READ
		MAP_FLAG  = syscall.MAP_SHARED
	)

	inputFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	fileInfo, err := inputFile.Stat()
	if err != nil {
		return nil, err
	}

	size := fileInfo.Size()
	if size < 0 {
		return nil, errors.New("invalid file size -- cannot be a negative number")
	}

	if potentiallyTruncatedSize := int64(int(size)); size != potentiallyTruncatedSize {
		return nil, errors.New("invalid file size -- too large, syscall will truncate data mapping")
	}

	fd := int(inputFile.Fd())
	if fd == 0 {
		return nil, errors.New("invalid file descriptor -- bailing early")
	}

	data, err := syscall.Mmap(fd, OFFSET, int(size), PROT_FLAG, MAP_FLAG)
	if err != nil {
		return nil, err
	}

	r := &MappedFile{
		data: data,
	}
	runtime.SetFinalizer(r, (*MappedFile).Close)
	return r, nil
}

func (f *MappedFile) Close() error {
	if f.data == nil {
		return nil
	}
	if len(f.data) == 0 {
		f.data = nil
		return nil
	}

	var err error
	if !f.isClosed {
		f.closeSync.Do(func() {
			temp := f.data
			f.data = nil

			runtime.SetFinalizer(f, nil)
			err = syscall.Munmap(temp)
		})
	}
	return err
}

func (f *MappedFile) Len() int {
	return len(f.data)
}

func (f *MappedFile) Read(offset int) byte {
	if offset > f.Len() {
		return 0
	}
	return f.data[offset]
}

func (f *MappedFile) ReadChunk(fromOffset int, length int) []byte {
	if fromOffset > f.Len() {
		return make([]byte, 0)
	}

	if fromOffset+length > f.Len() {
		return f.data[fromOffset:]
	}

	return f.data[fromOffset : fromOffset+length]
}
