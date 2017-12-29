package main

import (
	"database/sql"
	_ "go-oci8"
	"log"
	"os"
	"fmt"
	"myutil"
)

func main() {

	myConfig := new(myutil.Config)
	myConfig.InitConfig("c:/config.ini")
	fmt.Println(myConfig.Read("default", "path"))
	fmt.Printf("%v", myConfig.Mymap)

	fmt.Printf("%v", myConfig.Mymap["path"])

	path := myutil.GetCurrentDirectory()

	fmt.Println(path)

	files, err := myutil.ListDir(path, ".txt")
	fmt.Println(files, err)



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
	defer db.Close()


	rows, err  := db.Query("select  * from fin_base_contract t ")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close();


	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range columns {
		fmt.Println(columns[i])
		scanArgs[i] = &values[i]
	}


	fmt.Println("-----------------------------------")

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {

			// Here we can check if the value is nil (NULL value)


			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Print(columns[i], ": ", value + "   ")
		}
		fmt.Println("")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
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


