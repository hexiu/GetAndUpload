package main

import (
	// "bufio"
	// "code.google.com/p/mahonia"
	"database/sql"
	"fmt"
	_ "go-odbc/driver"
	// "io"
	// "os"
	// "strconv"
	// "time"
)

type Dept struct {
	deptnum string
	isExist bool
}

//部门信息组成的字典
var dept map[int]Dept = make(map[int]Dept, 0)

//部门存在的凭证，用来创建对应的目录
var deptjud map[string]bool = make(map[string]bool, 0)

var (
	DSN      string = "renwu"
	SERVER   string = "222.24.24.91"
	USERNAME string = "Admin"
	PASSWORD string = ""
	// FilePath   string = "C:/QianDao"
	DATABASE string = "renwu2"
	QIANZHUI string = "rmsTA_"
	ALL_NAME string = "AllMessge"
	// UpdateHour int    = 00
	// UpdateMin  int    = 01
)

var localDebug bool = true

func main() {

	// var dec mahonia.Decoder
	// var enc mahonia.Encoder
	if localDebug {
		fmt.Println("driver={SQL Server};DSN=" + DSN + ";SERVER=" + SERVER + ";Database=" + DATABASE + ";UID=" + USERNAME + ";PWD=" + PASSWORD)
	}
	conn, err := sql.Open("odbc", "driver={SQL Server};DSN="+DSN+";SERVER="+SERVER+";Database="+DATABASE+";UID="+USERNAME+";PWD="+PASSWORD) //

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// fmt.Println(conn)
	defer conn.Close()

	stmt, err := conn.Prepare("SELECT deptname,badgenumber,userid from USERINFO left join DEPARTMENTS ON USERINFO.DEFAULTDEPTID = DEPARTMENTS.DEPTID")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// fmt.Println(stmt)
	defer stmt.Close()

	row, err := stmt.Query()
	if err != nil {
		fmt.Println(err.Error())
		// writeErrorLog("Query Error" + err.Error())
		return
	}
	// fmt.Println(row)
	defer row.Close()

	//dept flag
	// flag := 0

	for row.Next() {
		var deptnum string
		var userid string
		var idstr string
		// var dept1 Dept
		if err := row.Scan(&deptnum, &userid, &idstr); err == nil {

			fmt.Println(deptnum, userid, idstr)
			/*
				    		dec = mahonia.NewDecoder("gbk")
							if ret, ok := dec.ConvertStringOK(useridstr); ok {

								fmt.Println("GBK to UTF-8: ", ret, " bytes:", useridstr)

							}
							// fmt.Println(idint)
			*/ //  idstr := strconv.Itoa(idint)

			// test, _ := strconv.ParseInt(useridstr, 10, 32)
			// fmt.Println("UserId======", test)
			// userid := strconv.Itoa(useridint)
			// userid, _ := strconv.ParseInt(useridstr, 10, 32)
			//  fmt.Println(useridstr)
			//idstr := strconv.Itoa(idint)
			//      id, _ := strconv.ParseInt(idstr, 10, 32)

		}

	}
}
