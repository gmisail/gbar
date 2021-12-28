package modules

type CPU struct {
	Usage float64
}

/*
 *	Returns the CPU usage as a percentage
 */
func (cpu *CPU) Run() string {
	times, err := cpu.Percent(0, false)
	cpu.Usage = times[0]

	return fmt.Sprintf("%.2f%%", cpu.Usage)
}
