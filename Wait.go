package nerdweb

import (
	"os"
	"os/signal"
	"syscall"
)

/*
WaitForKill returns a channel that waits for an OS interrupt or terminate.
*/
func WaitForKill() chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return quit
}
