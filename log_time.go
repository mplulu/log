package log

import (
	"fmt"
	"time"
)

type LogTime struct {
	startTime time.Time
	logList   []string
}

func NewLogTime() *LogTime {
	return &LogTime{
		startTime: time.Now(),
		logList:   []string{},
	}
}

func (log *LogTime) Log(text string) {
	log.logList = append(log.logList, fmt.Sprintf("%s: %s", text, time.Now().Sub(log.startTime).String()))
}

func (log *LogTime) Print() {
	fmt.Println(log.logList)
}
