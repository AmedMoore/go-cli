package cli

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	info *log.Logger
	warn *log.Logger
	err  *log.Logger
}

func NewLogger(out io.Writer, flags ...int) *Logger {
	l := new(Logger)

	flag := log.Ldate | log.Ltime | log.Llongfile
	if len(flags) > 0 {
		flag = flags[0]
	}

	l.info = log.New(out, "[INFO]: ", flag)
	l.warn = log.New(out, "[WARN]: ", flag)
	l.err = log.New(out, "[ERROR]: ", flag)

	return l
}

func NewStdLogger(flags ...int) *Logger {
	return NewLogger(os.Stdout, flags...)
}

func NewFileLogger(name string, flags ...int) *Logger {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(f)
	}
	return NewLogger(f, flags...)
}

func (l *Logger) Info() *log.Logger {
	return l.info
}

func (l *Logger) Warn() *log.Logger {
	return l.warn
}

func (l *Logger) Error() *log.Logger {
	return l.err
}
