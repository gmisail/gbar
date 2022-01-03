package modules

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
)

type CPU struct {}

/*
 *	Store icons and their colors so you don't need to regenerate every time
 */
/*var (
	CPU_ICON = style.Color("-", "#D6E3F8", "")
)*/

/*
 *	Returns the CPU usage as a percentage
 */
func (c CPU) Run() map[string] interface{} {
	times, _ := cpu.Percent(0, false)
	temperatures, _ := host.SensorsTemperatures()

	cpuTemperature := 0.0
 	for i := 0; i < len(temperatures); i++ {
    	if temperatures[i].SensorKey == "k10temp_tdie" {
 			cpuTemperature = temperatures[i].Temperature
			break
		}
 	}

	return map[string] interface{} {
		"cpu-usage": fmt.Sprintf("%.2f", times[0]),
		"cpu-temperature": fmt.Sprintf("%.2f", cpuTemperature),
	}
}
