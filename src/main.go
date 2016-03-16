package main

import (
	"bufio"
	// "code.google.com/p/mahonia"
	"database/sql"
	"fmt"
	_ "go-odbc/driver"
	"io"
	"os"
	"strconv"
	"time"
)

//获取的用户信息
type User struct {
	id       string
	userid   string
	deptment string
}

//用户组成的一个查询字典
var user map[string]User = make(map[string]User, 0)

//获取签到信息
type Info struct {
	uid   string
	inout string
	date  string
	time  string
}

//签到信息记录的字典
var info map[string]Info = make(map[string]Info, 0)

//获取部门信息
type Dept struct {
	deptnum string
	isExist bool
}

//部门信息组成的字典
var dept map[int]Dept = make(map[int]Dept, 0)

//部门存在的凭证，用来创建对应的目录
var deptjud map[string]bool = make(map[string]bool, 0)

//配置信息字典
var Conf map[string]string = make(map[string]string)

//更新标识
var Update bool

var localDebug bool = true

// 配置信息全局变量（有默认值）
var (
	DSN        string = "renwu"
	SERVER     string = "222.24.24.91"
	USERNAME   string = "Admin"
	PASSWORD   string = ""
	FilePath   string = "C:/QianDao"
	DATABASE   string = "renwu2"
	QIANZHUI   string = "rmsTA_"
	ALL_NAME   string = "AllMessge"
	UpdateHour int    = 00
	UpdateMin  int    = 01
)

// 主函数
func main() {
	//localDebug = false

	// 创建签到数据存放主目录
	createHomeDir()
	// 创建配置文件目录
	readConfDir()
	// 加载配置文件
	initConf()
	// fmt.Println(Conf)

	//获取数据
	getData()

	//创建各部门的目录
	createPartDir()

	// 数据处理引擎，并且写入文档
	message()

	time.Sleep(100 * time.Second)
}

// 加载配置文件函数
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

//创建日志文件目录,读取配置文件目录
func readConfDir() {
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

	dirName := "conf"
	dir, err := os.Open(dirName)
	if err != nil {
		/*
			err := os.Mkdir(dirName, os.ModePerm)
			if err != nil {
				writeErrorLog(err.Error() + "  Create Conf Dir " + FilePath + " Error")
			}
		*/
		return
	}
	defer dir.Close()

	fileName := dirName + "/" + "config.txt"
	file, err := os.Open(fileName)
	if err != nil {
		writeErrorLog("No Config File.")
		return
	}
	defer file.Close()

	inputReader := bufio.NewReader(file)
	lineCounter := 0
	for {
		inputString, readerError := inputReader.ReadString('\n')

		// fmt.Println(inputString)
		lenStr := len(inputString)
		for i, v := range inputString {
			if inputString[0:2] == "//" {
				continue
			}
			if v == '=' {
				Conf[inputString[:i]] = inputString[i+1 : lenStr-2]
			}
		}
		//	fmt.Println(Conf)

		if readerError == io.EOF {
			return
		}
		lineCounter++
	}
	if localDebug {
		fmt.Println(Conf)
	}

}

//以及程序存放数据主目录
func createHomeDir() {
	// fmt.Println(FilePath)
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
		for q := 0; q < infolen; q++ {
			// message := user[info[i].uid].deptment + "," + strconv.Itoa(user[info[i].uid].userid) + "," + info[i].inout + "," + info[i].date + "," + info[i].time + "\n"
			i := strconv.Itoa(q)

			//			message := user[info[i].uid].deptment + "," + user[info[i].uid].userid + "," + info[i].inout + "," + info[i].date + "," + info[i].time + "\n"
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
	//	writeRunLog("Write Department  success.")

	writeAll(m_all)
	if Update {
		writeRunLog("Update lastDay Success.")
		Update = false
	}
}

//将当天所有部门的信息写入汇总文件
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
	/*
		if Update {
			filename := FilePath + "/" + "AllMessage" + "/" + QIANZHUI + info[0].date + ALL_NAME + ".csv"
		}
	*/
	filename := FilePath + "/" + "AllMessage" + "/" + QIANZHUI + info["0"].date + ALL_NAME + ".csv"

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
	file.WriteString("deptnum" + "," + "userid" + "," + "in/out" + "," + "date" + "," + "time" + "\n")
	for _, v := range mess {
		file.WriteString(v)
	}
	writeRunLog("Write Day AllMessage Success.")
	defer file.Close()
}

//从数据库获取数据，并进行相应的处理
func getData() {

	//转码
	// var dec mahonia.Decoder
	// var enc mahonia.Encoder

	if localDebug {
		fmt.Println("driver={SQL Server};DSN=" + DSN + ";SERVER=" + SERVER + ";Database=" + DATABASE + ";UID=" + USERNAME + ";PWD=" + PASSWORD)
	}
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
		var useridstr string
		var idstr string
		var dept1 Dept
		if err := row.Scan(&deptnum, &useridstr, &idstr); err == nil {
			/*
				dec = mahonia.NewDecoder("gbk")
				if ret, ok := dec.ConvertStringOK(useridstr); ok {

					fmt.Println("GBK to UTF-8: ", ret, " bytes:", useridstr)

				}
			*/ // fmt.Println(idint)
			//  idstr := strconv.Itoa(idint)

			// test, _ := strconv.ParseInt(useridstr, 10, 32)
			// fmt.Println("UserId======", test)
			// userid := strconv.Itoa(useridint)
			// userid, _ := strconv.ParseInt(useridstr, 10, 32)
			//	fmt.Println(useridstr)
			//idstr := strconv.Itoa(idint)
			//		id, _ := strconv.ParseInt(idstr, 10, 32)
			if localDebug {
				fmt.Println(deptnum, useridstr, idstr)
			}
			//			fmt.Println(deptnum, userid, id)

			if deptjud[deptnum] == false {
				deptjud[deptnum] = true
				dept1.isExist = true
				dept1.deptnum = deptnum
				dept[flag] = dept1
				flag++
			}

			if localDebug == true {
				fmt.Println(dept1)
			}

			//			fmt.Println("--------------------")
			/*
				i, _ := strconv.ParseInt(id, 10, 32)
				j := int(i)
			*/ /*
				l, _ := strconv.ParseInt(id, 10, 32)
				k := int(id)
			*/
			//		j := int(id)
			var a User
			//		a.id = int(id)
			// a.id = int(id)
			a.id = idstr
			a.userid = useridstr
			a.deptment = deptnum
			//			q := strconv.Itoa(i)
			user[idstr] = a

			if localDebug == true {
				fmt.Println(a)
			}
		}

		if localDebug == true {
			fmt.Println(user)
		}

		//	fmt.Println(user)

	}
	//	fmt.Println(user)

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

				if localDebug == true {
					fmt.Println(checktype, checktime, uid)
				}

				switch checktype {
				case "I":
				case "O":
				default:
					continue
				}
				/*
					i, _ := strconv.ParseInt(uid, 10, 64)
				*/
				// j := uid
				var b Info
				//				b.uid = int(uid)
				b.uid = uid
				b.inout = checktype
				b.date = checktime[0:4] + checktime[5:7] + checktime[8:10]
				b.time = checktime[11:13] + checktime[14:16]
				indexstr := strconv.Itoa(index)
				info[indexstr] = b
				index++
			}

		}

		if localDebug == true {
			fmt.Println(info)
		}

		//fmt.Printf("%s\n", "finish")

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
				if localDebug == true {
					fmt.Println(checktype, checktime, uid)
				}

				switch checktype {
				case "I":
				case "O":
				default:
					continue
				}
				/*
				 i, _ := strconv.ParseInt(uid, 10, 64)
				*/
				// j := uid
				var b Info
				// b.uid = int(uid)
				b.uid = uid
				b.inout = checktype
				b.date = checktime[0:4] + checktime[5:7] + checktime[8:10]
				b.time = checktime[11:13] + checktime[14:16]
				indexstr := strconv.Itoa(index)
				info[indexstr] = b
				index++
			}
		}
		//	fmt.Println(info)

		if localDebug == true {
			fmt.Println(info)
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
		filename = QIANZHUI + info["0"].date + ".csv"

	} else {
		filename = QIANZHUI + time.Now().String()[0:4] + time.Now().String()[5:7] + time.Now().String()[8:10] + ".csv"
	}

	filepath := FilePath + "/" + deptment + "/" + filename

	/*
		file, err := os.Open(filepath)

		if err != nil {
			file, _ = os.Create(filepath)
			file.Close()
		} else {
			file.Close()
		}
	*/
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		writeErrorLog("Create File " + filepath + " Error !")
	}
	file.WriteString("deptnum" + "," + "userid" + "," + "in/out" + "," + "date" + "," + "time" + "\n")

	//fmt.Println(lenth)
	for i := 0; i < lenth; i++ {
		file.WriteString(message[i])
		//	fmt.Println(message[i])
	}

	defer file.Close()
}

//错误日志记录函数
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

//运行日志记录函数
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
