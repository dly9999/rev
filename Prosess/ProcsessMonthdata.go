package Prosess

import (
	"cloud/Cloudlog"
	"cloud/Db"
	"cloud/RabitMq"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

type ProMonthdata struct {
}
type MonthData struct {
	//tenantId  string
	OrgCode   string
	OrgName   string
	OrgType   string
	OutputSum string
	MesSort   string
	SortType  string
	StepCode  string
	OwnMonth  string
	OwnYear   string
}
type MonthDatalist struct {
	Monthlist []MonthData
}

func init() {
	RegisterProc("Monthdata", &ProMonthdata{})
	/*	var a = MonthData{"a", "b", "1", "342", "da", "", "123", "2019/6", time.Now().Format("2006")}
		fmt.Println(a)
		js, er := json.Marshal(a)
		if er != nil {
			panic(er)
		}
		log.Println(string(js[:]))*/
}

func (p *ProMonthdata) Procsessdata(Constrs string, quename string, datachan chan interface{}) {
	go ProcsessMonthchan(datachan, quename)
	RabitMq.RevMq(Constrs, quename, datachan)
}
func ProcsessMonthchan(datachan chan interface{}, quname string) {
	for {
		select {
		case v := <-datachan:
			//var by []byte = v
			sqlProcssforMonthdata(v, quname)

		case <-time.NewTicker(time.Second * 8).C:
			fmt.Println("watiing for month mesage")
		}
	}
}
func sqlProcssforMonthdata(arr interface{}, quname string) {
	defer Procsessrecover()
	var slice []MonthData
	ar, tf := arr.([]byte)
	if tf == true {
		//fmt.Println(string(ar[:]))
		er := json.Unmarshal(ar, &slice)
		if er != nil {
			panic("json er+" + string(ar[:]))
		} else {
			log.Println(slice)
			for _, v := range slice {
				ProcsessMonthdata(&v, quname)
				Cloudlog.Logprint("monthdata", v)
			}
			/*	ProcsessMonthdata(&slice, quname)
				Cloudlog.Logprint("monthdata", slice)*/
		}
	} else {
		log.Println("error", ar)
	}
}
func ProcsessMonthdata(data *MonthData, quname string) {
	dbformonthdata := Db.ReDb()
	log.Println(data)
	con, er := dbformonthdata.Sql.Begin()
	if er != nil {
		panic(er)
	}
	if func(snt string) int {
		in, _ := strconv.Atoi(snt)
		return in
	}(data.OrgType) == 1 {
		InsertMothdata(con, "Work_MonthCountDeptment", data, quname)
	} else if func(snt string) int {
		in, _ := strconv.Atoi(snt)
		return in
	}(data.OrgType) == 2 {
		InsertMothdata(con, "Work_MonthCountTeam", data, quname)
	} else if func(snt string) int {
		in, _ := strconv.Atoi(snt)
		return in
	}(data.OrgType) == 3 {
		InsertMothdata(con, "Work_MonthCountEmployee", data, quname)
	} else {
		con.Commit()
		return
	}
}
func InsertMothdata(con *sql.Tx, tablename string, data *MonthData, quename string) {
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
	str := "select * from  " + tablename + " where tenantId =\"" + userid + "\" and OrgCode=\"" +
		data.OrgCode + "\" and OwnMonth=\"" + data.OwnYear + data.OwnMonth + "\" and MESSort=\"" + data.MesSort + "\" and StepCode=\"" + data.StepCode + "\""
	//	Cloudlog.Logprint("str", str)
	row, er := con.Query(str)
	if er != nil {
		con.Rollback()
		Cloudlog.Logprint("month", er)
	}
	if !row.Next() {
		log.Println("11111")
		row.Close()
		r, er := con.Exec("insert into "+tablename+
			" (tenantId,OrgCode,OrgName,OutputSum,MESSort,SortType,OrgType,OwnMonth,StepCode,CreateBy,CreateDate,UpdateBy,UpdateDate)"+
			"values(?,?,?,?,?,?,?,?,?,?,?,?,?)",
			userid,
			data.OrgCode,
			data.OrgName,
			data.OutputSum,
			data.MesSort,
			data.SortType,
			data.OrgType,
			data.OwnYear+data.OwnMonth,
			data.StepCode,
			"",
			//time.Parse("2006-01-02 15:04:05", time.Now().UTC()),
			time.Now().Local(),
			"",
			time.Now().Local())
		if er != nil {
			Cloudlog.Logprint("montdata", er)
			con.Rollback()
			return
		}
		_, er = r.LastInsertId()
		if er != nil {
			Cloudlog.Logprint("montdata", er)
			con.Rollback()
			return

		}
		con.Commit()
		return
	} else {

		row.Close()
		str := "UPDATE " + tablename + " SET OutputSum=?, UpdateDate= ? where tenantId=? and OrgCode=? and OwnMonth= ? and StepCode= ? and MESSort= ?"
		Cloudlog.Logprint("str", str)
		_, er := con.Exec(str,
			data.OutputSum,
			time.Now(),
			userid,
			data.OrgCode,
			data.OwnYear+data.OwnMonth,
			data.StepCode,
			data.MesSort)
		log.Println("2222")
		if er != nil {
			Cloudlog.Logprint("dddd", er)
			con.Rollback()
			return
		}
		/*_, er = r.RowsAffected()
		if er != nil {
			Cloudlog.Logprint("montdata", er)
			con.Rollback()
			return
		}*/
		con.Commit()
		return
	}

}
