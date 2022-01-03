package modules

import ( 
	"strconv"
	"time"

	//"gmisail.me/gbar/style"
)

type Time struct { }

func (t Time) Run() map[string] interface{} {
	currentTime := time.Now()

	return map[string] interface{} {
		"time-full": time.Now().Format("3:04:05 PM"),
		"time-day": strconv.Itoa(currentTime.Day()),
		"time-hour": strconv.Itoa(currentTime.Hour()),
		"time-minute": strconv.Itoa(currentTime.Minute()),
		"time-second": strconv.Itoa(currentTime.Second()),
	}
}
