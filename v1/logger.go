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

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.info != nil {
		l.info.Printf(format, v...)
	}
}

func (l *Logger) Infoln(v ...interface{}) {
	if l.info != nil {
		l.info.Println(v...)
	}
}

func (l *Logger) Warn() *log.Logger {
	return l.warn
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.warn != nil {
		l.warn.Printf(format, v...)
	}
}

func (l *Logger) Warnln(v ...interface{}) {
	if l.warn != nil {
		l.warn.Println(v...)
	}
}

func (l *Logger) Error() *log.Logger {
	return l.err
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.err != nil {
		l.err.Printf(format, v...)
	}
}

func (l *Logger) Errorln(v ...interface{}) {
	if l.err != nil {
		l.err.Println(v...)
	}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.err != nil {
		l.err.Fatalf(format, v...)
	}
}

func (l *Logger) Fatalln(v ...interface{}) {
	if l.err != nil {
		l.err.Fatalln(v...)
	}
}
