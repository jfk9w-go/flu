package flu

import (
	"os"
	"os/signal"
)

func HandleSignals(signals ...os.Signal) {
	c := make(chan os.Signal)
	signal.Notify(c, signals...)
	<-c
}
