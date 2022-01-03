package modules

import ( 
//	"fmt"
	"time"

	//"gmisail.me/gbar/style"
)

type Time struct { }

/*
var (
	TIME_ICON string = style.Color("-", "#D5A021", "")
)*/

func (t Time) Run() map[string] interface{} {
	currentTime := time.Now()

	return map[string] interface{} {
		"time-full": time.Now().Format("3:04:05 PM"),
		"time-year": currentTime.Year(),
		"time-month": currentTime.Month(),
		"time-day": currentTime.Day(),
		"time-hour": currentTime.Hour(),
		"time-minute": currentTime.Minute(),
		"time-second": currentTime.Second(),
		"time-weekday": currentTime.Weekday(),
	}

//	return fmt.Sprintf("%s %s", TIME_ICON, time.Now().Format("3:04:05 PM"))
}
