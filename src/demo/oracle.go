package main

import (
	"database/sql"
	_ "go-oci8"
	"log"
	"os"
	"fmt"
)

func main() {
	// 为log添加短文件名,方便查看行数
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	log.Println("Oracle Driver example")

	os.Setenv("NLS_LANG", "AMERICAN_AMERICA.UTF8")
	os.Setenv("TNS_ADMIN", "D:/evn/instantclient-basic-windows.x64-12.2.0.1.0/instantclient_12_2/tns/")

	// 用户名/密码@实例名  跟sqlplus的conn命令类似
	db, err := sql.Open("oci8", "sitbas/qW5ekFKYTZI=@bas_sit")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("select 3.14, 'foo' from dual")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for rows.Next() {
		var f1 float64
		var f2 string
		rows.Scan(&f1, &f2)
		log.Println(f1, f2) // 3.14 foo
	}
	rows.Close()

	// 先删表,再建表
	/*
	db.Exec("drop table sdata")
	db.Exec("create table sdata(name varchar2(256))")

	db.Exec("insert into sdata values('中文')")
	db.Exec("insert into sdata values('1234567890ABCabc!@#$%^&*()_+')")
	*/

	//rows, err = db.Query("select * from sdata")
	rows, err = db.Query("select t.LOANNO name,t.CUSTOMERNAME from fin_base_contract t ")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close();


	//fmt.Println( rows.Columns())


	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for i := range columns {
		fmt.Println(columns[i])
	}



	/*for rows.Next() {
		var name string
		var customername1 string
		rows.Scan(&name,&customername1)
		fmt.Println(name)
		fmt.Println(customername1)

	}*/
	rows.Close()
}
