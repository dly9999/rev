package Db

import (
	"cloud/config"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Dbc struct {
	Sql *sqlx.DB
}

var db Dbc

type Person struct {
	UserId   int    `db:"user_id"`
	Username string `db:"username"`
	Sex      string `db:"sex"`
	Email    string `db:"email"`
}

func init() {
	var connects config.Db
	connects, _, _, _ = config.RerCgval()
	costr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", connects.Usr, connects.Psd, connects.Server, connects.Port, connects.Dbname)
	fmt.Println(costr)
	db1, err := sqlx.Open("mysql", costr) //"root:root@tcp(127.0.0.1:3306)/test"
	if err != nil {
		fmt.Println(err)
	}

	db.Sql = db1
}
func ReDb() Dbc {
	return db
}
func (db Dbc) Sel(va interface{}, questring string, args ...interface{}) error {
	//var pers []Person

	er := db.Sql.Select(va, questring, args[0])
	if er != nil {
		fmt.Println(er, va, questring)
		return er
	}
	return nil
}
func (db Dbc) Inst(qstr string, args []interface{}) error {
	_, er := db.Sql.Exec(qstr, args)
	if er != nil {
		fmt.Println(er)
		return er
	}
	return nil
}
func (db Dbc) Cls() {
	db.Sql.Close()
	fmt.Println("资源释放")
}
