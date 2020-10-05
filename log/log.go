package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	errLog = log.New(os.Stdout, "\033[31m[ERROR]\033[0m ", log.LstdFlags | log.Lshortfile)
	infoLog = log.New(os.Stdout, "\033[34m[INFO]\033[0m ", log.LstdFlags | log.Lshortfile)
	loggers = []*log.Logger{errLog, infoLog}
	mu sync.Mutex
)

var (
	Error = errLog.Println
	Errorf = errLog.Printf
	Info = infoLog.Println
	Infof = infoLog.Printf
)

const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

func SetLevel(l int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if l > InfoLevel {
		infoLog.SetOutput(ioutil.Discard)
	}
	if l > ErrorLevel {
		errLog.SetOutput(ioutil.Discard)
	}
}