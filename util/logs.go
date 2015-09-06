//Package util
package util

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	BAK
)

var levelStr = []string{"DEBUG", "INFO", "WARN", "ERROR", "BAK"}

type GfLogs struct {
	logger  *log.Logger
	fp      *os.File
	level   int
	dir     string
	prefix  string
	dateStr string
	mutex   *sync.Mutex
}

// NewLogs returns an GfLogs instance.
func NewLogs(dir, prefix string, level int) (*GfLogs, error) {
	if dir == "-" {
		os.Exit(1)
	}

	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0775)
		if err != nil {
			log.Fatalf("ERROR: create dir error. dir: %s error: %s", dir, err.Error())
		}
	}

	glogs := GfLogs{level: level, dir: dir, prefix: prefix}
	glogs.mutex = &sync.Mutex{}
	glogs.dateStr = time.Now().Format("2006-01-02")

	return &glogs, glogs.init()
}

// Close the fd
func (l *GfLogs) Close() {
	if l.fp != nil {
		l.fp.Close()
	}
}

func (l *GfLogs) init() error {
	var blogName = l.prefix + "_" + l.dateStr + ".log"
	logFile := path.Join(l.dir, blogName)
	fp, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		log.Printf("ERROR: open blog error. log: %s error: %s", logFile, err.Error())
		return err
	}

	l.Close()

	l.fp = fp
	l.logger = log.New(fp, "", log.Ldate|log.Ltime)
	return nil
}

func (l *GfLogs) reopen() {
	nowDate := time.Now().Format("2006-01-02")
	l.mutex.Lock()
	if nowDate > l.dateStr {
		l.dateStr = nowDate
		l.init()
	}
	l.mutex.Unlock()
}

func (l *GfLogs) logging(level int, s *string) {
	if level < 0 || level >= len(levelStr) {
		return
	}

	if l.dir == "-" {
		return
	}

	l.reopen()

	l.logger.Println("[" + levelStr[level] + "] " + *s)
}

//Info info level logs
func (l *GfLogs) Info(format string, a ...interface{}) {
	if l.level > INFO {
		return
	}

	s := fmt.Sprintf(format, a...)
	l.logging(INFO, &s)
}

//Debug debug level logs
func (l *GfLogs) Debug(format string, a ...interface{}) {
	if l.level > DEBUG {
		return
	}

	s := fmt.Sprintf(format, a...)
	l.logging(DEBUG, &s)
}

//Warn warning level logs
func (l *GfLogs) Warn(format string, a ...interface{}) {
	if l.level > WARN {
		return
	}

	s := fmt.Sprintf(format, a...)
	l.logging(WARN, &s)
}

//Error error level logs
func (l *GfLogs) Error(format string, a ...interface{}) {
	if l.level > ERROR {
		return
	}

	s := fmt.Sprintf(format, a...)
	_, file, line, ok := runtime.Caller(1)

	if ok {
		_, file = path.Split(file)
		s = fmt.Sprintf("[%s:%d] %s", file, line, s)
	}
	l.logging(ERROR, &s)
}
