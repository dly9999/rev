package config

import (
	"fmt"

	"github.com/Unknwon/goconfig"
	//"github.com/Unknwon/goconfig"
)

type Db struct {
	Server string
	Port   int
	Usr    string
	Psd    string
	Dbname string
} //数据库
type MQconfig struct {
	Mqip     string
	Port     int
	Username string
	Password string
} //Rabitserver Mqserver
type Queconfig struct {
	Per      string
	Name     string
	Routkey  string
	Exchange string
	Constr   string
} //que Mqqueue

var (
	db          Db
	mq          MQconfig
	que         Queconfig
	queforbrush Queconfig
)
var cfg *goconfig.ConfigFile
var Quekeyvalue = make(map[string]string, 10)

func init() {
	cfg, err := goconfig.LoadConfigFile("../conf.ini")
	if err != nil {
		//panic("错误")
		fmt.Println(err)
	}
	var er error
	db.Server, er = cfg.GetValue("DB", "ip")
	if er != nil {
		fmt.Println(er)
	}
	db.Port, er = cfg.Int("DB", "port")
	if er != nil {
		fmt.Println(er)
	}
	db.Usr, er = cfg.GetValue("DB", "username")
	if er != nil {
		fmt.Println(er)
	}
	db.Psd, er = cfg.GetValue("DB", "password")
	if er != nil {
		fmt.Println(er)
	}
	db.Dbname, er = cfg.GetValue("DB", "Dbname")
	if er != nil {
		fmt.Println(er)
	}
	fmt.Println(db)
	mq.Mqip, er = cfg.GetValue("Mqserver", "mqip")
	if er != nil {
		fmt.Println(er)
	}
	mq.Port, er = cfg.Int("Mqserver", "port")
	if er != nil {
		fmt.Println(er)
	}
	mq.Username, er = cfg.GetValue("Mqserver", "username")
	if er != nil {
		fmt.Println(er)
	}
	mq.Password, er = cfg.GetValue("Mqserver", "password")
	if er != nil {
		fmt.Println(er)
	}
	fmt.Println(mq)
	keyvalue, er := cfg.GetSection("Que_Groutelist")
	if er != nil {
		panic(er)
	}
	fmt.Println("a=", keyvalue)
	for k, v := range keyvalue {
		Quekeyvalue[k] = v
	}
	que, er = getquefconf(cfg, "Mqqueue", "per", "Mqqueue", "name")
	queforbrush, er = getquefconf(cfg, "Mqqueueforbrush", "per", "Mqqueueforbrush", "name")
	fmt.Println(que)
}

func RerCgval() (Db, MQconfig, Queconfig, map[string]string) {
	return db, mq, que, Quekeyvalue

}
func getquefconf(cfg *goconfig.ConfigFile, quper string, qupername string, que string, quename string) (qucf Queconfig, er error) {
	var qu Queconfig
	qu.Per, er = cfg.GetValue(quper, qupername)
	if er != nil {
		fmt.Println(er)
		//que = nil
		er = er
		return
	}
	qu.Name, er = cfg.GetValue(que, quename)
	if er != nil {
		//que = nil
		er = er
		return
	}
	//fmt.Println(qucf)
	return qu, nil
}
func getquekeyvalue() {

}

/*func GetqueFromConfig(session string) (keylist []string, que map[string]string) {
	keylist, er := cfg.GetKeyList()

}*/
