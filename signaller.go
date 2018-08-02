package artemis

import (
	"os"
	"syscall"
)

// signalProcess is only used once a fatal log message has been called.
// The reason being is so that the process itself can have a chance to recover from
// this event or finish processing some shit
func signalProcess() error {
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	return proc.Signal(syscall.SIGABRT)
}
