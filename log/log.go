package log

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	LevelFlags = [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	TRACE = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

const tunnelSizeDefault = 1024

type Record struct {
	time  string
	code  string
	info  string
	level int
}

func (r *Record) String() string {
	return fmt.Sprintf("[%s][%s][%s] %s\n", LevelFlags[r.level], r.time, r.code, r.info)
}

type Writer interface {
	Init() error
	Write(*Record) error
}

type Rotator interface {
	Rotate() error
	SetPathPattern(string) error
}

type Flusher interface {
	Flush() error
}

type Logger struct {
	writers     []Writer
	tunnel      chan *Record
	level       int
	lastTime    int64
	lastTimeStr string
	c           chan bool
	layout      string
	recordPool  *sync.Pool
}

func NewLogger() *Logger {
	if loggerDefault != nil && takeUp == false {
		takeUp = true //默认启动标志
		return loggerDefault
	}
	l := new(Logger)
	l.writers = []Writer{}
	l.tunnel = make(chan *Record, tunnelSizeDefault)
	l.c = make(chan bool, 2)
	l.level = DEBUG
	l.layout = "2006/01/02 15:04:05"
	l.recordPool = &sync.Pool{New: func() interface{} {
		return &Record{}
	}}
	go boostrapLogWriter(l)

	return l
}

func (l *Logger) Register(w Writer) {
	if err := w.Init(); err != nil {
		panic(err)
	}
	l.writers = append(l.writers, w)
}

func (l *Logger) SetLevel(lvl int) {
	l.level = lvl
}

func (l *Logger) SetLayout(layout string) {
	l.layout = layout
}

func (l *Logger) Trace(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(TRACE, fmt, args...)
}

func (l *Logger) Debug(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(DEBUG, fmt, args...)
}

func (l *Logger) Warn(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(WARNING, fmt, args...)
}

func (l *Logger) Info(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(INFO, fmt, args...)
}

func (l *Logger) Error(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(ERROR, fmt, args...)
}

func (l *Logger) Fatal(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(FATAL, fmt, args...)
}

func (l *Logger) Close() {
	close(l.tunnel)
	<-l.c
	for _, w := range l.writers {
		if f, ok := w.(Flusher); ok {
			if err := f.Flush(); err != nil {
				log.Println(err)
			}
		}
	}
}

func (l *Logger) deliverRecordToWriter(level int, format string, args ...interface{}) {
	var inf, code string

	if level < l.level {
		return
	}

	if format != "" {
		inf = fmt.Sprintf(format, args...)
	} else {
		inf = fmt.Sprint(args...)
	}

	// source code, file and line num
	_, file, line, ok := runtime.Caller(2)
	if ok {
		code = path.Base(file) + ":" + strconv.Itoa(line)
	}

	// format time
	now := time.Now()
	if now.Unix() != l.lastTime {
		l.lastTime = now.Unix()
		l.lastTimeStr = now.Format(l.layout)
	}
	r := l.recordPool.Get().(*Record)
	r.info = inf
	r.code = code
	r.time = l.lastTimeStr
	r.level = level

	l.tunnel <- r
}

func boostrapLogWriter(logger *Logger) {
	if logger == nil {
		panic("logger is nil")
	}

	var (
		r  *Record
		ok bool
	)

	if r, ok = <-logger.tunnel; !ok {
		logger.c <- true
		return
	}

	for _, w := range logger.writers {
		if err := w.Write(r); err != nil {
			log.Println(err)
		}
	}

	flushTimer := time.NewTimer(time.Millisecond * 500)
	rotateTimer := time.NewTimer(time.Second * 10)

	for {
		select {
		case r, ok = <-logger.tunnel:
			if !ok {
				logger.c <- true
				return
			}
			for _, w := range logger.writers {
				if err := w.Write(r); err != nil {
					log.Println(err)
				}
			}

			logger.recordPool.Put(r)

		case <-flushTimer.C:
			for _, w := range logger.writers {
				if f, ok := w.(Flusher); ok {
					if err := f.Flush(); err != nil {
						log.Println(err)
					}
				}
			}
			flushTimer.Reset(time.Millisecond * 1000)

		case <-rotateTimer.C:
			for _, w := range logger.writers {
				if r, ok := w.(Rotator); ok {
					if err := r.Rotate(); err != nil {
						log.Println(err)
					}
				}
			}
			rotateTimer.Reset(time.Second * 10)
		}
	}
}

// default logger
var (
	loggerDefault *Logger
	takeUp        = false
)

func SetLevel(lvl int) {
	defaultLoggerInit()
	loggerDefault.level = lvl
}

func SetLayout(layout string) {
	defaultLoggerInit()
	loggerDefault.layout = layout
}

func Trace(fmt string, args ...interface{}) {
	defaultLoggerInit()
	loggerDefault.deliverRecordToWriter(TRACE, fmt, args...)
}

func Debug(fmt string, args ...interface{}) {
	defaultLoggerInit()
	loggerDefault.deliverRecordToWriter(DEBUG, fmt, args...)
}

func Warn(fmt string, args ...interface{}) {
	defaultLoggerInit()
	loggerDefault.deliverRecordToWriter(WARNING, fmt, args...)
}

func Info(fmt string, args ...interface{}) {
	defaultLoggerInit()
	loggerDefault.deliverRecordToWriter(INFO, fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	defaultLoggerInit()
	loggerDefault.deliverRecordToWriter(ERROR, fmt, args...)
}

func Fatal(fmt string, args ...interface{}) {
	defaultLoggerInit()
	loggerDefault.deliverRecordToWriter(FATAL, fmt, args...)
}

func Register(w Writer) {
	defaultLoggerInit()
	loggerDefault.Register(w)
}

func Close() {
	defaultLoggerInit()
	loggerDefault.Close()
	loggerDefault = nil
	takeUp = false
}

func defaultLoggerInit() {
	if takeUp == false {
		loggerDefault = NewLogger()
	}
}
