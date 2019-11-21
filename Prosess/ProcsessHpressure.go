package Prosess

import (
	"cloud/Cloudlog"
	"cloud/Db"
	"cloud/RabitMq"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
	//"github.com/jmoiron/sqlx"
)

type Hpressure struct {
	OrgnanizeType  string `json:"OrganizeType"`
	OrgnanizeCode  string `json:"OrganizeCode"`
	OverStockCount int
}

func init() {
	/*var a = Hpressure{"a", "b", 12}
	jsdata, er := json.Marshal(a)
	if er != nil {
		log.Println(er)
	}
	log.Println(string(jsdata[:]))*/
	RegisterProc("orstock", &Hpressure{})
}
func (p *Hpressure) Procsessdata(Constrs string, quename string, datachan chan interface{}) {
	go ProcsessHpressuredata(datachan, quename)
	RabitMq.RevMq(Constrs, quename, datachan)
}
func ProcsessHpressuredata(datachan chan interface{}, quname string) {
	for {
		select {
		case v := <-datachan:
			//var by []byte = v
			sqlProcssforHpressure(v, quname)

		case <-time.NewTicker(time.Second * 8).C:
			fmt.Println("watiing for overstock mesage")
		}
	}
}
func sqlProcssforHpressure(arr interface{}, quname string) {
	defer Procsessrecover()
	var slice []Hpressure
	ar, tf := arr.([]byte)
	if tf == true {
		//fmt.Println(string(ar[:]))
		er := json.Unmarshal(ar, &slice)
		if er != nil {
			panic("json er+" + string(ar[:]))
		} else {
			Cloudlog.Logprint("111", slice)
			for _, v := range slice {
				ProcsessHpressure(&v, quname)
				//procsessBrushdate(&slice)
				Cloudlog.Logprint("stock", slice)
			}

		}
	} else {
		log.Println("error", ar)
	}
}
func ProcsessHpressure(data *Hpressure, quname string) {
	dbforHpressure := Db.ReDb()
	con, er := dbforHpressure.Sql.Begin()
	if er != nil {
		Cloudlog.Logprint("dberror", er)
		return
	}
	InsertHpressure(con, data, quname)

}
func InsertHpressure(con *sql.Tx, data *Hpressure, quename string) {
	var userid string
	strslect := " select tenantid from  MESforCloud.FactoryQueueAdmin where queuename = \"" + quename + "\""
	Cloudlog.Logprint("strselect", strslect)
	rw, er := con.Query(strslect)
	if er != nil {
		panic(er)
	}
	for rw.Next() {
		rw.Scan(&userid)
		//log.Println(userid)
	}
	//log.Println(userid)
	rw.Close()
	str := "select * from MESforCloud.Work_TeamOverStock where OrganizeCode= \"" + data.OrgnanizeCode +
		"\" and OrganizeType= \"" + data.OrgnanizeType + "\" and tenantId=\"" + userid + "\""
	Cloudlog.Logprint("OverStock str", str)
	r, er := con.Query(str)
	if er != nil {
		Cloudlog.Logprint("OverStock select", er)
		con.Rollback()
		return
	}
	if r.Next() {
		r.Close()
		ex_r, er := con.Exec("update  MESforCloud.Work_TeamOverStock set OverStockCount=? where OrganizeCode= ? and OrganizeType=? and tenantId=? ",
			data.OverStockCount,
			data.OrgnanizeCode,
			data.OrgnanizeType,
			userid)
		if er != nil {
			Cloudlog.Logprint("Work_TeamOverStock", er)
			con.Rollback()
			return
		}
		_, er = ex_r.RowsAffected()
		if er != nil {
			Cloudlog.Logprint("Work_TeamOverStock", er)
			con.Rollback()
			return
		}
		con.Commit()
		return
	} else {
		ex_r, er := con.Exec("insert into MESforCloud.Work_TeamOverStock ( OrganizeType, OrganizeCode, OverStockCount,tenantId) values (?,?,?,?)",
			data.OrgnanizeType,
			data.OrgnanizeCode,
			data.OverStockCount,
			userid)
		if er != nil {
			Cloudlog.Logprint("Work_TeamOverStock", er)
			con.Rollback()
			return
		}
		_, er = ex_r.LastInsertId()
		if er != nil {
			Cloudlog.Logprint("Work_TeamOverStock", er)
			con.Rollback()
			return
		}
		con.Commit()
		return
	}
}
