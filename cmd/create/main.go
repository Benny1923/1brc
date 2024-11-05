package main

import (
	"1brc/utils"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var target = "measurements.txt"

var lines = 1_000_000_000

func main() {
	start := time.Now()
	file, _ := os.Create(target)
	defer file.Close()
	buffer := bufio.NewWriterSize(file, 1024*1024)
	defer buffer.Flush()
	for i := 0; i < lines; i++ {
		if i%50_000_000 == 0 && i != 0 {
			fmt.Printf("Wrote %d measurements in %s \n", i, time.Since(start))
		}
		ws := utils.WeatherStations[rand.Intn(len(utils.WeatherStations))]
		buffer.WriteString(ws.Name)
		buffer.Write([]byte{';'})
		buffer.WriteString(strconv.FormatFloat(ws.Measurement(), 'f', 1, 64))
		buffer.Write([]byte{'\n'})
	}
	fmt.Printf("Wrote %d measurements in %s \n", lines, time.Since(start))
}
