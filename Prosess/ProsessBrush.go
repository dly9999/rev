package Prosess

import (
	"cloud/Db"
	"cloud/RabitMq"

	"cloud/Cloudlog"
	"cloud/config"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	//	"strings"
	"database/sql"
	"time"
	//"github.com/shopspring/decimal"
	//	"github.com/jmoiron/sqlx"
)

/*

  处理刷卡数据和质检数据，判断IsCheck是否为true，true只处理质检数据，false 只插入刷卡表或者 重复刷卡表

{
    "EdgeBrushCardModel":{
        "DepartModuleId":"321",
        "ReaderId":"321",
        "EmployeeCode":"321",
        "Odid":"321",
        "StepCode":"321",
        "IsRebruch":true,
        "BrushCardDate":"2016-07-22 11:33:19.000",
        "IsCheck":false,
        "OutCheck":true,
        "EmployeeName":"321",
        "ListCheckSize":[
            {
                "CheckedNo":"1",
                "CheckPartName":"2",
                "StandardSizeValue":"3",
                "CheckedSizeValue":"4"
            }
        ],
        "ListCheckReason":[
            {
                "OutQualityCode":"5",
                "Memo":"6",
                "IsDelete":true
            }
        ]
    },
    "SortCode":"a",
    "SortType":"b",
    "SysCode":"c",
    "OrderCode":"d",
    "TenantId":"12wsq"
}
*/
type ProBrush struct {
}
type ProcsessBrush struct {
	EdgeBrushCardModel EdgeBrush
	SortCode           string
	SortType           string
	SysCode            string
	OrderCode          string
	TenantId           string
} //队列的数据类型
type EdgeBrush struct {
	DepartModuleId  string
	ReaderId        string
	EmployeeCode    string
	Odid            string
	StepCode        string
	IsRebruch       bool
	BrushCardDate   string
	IsCheck         bool
	OutCheck        bool
	EmployeeName    string
	ListCheckSize   []EdgeCheck
	ListCheckReason []EdgeCheckReason
}
type EdgeCheck struct {
	CheckedNo         string
	CheckPartName     string
	StandardSizeValue string
	CheckedSizeValue  string
}
type EdgeCheckReason struct {
	OutQualityID int
	Memo         string
	IsDelete     bool
}

var (
	Mq_brush  config.MQconfig
	Que_brush config.Queconfig
	Mquebrush string
)

//数据库
var brushchan = make(chan interface{}, 1024)

func init() {
	/*checkarr = append(checkarr, chec)
	reasonarr = append(reasonarr, reson)
	dge.ListCheckSize = checkarr
	dge.ListCheckReason = reasonarr
	br.EdgeBrushCardModel = dge
	fmt.Println(br)
	jo, er := json.Marshal(&br)
	if er != nil {
		fmt.Println(er)
	} else {
		fmt.Println(string(jo[:]))
	}*/
	/*_, Mq_brush, _, Que_brush = config.RerCgval()
	Que_brush.Constr = fmt.Sprintf("amqp://%s:%s@%s:%d/%s", Mq_brush.Username, Mq_brush.Password, Mq_brush.Mqip, Mq_brush.Port, Que_brush.Per) //"amqp://admin:admin@172.16.172.98:5672/Admin"
	fmt.Println("MQbrush:", Que_brush.Constr)*/
	RegisterProc("ProsessBrush", &ProBrush{}) //初始化注册到内存中

}

func (p *ProBrush) ProcsessforBrush() {
	go ProcsessBrushchan(brushchan)
	RabitMq.RevMq(Que_brush.Constr, Que_brush.Name, brushchan)

	log.Println("end")
}
func (p *ProBrush) Procsessdata(Constrs string, quename string, datachan chan interface{}) {
	go ProcsessBrushchan(datachan)
	RabitMq.RevMq(Constrs, quename, datachan)
	log.Println("end")
}
func ProcsessBrushchan(datachan chan interface{}) {
	for {
		select {
		case v := <-datachan:
			//var by []byte = v
			sqlProcssforbrush(v)

		case <-time.NewTicker(time.Second * 8).C:
			fmt.Println("watiing for brushcard mesage")
		}
	}
}
func sqlProcssforbrush(arr interface{}) {
	//defer Procsessrecover()
	var slice ProcsessBrush
	ar, tf := arr.([]byte)
	if tf == true {
		//fmt.Println(string(ar[:]))
		er := json.Unmarshal(ar, &slice)
		if er != nil {
			log.Println(er)
			panic("json er+" + string(ar[:]))

		} else {

			procsessBrushdate(&slice)
			Cloudlog.Logprint("bruhdata", slice)
		}
	} else {
		log.Println("error", ar)
	}
}
func procsessBrushdate(data *ProcsessBrush) {
	var db_brush = Db.ReDb()
	//defer db_brush.Cls()
	defer Procsessrecover()
	if data == nil {
		return
	}

	con, er := db_brush.Sql.Begin()
	if er != nil {
		log.Println("db error")
		return
	}
	if data.EdgeBrushCardModel.IsCheck == true {
		r, er := con.Exec("insert into MESforCloud.QC_CheckMain (tenantId,QCCheckDate,StepCode,IsOutCheck,EmployCode,EmployeeName,OdID) VALUES (?,?,?,?,?,?,?)",
			data.TenantId,
			func() time.Time {
				tim, _ := time.Parse("2006-01-02 15:04:05", data.EdgeBrushCardModel.BrushCardDate)
				return tim
			}(),
			data.EdgeBrushCardModel.StepCode,
			data.EdgeBrushCardModel.OutCheck,
			data.EdgeBrushCardModel.EmployeeCode,
			data.EdgeBrushCardModel.EmployeeName,
			func() int64 {
				in, _ := strconv.ParseInt(data.EdgeBrushCardModel.Odid, 10, 64)
				return in
			}())
		if er != nil {
			log.Println(er)
			con.Rollback()
			return
		}
		id, err := r.LastInsertId()
		if err != nil {
			con.Rollback()
			log.Println(err)
			return
		}
		for _, v := range data.EdgeBrushCardModel.ListCheckSize {
			r_checksize, er_checksize := con.Exec("insert into MESforCloud.QC_CheckSize (QCCheckID,CheckedNo,CheckPartName,StandardSizeValue,CheckedSizeValue) values(?,?,?,?,?)",
				id,
				v.CheckedNo,
				v.CheckPartName,
				func() float64 {
					fl, _ := strconv.ParseFloat(v.StandardSizeValue, 32)
					//fmt.Println(fl)
					return fl
				}(),
				func() float64 {
					fl, _ := strconv.ParseFloat(v.CheckedSizeValue, 32)
					return fl
				}())
			rowf, _ := r_checksize.RowsAffected()
			if er_checksize != nil || rowf < 1 {
				con.Rollback()
				return
			}
		}
		for _, v1 := range data.EdgeBrushCardModel.ListCheckReason {
			r_checkreason, er_checkreason := con.Exec("insert into MESforCloud.QC_OutCheckReason(QCCheckID,OutQualityID,Memo,IsDelete) values(?,?,?,?)",
				id,
				v1.OutQualityID,
				v1.Memo,
				v1.IsDelete)
			log.Println(er_checkreason)
			rowf, _ := r_checkreason.RowsAffected()
			if er_checkreason != nil || rowf < 1 {
				con.Rollback()
				log.Println(er_checkreason, "rowaffect:", rowf)
				return
			}
		}
		con.Commit()
		log.Println("sucess")
		return
	} else {
		if data.EdgeBrushCardModel.IsRebruch {
			Insertbrush(con, "MESforCloud.BrushCard_Rebrush_Material", data)
			return
		} else {
			str := "select * from MESforCloud.BrushCard_Material where  OdID= \"" + data.EdgeBrushCardModel.Odid + "\"and stepcode=\"" + data.EdgeBrushCardModel.StepCode + "\""
			log.Println(str)
			r_brush, er := con.Query(str)
			if er != nil {
				log.Println(er)
				return
			}
			if func() bool {
				return r_brush.Next()
			}() == true {
				r_brush.Close()
				Insertbrush(con, "MESforCloud.BrushCard_Rebrush_Material", data)
				return
			} else {
				r_brush.Close()
				Insertbrush(con, "MESforCloud.BrushCard_Material", data)
				return
			}

		}
	}
}
func Insertbrush(con *sql.Tx, tablename string, data *ProcsessBrush) {
	defer Procsessrecover()
	r_rebrush, er := con.Exec("insert into "+tablename+"(OrderCode,syscode,tenantId,DepartModuleID,ReaderID,EmpCode,StepCode,OdID,OdCount,OdState,BrushDate,LoadDate,Remarks)values(?,?,?,?,?,?,?,?,?,?,?,?,?)",
		data.OrderCode,
		data.SysCode,
		data.TenantId,
		data.EdgeBrushCardModel.DepartModuleId,
		data.EdgeBrushCardModel.ReaderId,
		data.EdgeBrushCardModel.EmployeeCode,
		data.EdgeBrushCardModel.StepCode,
		func() int64 {
			in, _ := strconv.ParseInt(data.EdgeBrushCardModel.Odid, 10, 64)
			return in
		}(),
		1,
		1,
		func() time.Time {
			tim, _ := time.Parse("2006-01-02 15:04:05", "2014-06-15 08:37:18")
			return tim
		}(),
		func() time.Time {
			tim, _ := time.Parse("2006-01-02 15:04:05", data.EdgeBrushCardModel.BrushCardDate)
			return tim
		}(),
		" ")
	log.Println(er)
	_, er_brush := r_rebrush.LastInsertId()
	if er != nil || er_brush != nil {
		log.Println(er, er_brush)
		con.Rollback()
		return
	}
	_, er_deitail := con.Exec("update MESforCloud.OrderDetail set StepCode=? where odid =?",
		data.EdgeBrushCardModel.StepCode,
		func() int64 {
			in, _ := strconv.ParseInt(data.EdgeBrushCardModel.Odid, 10, 64)
			return in
		}())
	if er_deitail != nil {
		log.Println(er_deitail)
		con.Rollback()
		return
	}
	log.Println("sucess")
	con.Commit()
}
func Procsessrecover() {
	if r := recover(); r != nil {
		Cloudlog.Logprint("error:", r)
	}
}
