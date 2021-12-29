package modules

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
)

type CPU struct {
	Usage float64
}

/*
 *	Returns the CPU usage as a percentage
 */
func (c CPU) Run() string {
	times, _ := cpu.Percent(0, false)
	c.Usage = times[0]

	return fmt.Sprintf("%.2f%%", c.Usage)
}
