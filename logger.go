package artemis

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"sync"
)

type instance struct {
	entries chan Entry
	done    chan bool
	output  []io.Writer
	level   Level
	running bool
}

var (
	isnt *instance
	once sync.Once
	// BufferSize defines how many log entries can store before locking the application
	BufferSize int
)

func init() {
	BufferSize = 100
	GetInstance().Start()
}

// GetInstance will returns the logger singleton
// ready for use.
func GetInstance() Log {
	once.Do(func() {
		isnt = &instance{
			level: Fatal,
		}
		runtime.SetFinalizer(isnt, func(inst *instance) {
			inst.Stop()
		})
	})
	return isnt
}

func (i *instance) Log(e Entry) Log {
	if !i.running {
		return i
	}
	_, fn, line, _ := runtime.Caller(1)
	e.Data = fmt.Sprintf("[%s:%d]\t%s", path.Base(fn), line, e.Data)
	i.entries <- e
	if i.level == Debug {
		var m runtime.MemStats
		convert := func(val uint64) uint64 {
			return val / 1024 / 1024
		}
		runtime.ReadMemStats(&m)
		i.entries <- Entry{
			Level: Debug,
			Data:  fmt.Sprintf("[Current Usage] %v MiB, [GC Count] %v, [GoRoutines] %v", convert(m.Alloc), m.NumGC, runtime.NumGoroutine()),
		}
	}
	return i
}

func (i *instance) Set(level Level, writers ...io.Writer) Log {
	i.output = append(i.output, writers...)
	i.level = level
	return i
}

func (i *instance) Start() {
	if i.running {
		// Don't try start an already running logger
		return
	}
	i.running = true
	i.entries = make(chan Entry, BufferSize)
	i.done = make(chan bool)
	go func() {
		for data := range i.entries {
			if data.Level <= i.level {
				for _, out := range i.output {
					fmt.Fprintf(out, "[%s]\t%s\n", data.Level.String(), data.Data)
				}
				if data.Level == Fatal {
					// Issue a signal to the process so that it knows it should die
					if err := signalProcess(); err != nil {
						panic(err)
					}
				}
			}
		}
		i.done <- true
	}()
}

func (i *instance) Stop() {
	if !i.running {
		// Don't try to stop an non running logger
		return
	}
	i.Log(Entry{
		Level: Info,
		Data:  "Logger is being stopped",
	})
	i.running = false
	for {
		if len(i.entries) == 0 {
			break
		}
		// Wait for the buffer to empty
	}
	close(i.entries)
	<-i.done
	close(i.done)
}
