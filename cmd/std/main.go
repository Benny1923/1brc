package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"time"
)

const TARGET = "measurements.txt"

type Station struct {
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

type Stations map[string]*Station

func (s Stations) Keys() []string {
	result := make([]string, 0, len(s))

	for key := range s {
		result = append(result, key)
	}

	return result
}

func (s Stations) String() string {
	var builder strings.Builder

	builder.WriteByte('{')

	keys := s.Keys()
	sort.Strings(keys)
	last := len(keys) - 1
	for idx, key := range keys {
		builder.WriteString(key)
		builder.WriteByte('=')
		builder.Write(s[key].report())
		if idx != last {
			builder.Write([]byte(", "))
		}
	}
	builder.WriteByte('}')

	return builder.String()
}

func main() {
	profile, _ := os.Create("std.pprof")
	pprof.StartCPUProfile(profile)
	defer pprof.StopCPUProfile()
	start := time.Now()

	stations := Stations{}
	file, _ := os.Open(TARGET)

	scanner := bufio.NewScanner(file)
	bufSize := 1024 * 1024
	scanner.Buffer(make([]byte, bufSize), bufSize)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ";")
		name := parts[0]
		temp, _ := strconv.ParseFloat(parts[1], 64)
		s, ok := stations[name]
		if !ok {
			s = &Station{
				min: temp,
				max: temp,
			}
			stations[name] = s
		}
		s.Add(temp)
	}

	fmt.Printf("%s\n", stations)
	fmt.Printf("%s\n", time.Since(start))
}
