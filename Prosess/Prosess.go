package Prosess

import (
	"cloud/Db"
	"cloud/RabitMq"
	"cloud/config"
	"encoding/json"
	"fmt"
	"log"
	"time"
	//"github.com/go-sql-driver/mysql"
	//"github.com/jmoiron/sqlx"
)

var (
	Mq_test  config.MQconfig
	Que_test config.Queconfig
	Mque     string
)

type Employees struct {
	FirstName string
	LastName  string
}
type Ep struct {
	Empl []Employees `json:"employees"`
}

var db = Db.ReDb()
var prochan = make(chan interface{}, 1024)

func init() {
	/*	_, Mq_test, Que_test, _ = config.RerCgval()
		Mque = fmt.Sprintf("amqp://%s:%s@%s:%d/%s", Mq_test.Username, Mq_test.Password, Mq_test.Mqip, Mq_test.Port, Que_test.Per) //"amqp://admin:admin@172.16.172.98:5672/Admin"
		fmt.Println("MQ:", Mque)*/
}
func Prosesstest() {
	RabitMq.RevMq(Mque, Que_test.Name, prochan)
	fmt.Println("ProsessforDb:")
}
func ProsessforDb() {
	fmt.Println("ProsessforDb:")
	for {
		select {
		case v := <-prochan:
			//var by []byte = v
			sqlPross(v)
			//fmt.Println("v1", string(v[:]))

		case <-time.NewTicker(time.Second * 5).C:
			fmt.Println("watiing for channle mesage")
		}
	}
}
func sqlPross(arr interface{}) {
	var slice Ep
	ar, tf := arr.([]byte)
	if tf == true {
		fmt.Println(string(ar[:]))
		er := json.Unmarshal(ar, &slice)
		if er != nil {
			log.Panicln(er)
			return
		} else {
			fmt.Println("01", slice.Empl)
		}
		for _, v := range slice.Empl {
			//	fmt.Println(v.FirstName)
			insertDb(v)
			time.Sleep(time.Second * 5)
		}
	} else {
		log.Println("error", ar)
	}
}
func insertDb(v Employees) {
	conn, err := db.Sql.Begin()
	if err != nil {
		return
	}
	r, er := conn.Exec("insert into name  (firstname,lastname) VALUES(?,?)", v.FirstName, v.LastName)
	if er != nil {
		log.Println(er)
		conn.Rollback()
		return
	}
	id, er := r.LastInsertId()
	if er != nil {
		log.Println(er)
		conn.Rollback()
		return
	}
	fmt.Println("insert suc", id)
	conn.Commit()
}
