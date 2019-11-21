package Cloudlog

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Cloudlog struct {
}

func createlogfile() *os.File {
	logFile, err := os.OpenFile("./"+time.Now().Format("20060102")+".txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	return logFile
}

func Logprint(str string, a interface{}) {
	file := createlogfile()
	defer file.Close()
	cloulog := log.New(file, "Cdebug", log.Ldate|log.Ltime|log.Lshortfile)

	cloulog.Println(str+":", a)
	return
}
