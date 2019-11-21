package main

import (
	"cloud/Db"
	_ "cloud/Prosess"
	_ "cloud/config"
	"cloud/gromanage"
	"fmt"
	"time"
)

func main() {
	db := Db.ReDb()
	defer db.Cls()
	fmt.Println(time.Now())
	time.Sleep(time.Second * 3)
	//	fmt.Println(config.RerCgval())
	//	a := Prosess.ProBrush{}
	//gromanage.Getqueues("upbrushdata")
	forever := make(chan bool)
	gromanage.Getqueueslist()
	//go a.ProcsessforBrush()
	<-forever
}
