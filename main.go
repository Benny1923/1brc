package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"slices"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"1brc/mmap"
	"1brc/utils"
)

const TARGET = "measurements.txt"

var WORKERS = runtime.NumCPU()

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

func (s *Station) Merge(o *Station) {
	if !strings.EqualFold(s.name, o.name) {
		panic(fmt.Errorf("%v not equal to %v", o.name, s.name))
	}

	if o.min < s.min {
		s.min = o.min
	}

	if o.max > s.max {
		s.max = o.max
	}

	s.sum += o.sum
	s.count += o.count
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

func main() {
	w, _ := os.Create("main.pprof")
	pprof.StartCPUProfile(w)
	start := time.Now()
	mm, _ := mmap.NewMmap(TARGET)

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

	ch := make(chan Stations, WORKERS)
	fRange := int(mm.Size()) / WORKERS
	for i := 0; i < WORKERS; i++ {
		go worker(mm, i*fRange, (i+1)*fRange, ch)
	}

	stations := Stations{}
	for i := 0; i < WORKERS; i++ {
		ss := <-ch
		for o, os := range ss {
			s, ok := stations[o]
			if ok {
				s.Merge(os)
			} else {
				stations[o] = os
			}
		}
	}

	keys := stations.Keys()
	slices.SortFunc(keys, func(a, b string) int {
		return strings.Compare(utils.RemoveDiacritics(a), utils.RemoveDiacritics(b))
	})

	var sb strings.Builder
	sb.WriteByte('{')
	for _, key := range keys {
		s := stations[utils.StringHash(key)]
		sb.WriteString(s.name)
		sb.WriteByte('=')
		sb.Write(s.report())
		if key != keys[len(keys)-1] {
			sb.Write([]byte(", "))
		}
	}
	sb.WriteByte('}')

	fmt.Printf("%s\n", sb.String())
	fmt.Println(time.Since(start).String())
	pprof.StopCPUProfile()
}

type Stations map[uint64]*Station

func (m Stations) Keys() []string {
	result := make([]string, 0, len(m))
	for _, s := range m {
		result = append(result, s.name)
	}
	return result
}

func worker(mm *mmap.Mmap, begin, end int, ch chan<- Stations) {
	var buf []byte
	var next int
	var err error
	stations := Stations{}
	if begin != 0 {
		// drop first line
		_, next, _ = mm.ReadLineIdx(begin)
		next += 1
	}

	for next < end {
		buf, next, err = mm.ReadLineIdx(next)
		if err != nil {
			break
		}
		next += 1
		name, tempStr := utils.Split(buf)
		temp := utils.ParseFloat(tempStr)
		key := utils.BytesHash(name)
		s, ok := stations[key]
		if !ok {
			s = &Station{
				name: strings.Clone(unsafe.String(unsafe.SliceData(name), len(name))),
				min:  temp,
				max:  temp,
			}
			stations[key] = s
		}
		s.Add(temp)
	}

	ch <- stations
}
