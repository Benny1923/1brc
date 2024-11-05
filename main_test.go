package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func BenchmarkScan(b *testing.B) {
	size := []int{1024}
	for _, s := range size {
		b.Run(genScanFunc(s))
	}
}

func genScanFunc(size int) (string, func(*testing.B)) {
	return fmt.Sprint(size), func(b *testing.B) {
		file, _ := os.Open(TARGET)
		reader := bufio.NewReaderSize(file, size*1024)
		for {
			_, _, err := reader.ReadLine()
			if err != nil {
				return
			}
		}
	}
}
