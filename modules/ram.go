package modules

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
)

type RAM struct {
	Usage float64
}

/*
 *	Returns the CPU usage as a percentage
 */
func (r RAM) Run() string {
	memory, _ := mem.VirtualMemory()
	return fmt.Sprintf("  %.2f%%", memory.UsedPercent)
}
