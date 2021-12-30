package modules

import (
	"fmt"
	"gmisail.me/gbar/style"
	"github.com/shirou/gopsutil/v3/mem"
)

type RAM struct {
	Usage float64
}

var (
	RAM_ICON string = style.Color("-", "#CC3F0C", "")
)

/*
 *	Returns the CPU usage as a percentage
 */
func (r RAM) Run() string {
	memory, _ := mem.VirtualMemory()
	return fmt.Sprintf("%s %.2f%%", RAM_ICON, memory.UsedPercent)
}
