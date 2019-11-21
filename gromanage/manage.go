package gromanage

import (
	"cloud/Db"
	"cloud/Prosess"
	"cloud/config"
	"fmt"

	//"log"
	"time"
)

var db_ma = Db.ReDb()
var Mq_pro config.MQconfig

var que = make(map[string][]config.Queconfig, 1)
var keyvalue map[string]string
var quelen = make(map[string]int, 1)

func init() {
	_, Mq_pro, _, keyvalue = config.RerCgval()
	for key, _ := range keyvalue {
		quelen[key] = 0
	}
}
func Getqueues(mesasgetype string) {
	r, err := db_ma.Sql.Query("select a.tenantId,queueName from mescloud.FactoryQueueAdmin AS a where messageType = \"" + mesasgetype + "\"and isDelete=0")
	if err != nil {
		panic("Db err:" + err.Error())
	}
	defer r.Close()
	for r.Next() {
		var rp bool
		rp = false
		var queconfig config.Queconfig
		r.Scan(&queconfig.Per, &queconfig.Name)
		fmt.Println(queconfig.Name)
		queconfig.Constr = fmt.Sprintf("amqp://%s:%s@%s:%d/%s", Mq_pro.Username, Mq_pro.Password, Mq_pro.Mqip, Mq_pro.Port, queconfig.Per)
		for _, v := range que[mesasgetype] {
			if v.Name == queconfig.Name {
				fmt.Println("repeate ...")
				rp = true
				break
			}
		}
		if !rp {
			que[mesasgetype] = append(que[mesasgetype], queconfig)
		}
	}
	/*if len(que) > 0 {
		for _, v := range que {

		}
	}*/
	//quelen[mesasgetype] = 0
	//log.Println("dasdas", que[mesasgetype])
}
func Getqueueslist() {
	/*for {
		Getqueues("aaa")
		go Prosess.DoProsess().Procsessdata()
	}*/

	for {
		//defer Prosess.Procsessrecover()
		fmt.Println("quekeyvalue", keyvalue)
		for key, v := range keyvalue {
			Getqueues(key)
			if len(que[key]) > 0 && quelen[key] == 0 {
				for _, va := range que[key] {
					var ch = make(chan interface{}, 1024)
					var a = Prosess.DoProsess(v)
					fmt.Println("aaaa", a)
					if a != nil {
						quelen[key] += 1
						go a.Procsessdata(va.Constr, va.Name, ch)
						//go Prosess.DoProsess(v).Procsessdata(va.Constr, va.Name, ch)
					}
				}
			} else if len(que[key]) > 0 && quelen[key] < len(que[key]) {
				var pt = quelen[key] - 1
				for quelen[key] < len(que[key]) {
					var ch = make(chan interface{}, 1024) //计算时从基础往上加防止一次增加两个以上
					pt += 1
					va := que[key][pt]
					fmt.Println("add ")
					var a = Prosess.DoProsess(v)
					go a.Procsessdata(va.Constr, va.Name, ch)
					quelen[key] += 1
				}
			} else {
			}
		}
		time.Sleep(time.Second * 8)
		fmt.Println("again")
	}
}
