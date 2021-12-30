package modules

import ( 
	"fmt"
	"time"

	"gmisail.me/gbar/style"
)

type Time struct { }

var (
	TIME_ICON string = style.Color("-", "#D5A021", "")
)

func (t Time) Run() string {
	return fmt.Sprintf("%s %s", TIME_ICON, time.Now().Format("3:04:05 PM"))
}
