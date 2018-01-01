package main

import (
	"database/sql"
	_ "go-oci8"
	"log"
	"os"
	"fmt"
	"myutil"
	"io/ioutil"
	"xlsx"
	"strings"
	"bar"
	"time"
)



func main() {

	var index int = 0










	log.SetFlags(log.Lshortfile | log.LstdFlags)

	files, _ := myutil.ListDir("D:\\export", ".sql")
	if(len(files) < 1) {
		log.Fatal("no sql file!")
	}
	mysql,_ := ioutil.ReadFile(files[0])

	//println(files[0])

	fname := strings.Split((strings.Split(files[0],"."))[0], string(os.PathSeparator))[len(strings.Split((strings.Split(files[0],"."))[0], string(os.PathSeparator)))-1]

	//println(fname)


	//println(string(mysql))

	//os.Rename(files[0], files[0]+".over")

	myConfig := new(myutil.Config)
	myConfig.InitConfig("D:/export/config.ini")
	//fmt.Println(myConfig.Read("default", "path"))
	//fmt.Printf("%v", myConfig.Mymap)

	//path := myutil.GetCurrentDirectory()
	//fmt.Println(path)
	//files, err := myutil.ListDir(path, ".txt")
	//fmt.Println(files, err)

	// 为log添加短文件名,方便查看行数

	os.Setenv("NLS_LANG", myConfig.Read("dbconfig", "NLS_LANG"))
	os.Setenv("TNS_ADMIN", myConfig.Read("dbconfig", "TNS_ADMIN"))
	db, err := sql.Open("oci8", myConfig.Read("dbconfig", "DB_USER")+"/"+myConfig.Read("dbconfig", "DB_PASS")+"@"+myConfig.Read("dbconfig", "DB_HOST"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()




	rows, err  := db.Query(string(mysql))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()




	excelFile := xlsx.NewFile()

	sheet, err := excelFile.AddSheet(fname)
	if err != nil {
		fmt.Printf(err.Error())
	}



	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	headrow := sheet.AddRow()

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range columns {
		//fmt.Println(columns[i])
		scanArgs[i] = &values[i]
		headrow.AddCell().Value = columns[i]
	}








	go func() {



		var countsql = "select count(1) sum from ( " + string(mysql) + " )"
		println(countsql)

		countRow, err  := db.Query(countsql)

		if err != nil {
			log.Fatal(err)
		}
		defer countRow.Close()

		var countAll int
		for countRow.Next() {
			countRow.Scan(&countAll)
		}



		b := bar.New("test").NewBar("exporting ...", countAll)

		b.InitNumber(index)
		var thisindex int  = index

		for  {
			if(thisindex != index && index < countAll) {
				var a  int  = index - thisindex
				b.AddNumber(a)
				thisindex = index
				time.Sleep(time.Second / 200)
			}
			if(index >= countAll){
				break
			}

		}




	}()



	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		headrow = sheet.AddRow()
		var value string
		for _, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//fmt.Print(columns[i], ": ", value + "   ")
			headrow.AddCell().Value = value
		}
		index ++
		//fmt.Println("")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error())
	}


	err = excelFile.Save("D:/export/"+fname+".xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}


}



