package flu

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Clock interface {
	Now() time.Time
}

type ClockFunc func() time.Time

func (fun ClockFunc) Now() time.Time {
	return fun()
}

var DefaultClock Clock = ClockFunc(time.Now)

func AwaitSignal(signals ...os.Signal) {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM}
	}

	c := make(chan os.Signal)
	signal.Notify(c, signals...)
	<-c
}
