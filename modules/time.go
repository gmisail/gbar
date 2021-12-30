package modules

import ( 
	"fmt"
	"time"
)

type Time struct { }

func (t Time) Run() string {
	return fmt.Sprintf(time.Now().Format("3:04:05 PM"))
}
