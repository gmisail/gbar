package modules

import (
	"fmt"
	"gmisail.me/gbar/style"
//	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
)

type CPU struct {
	Usage float64
}

/* 
 *	Store icons and their colors so you don't need to regenerate every time
 */
var (
	CPU_ICON = style.Color("-", "#D6E3F8", "")
    CPU_TEMP_ICON = style.Color("-", "#FCFC62", "") 
)

/*
 *	Returns the CPU usage as a percentage
 */
func (c CPU) Run() string {
//	times, _ := cpu.Percent(0, false)
//	c.Usage = times[0]

	temperatures, _ := host.SensorsTemperatures()

	cpuTemperature := 0.0                                                                                                                             
 	for i := 0; i < len(temperatures); i++ {                                                                                                          
    	if temperatures[i].SensorKey == "k10temp_tdie" {                                                                                              
 			cpuTemperature = temperatures[i].Temperature                                                                                              
			break                                                                                                                                   
		}                                                                                                                                             
 	} 

	//return fmt.Sprintf("%s %.2f%% %s %2.f C", CPU_ICON, c.Usage, CPU_TEMP_ICON, cpuTemperature)

	return fmt.Sprintf(" %s %2.f C", CPU_TEMP_ICON, cpuTemperature)
}
