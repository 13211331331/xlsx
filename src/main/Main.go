package main

import (
	"database/sql"
	_ "go-oci8"
	"log"
	"os"
	"fmt"
	"myutil"
	"io/ioutil"
	"excel"
	"strings"
	"time"
	"runtime"
)



func main() {

	var index = 0

	runtime.GOMAXPROCS(4)


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

	excelFile := excel.NewFile()

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	header := make(map[string]string)

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range columns {
		//fmt.Println(columns[i])
		scanArgs[i] = &values[i]
		header[myutil.CountToExcel(i+1)] = columns[i]
	}

	for k, v := range header {
		excelFile.SetCellValue("Sheet1", k+"1", v)
	}

	go func() {
		var countsql = "select count(1) sum from ( " + string(mysql) + " )"
		countRow, err  := db.Query(countsql)

		if err != nil {
			log.Fatal(err)
		}
		defer countRow.Close()

		var countAll int
		for countRow.Next() {
			countRow.Scan(&countAll)
		}


		var thisindex int  = index
		for  {
			if(thisindex <= countAll) {
				thisindex = index
				//fmt.Println(thisindex)
				str :=  fmt.Sprintf("当前%d/总数%d", thisindex,countAll)
				log.Println(str)
				time.Sleep(time.Second * 10)
			}
			if(thisindex >=countAll){
				break
			}

		}
	}()


	for rows.Next() {

		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}


		var value string
		var rowindex int = 1
		//mapVal := make(map[string]string)
		for _, col := range values {
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			//fmt.Print(columns[i], ": ", value + "   ")
			var setStr string = fmt.Sprintf("%d",index+2)
			 excelFile.SetCellValue("Sheet1", myutil.CountToExcel(rowindex)+ setStr, value)
			rowindex ++
		}
		index ++
		//fmt.Println("")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

	excelFile.SaveAs("D:/export/"+fname+".xlsx")

}



