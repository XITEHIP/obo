package support

import (
	"time"
	"fmt"
	"log"
	"os"
)

const (
	LOG_LEVEL_MSG = iota
	LOG_LEVEL_WARN
	LOG_LEVEL_ERR
)


type Console struct {
	IsWrite bool
}

var wm *WriterM
var path = "/tmp"
var pre = "obo"

type WriteHandler struct {
	outputDrive string

}

var cl *Console
var wh *WriteHandler

func init()  {
	cl = &Console{}
	cl.IsWrite = true

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
	formatStr := format(msg, levelStr)
	if levelStr == "ERR" {
		log.Println()
	} else {
		fmt.Println(formatStr)
	}
	if o.IsWrite == true {
		WM().writeFile(formatStr + "\n")
	}
}

func format(msg string, level string) string  {
	return "[[" + getTime() + "  " + level + "]]" + msg
}

func getTime() string {
	return time.Now().Format("2006-01-02 15:04:05")//
}

type WriterM struct {
	path string
	flag int
	model os.FileMode
	file *os.File
	pre string
}

func WM() *WriterM  {
	NewWM()
	return wm
}

func NewWM() {
	if wm == nil {
		wm = &WriterM{}
		wm.path = path
		wm.flag = os.O_RDWR|os.O_CREATE|os.O_APPEND
		wm.model = 0644
		wm.pre = pre

		wm.createFile("")
	}
}

func (o * WriterM)writeFile(strContent string)  {
	o.isNewFile()
	buf := []byte(strContent)
	o.file.Write(buf)
}

func (o * WriterM)isNewFile()  {

	file := pre + "_" + time.Now().Format("2006-01-02")
	file = path + "/" + file + ".log"
	_, err := os.Stat(file)
	if err != nil {
		o.createFile(file)
	}
}

func (o * WriterM)createFile(file string) {

	if file == "" {
		file = pre + "_" + time.Now().Format("2006-01-02")
		file = path + "/" + file + ".log"
	}

	of, err := os.OpenFile(file, wm.flag, wm.model)
	if err != nil {
		fmt.Println(err)
		of, _ := os.Create(file)
		o.file = of
	}
	o.file = of
}

func (o * WriterM)Close() {
	o.file.Close()
}




