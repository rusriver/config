package config

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/dustin/go-humanize"
)

func (c *Config) PrintJson(tag string) {
	bb, err := json.MarshalIndent(c, "", "    ")
	fmt.Printf("PRINT JSON, tag='%v' err='%v'        -- %v\n%s\n", tag, err, tag, bb)
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("MEM USAGE: Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tHeapObjects = %v", humanize.SIWithDigits(float64(m.HeapObjects), 2, ""))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
