package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	_ "time/tzdata"
)

type Logger struct {
	base *log.Logger
	loc  *time.Location
}

func NewLogger() *Logger {
	return &Logger{
		base: log.New(os.Stdout, "", 0),
		loc:  loadColomboLocation(),
	}
}

func loadColomboLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Colombo")
	if err != nil {
		return time.FixedZone("+0530", 5*60*60+30*60)
	}
	return loc
}

func (l *Logger) Info(msg string, kv ...any) {
	l.write("INFO", msg, kv...)
}

func (l *Logger) Warn(msg string, kv ...any) {
	l.write("WARN", msg, kv...)
}

func (l *Logger) Error(msg string, kv ...any) {
	l.write("ERROR", msg, kv...)
}

func (l *Logger) Fatal(msg string, kv ...any) {
	l.write("FATAL", msg, kv...)
	os.Exit(1)
}

func (l *Logger) write(level, msg string, kv ...any) {
	timestamp := time.Now().In(l.loc).Format("2006-01-02 15:04:05")

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("[%s] %s : %s", level, timestamp, msg))

	for i := 0; i+1 < len(kv); i += 2 {
		key, ok := kv[i].(string)
		if !ok || key == "" {
			key = fmt.Sprintf("field_%d", i)
		}
		builder.WriteString(fmt.Sprintf(" %s=%v", key, kv[i+1]))
	}

	if len(kv)%2 == 1 {
		builder.WriteString(fmt.Sprintf(" extra=%v", kv[len(kv)-1]))
	}

	l.base.Println(builder.String())
}

var Log = NewLogger()
