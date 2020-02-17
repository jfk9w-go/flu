package flu

import (
	"os"
	"os/signal"
)

func AwaitSignal(signals ...os.Signal) {
	c := make(chan os.Signal)
	signal.Notify(c, signals...)
	<-c
}
