package main

import (
	"database/sql"
	"fmt"
	// "go-odbc"
	_ "go-odbc/driver"
	"strconv"
	// "strings"
	"bufio"
	"io"
	"os"
	"time"
)

type User struct {
	uid      int
	userid   string
	deptment string
}

var user map[int]User = make(map[int]User, 0)

//var user map[int]User = make(map[int]User, 0)

type Info struct {
	uid   int
	inout string
	date  string
	time  string
}

var info map[int]Info = make(map[int]Info, 0)

//var info map[int]Info = make(map[int]Info, 0)

type Dept struct {
	deptnum string
	isExist bool
}

var dept map[int]Dept = make(map[int]Dept, 0)
var deptjud map[string]bool = make(map[string]bool, 0)

var Conf map[string]string = make(map[string]string)

var Update bool

/*
const (
	FilePath        = "C:/QianDao"
	DSN      string = "renwu"
	SERVER   string = "222.24.24.91"
	USERNAME string = "Admin"
	PASSWORD string = ""

	DATABASE   string = "renwu2"
	QIANZHUI   string = "rmsTA_"
	ALL_NAME   string = "AllMessge"
	UpdateHour int    = 00
	UpdateMin  int    = 01
)
*/

var (
	DSN        string = "renwu"
	SERVER     string = "222.24.24.91"
	USERNAME   string = "Admin"
	PASSWORD   string = ""
	FilePath          = "C:/QianDao"
	DATABASE   string = "renwu2"
	QIANZHUI   string = "rmsTA_"
	ALL_NAME   string = "AllMessge"
	UpdateHour int    = 00
	UpdateMin  int    = 01
)

func main() {

	createHomeDir()
	readConfDir()
	initConf()
	fmt.Println(Conf)
	for true {

		//获取数据
		getData()

		//创建各部门的目录
		createPartDir()

		message()
		//	if time.Now().Hour() != 1 && time.Now().Minute() != 41 {
		time.Sleep(60 * time.Second)
		//	}

	}
}

func defaultConf() {

}

func initConf() {
	if Conf["DSN"] != "" {
		DSN = Conf["DSN"]

	}
	if Conf["SERVER"] != "" {
		SERVER = Conf["SERVER"]
	}
	if Conf["USERNAME"] != "" {
		USERNAME = Conf["USERNAME"]
	}
	if Conf["PASSWORD"] != "" {
		PASSWORD = Conf["PASSWORD"]
	}

	if Conf["DATABASE"] != "" {
		DATABASE = Conf["DATABASE"]
	}
	if Conf["QIANZHUI"] != "" {
		QIANZHUI = Conf["QIANZHUI"]
	}
	if Conf["ALL_NAME"] != "" {
		ALL_NAME = Conf["ALL_NAME"]
	}
	Hour, err := strconv.ParseInt(Conf["UpdateHour"], 10, 32)
	if err != nil {
		writeErrorLog(err.Error())
	}
	Min, err := strconv.ParseInt(Conf["UpdateMin"], 10, 32)
	if err != nil {
		writeErrorLog(err.Error())
	}
	UpdateHour = int(Hour)
	UpdateMin = int(Min)
}

func readConfDir() {
	dirName := "conf"
	dir, err := os.Open(dirName)
	if err != nil {
		err := os.Mkdir(dirName, os.ModePerm)
		if err != nil {
			writeErrorLog(err.Error() + "  Create Conf Dir " + FilePath + " Error")
		}
	}
	dir.Close()

	fileName := dirName + "/" + "config.txt"
	file, err := os.Open(fileName)
	if err != nil {
		writeErrorLog(err.Error())
		return
	}
	defer file.Close()

	inputReader := bufio.NewReader(file)
	lineCounter := 0
	for {
		inputString, readerError := inputReader.ReadString('\n')
		//inputString,  := inputReader.ReadBytes('\n')             if readerError == io.EOF {
		//fmt.Printf("%d : %s", lineCounter, inputString)
		lenStr := len(inputString)
		for i, v := range inputString {
			if inputString[0:2] == "//" {
				continue
			}
			if v == '=' {
				Conf[inputString[:i]] = inputString[i+1 : lenStr-2]
			}
		}
		if readerError == io.EOF {
			//fmt.Println(Conf)
			return
		}

		lineCounter++
	}

}

func createHomeDir() {
	// fmt.Println(FilePath)
	logPath := "log"
	logdir, err := os.Open(logPath)
	if err != nil {
		writeErrorLog(err.Error())
		err := os.Mkdir(logPath, os.ModePerm)
		if err != nil {
			writeErrorLog("Create Log Dir " + logPath + " Error")
		}
	}
	logdir.Close()
	logPath = "log/run"
	logdir, err = os.Open(logPath)
	if err != nil {
		writeErrorLog(err.Error())
		err := os.Mkdir(logPath, os.ModePerm)
		if err != nil {
			writeErrorLog("Create Log Dir " + logPath + " Error")
		}
	}
	logdir.Close()
	logPath = "log/error"
	logdir, err = os.Open(logPath)
	if err != nil {
		writeErrorLog(err.Error())
		err := os.Mkdir(logPath, os.ModePerm)
		if err != nil {
			writeErrorLog("Create Log Dir " + logPath + " Error")
		}
	}
	logdir.Close()

	dir, err := os.Open(FilePath)
	if err != nil {
		writeErrorLog(err.Error())
		err := os.Mkdir(FilePath, os.ModePerm)
		if err != nil {
			writeErrorLog("Create Home Dir " + FilePath + " Error")
		}
	}
	dir.Close()
}

//将信息汇总在这里
func message() {
	type mesg []string
	mess := make(map[string]mesg, 0)
	messkey := make(map[string]int)
	infolen := len(info)
	deptmaplen := len(dept)
	m_all := make(mesg, infolen)
	k := 0
	for j := 0; j < deptmaplen; j++ {
		m := make(mesg, infolen)
		for i := 0; i < infolen; i++ {
			message := user[info[i].uid].deptment + "," + user[info[i].uid].userid + "," + info[i].inout + "," + info[i].date + "," + info[i].time + "\n"
			if k != infolen {
				m_all[k] = message
				k++

			}

			if user[info[i].uid].deptment == dept[j].deptnum {

				// fmt.Println("................................")
				m[messkey[dept[j].deptnum]] = message

				messkey[dept[j].deptnum]++
			} else {
				continue
			}

		}
		//fmt.Println(m)
		mess[dept[j].deptnum] = m
		write(dept[j].deptnum, m, messkey[dept[j].deptnum])
		writeRunLog("Write Department " + dept[j].deptnum + " success.")
	}
	writeAll(m_all)
	if Update {
		writeRunLog("Update lastDay Success.")
		Update = false
	}
}

func writeAll(mess []string) {
	dirname := FilePath + "/" + "AllMessage"
	dir, err := os.Open(dirname)
	if err != nil {
		err := os.Mkdir(dirname, os.ModePerm)
		if err != nil {
			writeErrorLog("Create Dir " + dirname + " Error!")
		}
	}
	dir.Close()
	filename := FilePath + "/" + "AllMessage" + "/" + QIANZHUI + time.Now().String()[0:10] + ALL_NAME + ".csv"
	file, err := os.Open(filename)
	if err != nil {
		file, err := os.Create(filename)
		if err != nil {
			writeErrorLog("Create File " + filename + " Error!")
		}
		file.Close()

	}
	file.Close()

	file, err = os.OpenFile(filename, os.O_CREATE, os.ModePerm)
	if err != nil {
		writeErrorLog(err.Error())
	}
	for _, v := range mess {
		file.WriteString(v)
	}
	writeRunLog("WriteAllMessage Success.")
	defer file.Close()
}

//从数据库获取数据，并进行相应的处理
func getData() {
	fmt.Println("driver={SQL Server};DSN=" + DSN + ";SERVER=" + SERVER + ";Database=" + DATABASE + ";UID=" + USERNAME + ";PWD=" + PASSWORD)
	conn, err := sql.Open("odbc", "driver={SQL Server};DSN="+DSN+";SERVER="+SERVER+";Database="+DATABASE+";UID="+USERNAME+";PWD="+PASSWORD) //

	if err != nil {
		writeErrorLog("Connecting Error")
		return
	}
	// fmt.Println(conn)
	defer conn.Close()
	stmt, err := conn.Prepare("SELECT deptname,badgenumber,userid from USERINFO left join DEPARTMENTS ON USERINFO.DEFAULTDEPTID = DEPARTMENTS.DEPTID")
	if err != nil {
		writeErrorLog("Query Error" + err.Error())
		return
	}
	// fmt.Println(stmt)
	defer stmt.Close()

	row, err := stmt.Query()
	if err != nil {
		writeErrorLog("Query Error" + err.Error())
		return
	}
	// fmt.Println(row)
	defer row.Close()

	//dept flag
	flag := 0

	for row.Next() {
		var deptnum string
		var id string
		var uid string
		var dept1 Dept
		if err := row.Scan(&deptnum, &id, &uid); err == nil {

			if deptjud[deptnum] == false {
				deptjud[deptnum] = true
				dept1.isExist = true
				dept1.deptnum = deptnum
				dept[flag] = dept1
				flag++
			}

			i, _ := strconv.ParseInt(uid, 10, 64)
			j := int(i)
			var a User
			a.uid = j
			a.userid = id
			a.deptment = deptnum
			user[j] = a
		}
		fmt.Println(user)
	}

	if time.Now().Hour() == UpdateHour && time.Now().Minute() == UpdateMin {
		Update = true
	}
	if Update {
		stmt, err = conn.Prepare("SELECT checktype,checktime,userid from CHECKINOUT where DateDiff(dd,CHECKTIME,getdate())=1")
		if err != nil {
			writeErrorLog("Query Error" + err.Error())
			return
		}
		// fmt.Println(stmt)
		defer stmt.Close()

		row, err = stmt.Query()
		if err != nil {
			writeErrorLog("Query Error" + err.Error())
			return
		}
		// fmt.Println(row)
		defer row.Close()

		index := 0
		for row.Next() {
			var checktype string
			var checktime string
			var uid string
			if err := row.Scan(&checktype, &checktime, &uid); err == nil {

				switch checktype {
				case "I":
				case "O":
				default:
					continue
				}

				i, _ := strconv.ParseInt(uid, 10, 64)
				j := int(i)
				var b Info
				b.uid = j
				b.inout = checktype
				b.date = checktime[0:4] + checktime[5:7] + checktime[8:10]
				b.time = checktime[11:13] + checktime[14:16]
				info[index] = b
				index++
			}
			//			fmt.Println(info)

		}
		fmt.Printf("%s\n", "finish")
		return
	} else {
		// stmt, err = conn.Prepare("SELECT checktype,checktime,userid from CHECKINOUT ")
		stmt, err = conn.Prepare("SELECT checktype,checktime,userid from CHECKINOUT where DateDiff(dd,CHECKTIME,getdate())=0")
		if err != nil {
			writeErrorLog("Query Error" + err.Error())
			return
		}
		// fmt.Println(stmt)
		defer stmt.Close()

		row, err = stmt.Query()
		if err != nil {
			writeErrorLog("Query Error" + err.Error())
			return
		}
		// fmt.Println(row)
		defer row.Close()

		index := 0
		for row.Next() {
			var checktype string
			var checktime string
			var uid string
			if err := row.Scan(&checktype, &checktime, &uid); err == nil {
				//fmt.Println(checktype, checktime, uid)
				// fmt.Println("***************...")

				switch checktype {
				case "I":
				case "O":
				default:
					continue
				}

				i, _ := strconv.ParseInt(uid, 10, 64)
				j := int(i)
				var b Info
				b.uid = j
				b.inout = checktype
				b.date = checktime[0:4] + checktime[5:7] + checktime[8:10]
				b.time = checktime[11:13] + checktime[14:16]
				info[index] = b
				index++
			}
		}
		//fmt.Printf("%s\n", "finish====")
		return
	}
}

//创建部门目录
func createPartDir() {

	deptMapLen := len(dept)
	for i := 0; i < deptMapLen; i++ {
		dir, err := os.Open(FilePath + "/" + dept[i].deptnum)
		if err != nil {
			err := os.Mkdir(FilePath+"/"+dept[i].deptnum, os.ModePerm)
			if err != nil {
				writeErrorLog("Mkdir " + FilePath + dept[i].deptnum + " Error")

			}
		}

		dir.Close()

	}

}

//创建文件，函数
func createFile(path, name string) {
	filepath := path + name
	file, err := os.Open(filepath)
	if err != nil {
		_, err = os.Create(filepath)
		if err != nil {
			writeErrorLog("Create " + filepath + " Error")
		}
	}
	defer file.Close()
}

//写入相应的文件
func write(deptment string, message []string, lenth int) {
	var filename string
	if Update {
		filename = QIANZHUI + info[0].date + ".csv"

	} else {
		filename = QIANZHUI + time.Now().String()[0:4] + time.Now().String()[5:7] + time.Now().String()[8:10] + ".csv"
	}

	filepath := FilePath + "/" + deptment + "/" + filename

	file, err := os.Open(filepath)

	if err != nil {
		file, _ = os.Create(filepath)
		file.WriteString("deptnum\t" + "," + "userid\t" + "," + "in/out\t" + "," + "date\t" + "," + "time\t" + "\n")
		file.Close()
	} else {
		file.Close()
	}

	file, err = os.OpenFile(filepath, os.O_WRONLY, os.ModePerm)
	fmt.Println(lenth)
	for i := 0; i < lenth; i++ {
		file.WriteString(message[i])
		//	fmt.Println(message[i])
	}

	file.Close()
}

func writeErrorLog(logMessage string) {
	logfile := "log/error/" + "log_" + time.Now().String()[0:10] + ".log"
	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("Create Log File Error" + time.Now().String())
		return
	}
	file.WriteString(time.Now().String() + " " + logMessage + "\n")
	defer file.Close()
}

func writeRunLog(logMessage string) {
	logfile := "log/run/" + "log_" + time.Now().String()[0:10] + ".log"
	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("Create Log File Error" + time.Now().String())
		return
	}
	file.WriteString(time.Now().String() + " " + logMessage + "\n")
	defer file.Close()
}
