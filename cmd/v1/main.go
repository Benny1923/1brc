package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime/pprof"
	"slices"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"1brc/utils"
)

const TARGET = "measurements.txt"

type Station struct {
	name  string
	min   float64
	max   float64
	sum   float64
	count int
}

func (s *Station) Add(v float64) {

	if v < s.min {
		s.min = v
	}

	if v > s.max {
		s.max = v
	}

	s.sum += v
	s.count += 1
}

func (s *Station) report() []byte {
	mean := s.sum / float64(s.count)

	buf := []byte{}
	buf = strconv.AppendFloat(buf, s.min, 'f', 1, 64)
	buf = append(buf, '/')
	buf = strconv.AppendFloat(buf, mean, 'f', 1, 64)
	buf = append(buf, '/')
	buf = strconv.AppendFloat(buf, s.max, 'f', 1, 64)
	return buf
}

// not that fast implementation
// may have hash collision
func main() {
	w, _ := os.Create("v1.pprof")
	pprof.StartCPUProfile(w)
	start := time.Now()
	file, _ := os.Open(TARGET)

	reader := bufio.NewReaderSize(file, 1024*1024)

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
		pprof.StopCPUProfile()
	}()

	go func() {
		time.Sleep(3 * time.Minute)
		panic("time out")
	}()

	mod := uint64(500000)
	ss := make([]*Station, mod)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		name, tempStr := split(line)
		temp := parseFloat(tempStr)
		idName := utils.BytesHash(name)

		s := ss[idName%mod]
		if s == nil {
			s = &Station{
				name: strings.Clone(unsafe.String(unsafe.SliceData(name), len(name))),
				min:  temp,
				max:  temp,
			}
			ss[idName%mod] = s
		}

		s.Add(temp)
	}

	keys := []string{}
	for _, v := range ss {
		if v != nil {
			keys = append(keys, v.name)
		}
	}
	slices.Sort(keys)

	var sb strings.Builder
	sb.WriteByte('{')
	for _, key := range keys {
		s := ss[utils.StringHash(key)%mod]
		sb.WriteString(s.name)
		sb.WriteByte('=')
		sb.Write(s.report())
		if key != keys[len(keys)-1] {
			sb.Write([]byte(", "))
		}
	}
	sb.WriteByte('}')

	fmt.Printf("%s\n", sb.String())
	elapsed := time.Since(start)
	fmt.Println(elapsed.String())
	pprof.StopCPUProfile()
}

func split(text []byte) (name, temp []byte) {
	idx := bytes.IndexByte(text, ';')
	name = text[0:idx]
	temp = text[idx+1:]
	return
}

func parseFloat(text []byte) float64 {
	var result float64
	infloat := 0.0
	minus := false
	for _, b := range text {
		if b == '.' {
			infloat = 0.1
			continue
		} else if b == '-' {
			minus = true
			continue
		}

		if infloat != 0.0 {
			result = result + float64(b-48)*infloat
			infloat *= 0.1
		} else {
			result = result*10 + float64(b-48)
		}
	}

	if minus {
		return -result
	}

	return result
}
