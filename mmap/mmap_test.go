package mmap

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const TARGET = "../measurements.txt"

func BenchmarkMmap(b *testing.B) {
	mm, err := NewMmap(TARGET)
	if err != nil {
		b.Errorf("mmap failed %v", err)
	}

	for {
		_, err := mm.ReadLine()
		if err != nil {
			break
		}
	}
}

func BenchmarkMmapMulti(b *testing.B) {
	mm, err := NewMmap(TARGET)
	if err != nil {
		b.Errorf("mmap failed %v", err)
	}

	c := 4
	var wg sync.WaitGroup
	ch := make(chan int, 4)

	worker := func(begin, end int) {
		defer wg.Done()
		var err error
		next := 0
		if begin != 0 {
			_, next, _ = mm.ReadLineIdx(begin)
		}
		s := 0
		for {
			_, next, err = mm.ReadLineIdx(next + 1)
			if err != nil {
				break
			}
			s += 1
			if next+1 > end {
				break
			}
		}
		ch <- s
	}

	wg.Add(c)
	for i := 0; i < c; i++ {
		r := mm.Size() / 4
		go worker(int(r)*i, int(r)*(i+1))
	}
	wg.Wait()

	sum := 0
	for i := 0; i < c; i++ {
		sum += <-ch
	}

	fmt.Println(sum)
}

func TestMmapReadLine(t *testing.T) {
	TARGET := "../main.go"
	mm, err := NewMmap(TARGET)
	if err != nil {
		t.Errorf("mmap failed %v", err)
	}
	file, _ := os.Open(TARGET)
	reader := bufio.NewReaderSize(file, 1024*1024)

	i := 0
	for {
		expect, _, err := reader.ReadLine()
		line, mErr := mm.ReadLine()
		if err != nil && mErr != nil {
			break
		} else if err != nil && mErr == nil {
			t.Errorf("mmap except err")
		} else if mErr != nil {
			t.Errorf("mmap failed %v", err)
		}
		require.EqualValues(t, expect, line, "in round %d", i)
		i++
	}
}
