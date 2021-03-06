package modules

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
)

type RAM struct {}

/*
 *	Returns RAM statistics
 */
func (r RAM) Run() map[string] interface{} {
	memory, _ := mem.VirtualMemory()

	return map[string] interface{} {
		"mem-total": memory.Total,
		"mem-available": memory.Available,
		"mem-used": memory.Used,
		"mem-used-percentage": fmt.Sprintf("%.2f", memory.UsedPercent),
		"mem-free": memory.Free,
	}
}
