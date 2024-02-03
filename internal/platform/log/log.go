package log

import (
	"fmt"
	"log"
)

func Error(msg string, err error) {
	log.Printf("[ERROR] [msg:%s] [err:%s]\n", msg, err.Error())
}

func Info(msg string, tags ...string) {
	logMsg := fmt.Sprintf("[INFO] [msg:%s]", msg)
	for _, tag := range tags {
		logMsg += fmt.Sprintf(" [%s]", tag)
	}
	log.Println(logMsg)
}
