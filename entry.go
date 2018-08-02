package artemis

import "fmt"

type Entry struct {
	Level Level
	Data  string
}

func defaultWrapper(l Level, content ...interface{}) Entry {
	return Entry{
		Level: l,
		Data:  fmt.Sprint(content...),
	}
}

func FatalEntry(contents ...interface{}) Entry {
	return defaultWrapper(Fatal, contents...)
}

func InfoEntry(contents ...interface{}) Entry {
	return defaultWrapper(Info, contents...)
}

func DebugEntry(contents ...interface{}) Entry {
	return defaultWrapper(Debug, contents...)
}

func TraceEntry(contents ...interface{}) Entry {
	return defaultWrapper(Trace, contents...)
}
