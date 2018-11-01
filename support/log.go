package support

import (
	"time"
	"fmt"
	"log"
)

const (
	LOG_LEVEL_MSG = iota
	LOG_LEVEL_WARN
	LOG_LEVEL_ERR
)

type Console struct {
	IsWrite bool
}

type WriteHandler struct {
	outputDrive string

}

var cl *Console
var wh *WriteHandler

func init()  {
	cl = &Console{}
	cl.IsWrite = false

	wh = &WriteHandler{}
	wh.outputDrive = "file"
}

func Cl() *Console  {
	return cl
}

func (o *Console)Message(msg string)  {
	o.output(msg, LOG_LEVEL_MSG)
}

func (o *Console)Error(msg string)  {
	o.output(msg, LOG_LEVEL_ERR)
}

func (o *Console)output(msg string, level int)  {
	levelStr := "MSG"
	if level == LOG_LEVEL_WARN {
		levelStr = "WARN"
	} else if level == LOG_LEVEL_ERR {
		levelStr = "ERR"
	}
	if levelStr == "ERR" {
		log.Println(format(msg, levelStr))
	} else {
		fmt.Println(format(msg, levelStr))
	}
	if o.IsWrite == true {

	}
}

func format(msg string, level string) string  {
	return "[[" + getTime() + "  " + level + "]]" + msg
}

func getTime() string {
	return time.Now().Format("2006-01-02 15:04:05")//
}


