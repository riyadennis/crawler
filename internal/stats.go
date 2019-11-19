package internal

import (
	"fmt"
	"runtime"
)

type MemStats struct {
	Alloc      uint64
	TotalAlloc uint64
}

func GetMemStats() *MemStats {
	var mem runtime.MemStats
	fmt.Println("memory baseline ...")
	runtime.ReadMemStats(&mem)
	return &MemStats{
		Alloc:      mem.Alloc,
		TotalAlloc: mem.TotalAlloc,
	}
}
