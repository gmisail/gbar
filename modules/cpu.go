package modules

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
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

	temperatures, _ := host.SensorsTemperatures()

	cpuTemperature := 0.0                                                                                                                             
 	for i := 0; i < len(temperatures); i++ {                                                                                                          
    	if temperatures[i].SensorKey == "k10temp_tdie" {                                                                                              
 			cpuTemperature = temperatures[i].Temperature                                                                                              
			break                                                                                                                                   
		}                                                                                                                                             
 	} 

	return fmt.Sprintf(" %.2f%%  %2.f C", c.Usage, cpuTemperature)
}
