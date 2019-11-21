package test

import (
	"cloud/Db"
	"fmt"
)

type Person struct {
	UserId   int    `db:"user_id"`
	Username string `db:"username"`
	Sex      string `db:"sex"`
	Email    string `db:"email"`
}

func Testsel() {
	str := "select user_id, username, sex, email from person where user_id=?"
	var va []Person
	db := Db.ReDb()
	er := db.Sel(&va, str, 5)
	if er != nil {
		fmt.Println(er)
		return
	}
	fmt.Println(va)
}

/*func TestInsert() {
	str := "insert into person(username, sex, email)values(?, ?, ?)"
	db := Db.ReDb()
	var a []string
	a = append(a, "stu001", "man", "stu01@qq.com")
	db.Inst(str, a)
}*/
