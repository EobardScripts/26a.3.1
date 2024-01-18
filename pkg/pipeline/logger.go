package pipeline

import (
	"fmt"
	"os"
)

// Debug - выводить ли отладочные сообщения.
var Debug = true

// DebugWriter - дескриптор для вывода отладочных сообщений.
var DebugWriter = os.Stderr

// logf - выводит отформатированное отладочное сообщение в DebugWriter.
func logf(format string, a ...interface{}) {
	if Debug {
		_, _ = fmt.Fprintf(DebugWriter, format+"\n", a)
	}
}
