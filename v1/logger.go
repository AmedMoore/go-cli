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

	flag := log.Ldate | log.Ltime
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

func (ref *Logger) Info() *log.Logger {
	return ref.info
}

func (ref *Logger) Infof(format string, v ...interface{}) {
	if ref.info != nil {
		ref.info.Printf(format, v...)
	}
}

func (ref *Logger) Infoln(v ...interface{}) {
	if ref.info != nil {
		ref.info.Println(v...)
	}
}

func (ref *Logger) Printf(format string, v ...interface{}) {
	ref.Infof(format, v...)
}

func (ref *Logger) Println(v ...interface{}) {
	ref.Infoln(v...)
}

func (ref *Logger) Warn() *log.Logger {
	return ref.warn
}

func (ref *Logger) Warnf(format string, v ...interface{}) {
	if ref.warn != nil {
		ref.warn.Printf(format, v...)
	}
}

func (ref *Logger) Warnln(v ...interface{}) {
	if ref.warn != nil {
		ref.warn.Println(v...)
	}
}

func (ref *Logger) Error() *log.Logger {
	return ref.err
}

func (ref *Logger) Errorf(format string, v ...interface{}) {
	if ref.err != nil {
		ref.err.Printf(format, v...)
	}
}

func (ref *Logger) Errorln(v ...interface{}) {
	if ref.err != nil {
		ref.err.Println(v...)
	}
}

func (ref *Logger) Fatalf(format string, v ...interface{}) {
	if ref.err != nil {
		ref.err.Fatalf(format, v...)
	}
}

func (ref *Logger) Fatalln(v ...interface{}) {
	if ref.err != nil {
		ref.err.Fatalln(v...)
	}
}
