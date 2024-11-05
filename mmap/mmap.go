package mmap

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"syscall"
)

type Mmap struct {
	data []byte
	idx  int64
}

func (r *Mmap) ReadLine() ([]byte, error) {
	buf, end, err := r.ReadLineIdx(int(r.idx))
	r.idx = int64(end + 1)
	return buf, err
}

func (r *Mmap) ReadLineIdx(start int) (buf []byte, end int, err error) {
	size := len(r.data)
	if start >= size {
		return nil, 0, io.EOF
	}
	n := bytes.IndexByte(r.data[start:], '\n')
	if n >= 0 {
		end = start + n
	} else {
		end = size
	}
	buf = r.data[start:end]
	return
}

func (r *Mmap) Size() int64 {
	return int64(len(r.data))
}

func (r *Mmap) Close() error {
	data := r.data
	runtime.SetFinalizer(r, nil)
	return syscall.Munmap(data)
}

func NewMmap(path string) (*Mmap, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := stat.Size()

	data, err := syscall.Mmap(int(file.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	r := &Mmap{data, 0}
	runtime.SetFinalizer(r, (*Mmap).Close)
	return r, nil
}
