package modules

import (
	"fmt"
	//"gmisail.me/gbar/style"
	"github.com/shirou/gopsutil/v3/mem"
)

type RAM struct {}

/*
var (
	RAM_ICON string = style.Color("-", "#CC3F0C", "")
)
*/

/*
 *	Returns RAM statistics
 */
func (r RAM) Run() map[string] interface{} {
	memory, _ := mem.VirtualMemory()
//	return fmt.Sprintf("%s %.2f%%", RAM_ICON, memory.UsedPercent)

	return map[string] interface{} {
		"mem-total": memory.Total,
		"mem-available": memory.Available,
		"mem-used": memory.Used,
		"mem-used-percentage": fmt.Sprintf("%.2f", memory.UsedPercent),
		"mem-free": memory.Free,
	}
}
