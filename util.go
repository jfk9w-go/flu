package flu

import (
	"os"
	"os/signal"
	"syscall"
)

func AwaitSignal(signals ...os.Signal) {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM}
	}

	c := make(chan os.Signal)
	signal.Notify(c, signals...)
	<-c
}
