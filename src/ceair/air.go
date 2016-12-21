package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "go-odbc/driver"
	"io"
	"os"
	"strconv"
	"time"
	 "code.google.com/p/mahonia"
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
	badgenumber string 
	deptname string
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
//var Update bool

// 配置信息全局变量（有默认值）
var (
	DSN         string = "ceair"
	SERVER      string = "192.168.230.241"
	USERNAME    string = "sa"
	PASSWORD    string = "123456"
	FilePath    string = "C:/QianDao"
	FilePathBak string = "C:/CLOCKING"
	DATABASE    string = "CEAIR"
	QIANZHUI    string = "rmsTA_"
	ALL_NAME    string = "AllMessge"
	// UpdateHour  int    = 00
	// UpdateMin   int    = 01
	UpdateTime string = "20000"
	localDebug bool   = false
	rtRunLog   bool   = false
)

// 主函数
func main() {
	//localDebug = false
	// 创建配置文件目录
	readConfDir()
	// 加载配置文件
	initConf()

	// 创建签到数据存放主目录
	createHomeDir()

	//获取数据
	getData()

	// 数据处理引擎，并且写入文档
		// message()

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
	if Conf["FilePath"] != "" {
		FilePath = Conf["FilePath"]
	}
	if Conf["rtRunLog"] != "" {
		jud := "true"
		if jud == Conf["rtRunLog"] {
			rtRunLog = true
		} else {
			rtRunLog = false
		}

	}
	if Conf["localDebug"] != "" {
		jud := "true"
		if jud == Conf["localDebug"] {
			localDebug = true
		} else {
			localDebug = false
		}
	}

	if Conf["FilePathBak"] != "" {
		FilePathBak = Conf["FilePathBak"]
	}
	if Conf["UpdateTime"] != "" {
		UpdateTime = Conf["UpdateTime"]
	}

	if localDebug {
		fmt.Println(Conf,UpdateTime)
	}
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

		err := os.Mkdir(dirName, os.ModePerm)
		if err != nil {
			writeErrorLog(err.Error() + "  Create Conf Dir " + FilePath + " Error")
		}

	}
	defer dir.Close()

	fileName := dirName + "/" + "config.txt"
	file, err := os.Open(fileName)
	if err != nil {
		writeErrorLog("No Config File.")
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

	// fmt.Println("dept", dept, info)

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

			// fmt.Println("................................", user[info[i].uid].deptment)

			if user[info[i].uid].deptment == dept[j].deptnum {

				m[messkey[dept[j].deptnum]] = message

				messkey[dept[j].deptnum]++
			} else {
				continue
			}

		}
		// fmt.Println("messkey[dept[j].deptnum]", messkey[dept[j].deptnum])
		mess[dept[j].deptnum] = m
		write(dept[j].deptnum, m, messkey[dept[j].deptnum])
		if rtRunLog {
			writeRunLog("Write Department " + dept[j].deptnum + " success.")
		}
	}
	if rtRunLog {
		writeRunLog("Write Department  success.")
	}
	//writeAll(m_all)
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

	var filename string
	// filename = FilePath + "/" + "AllMessage" + "/" + QIANZHUI + time.Now().String()[0:4] + time.Now().String()[5:7] + time.Now().String()[8:10] + "_" + time.Now().String()[11:13] + time.Now().String()[14:16] + "_" + ALL_NAME + ".csv"
	filename = FilePath + "/" + "AllMessage" + "/" + QIANZHUI + time.Now().String()[0:4] + time.Now().String()[5:7] + time.Now().String()[8:10] + "_" + ALL_NAME + ".csv"

	file, err := os.Open(filename)
	if err != nil {
		file, err := os.Create(filename)
		if err != nil {
			writeErrorLog("Create File " + filename + " Error!")
		}

		file.Close()

	}
	file.Close()

	file, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		writeErrorLog(err.Error())
	}
	// file.WriteString("deptnum" + "," + "userid" + "," + "in/out" + "," + "date" + "," + "time" + "\n")
	for _, v := range mess {
		file.WriteString(v)
	}
	if rtRunLog {
		writeRunLog("Write Day AllMessage Success.")
	}
	defer file.Close()
}

//从数据库获取数据，并进行相应的处理
func getData() {

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

	stmt, err := conn.Prepare("SELECT deptname from  DEPARTMENTS")
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
	fmt.Println(row)
	defer row.Close()

	//dept flag
	flag := 0

	var deptnumber map[int]string = make(map[int]string,0)

	for row.Next() {
		var deptname string
		if err := row.Scan(&deptname); err == nil {
			if localDebug {
				fmt.Println(deptname)
			}
			//			fmt.Println(deptnum, userid, id)
			deptname = gbk_utf_8(deptname)
			deptnumber[flag]=deptname
			flag=flag+1
			createPartDir(deptname)
	
		}
	}

	fmt.Println("test")

	// stmt, err = conn.Prepare("SELECT checktype,checktime,userid from CHECKINOUT ")
	for _,deptname := range deptnumber {
		stmt, err = conn.Prepare("select checktype,checktime,qiandaouser.userid,deptname,badgenumber from qiandaouser,kaoqin where qiandaouser.userid = kaoqin.userid and deptname=" + "'"+deptname+"'")
		if err != nil {
			writeErrorLog("Query Error" + err.Error())
			return
		}
		fmt.Println(stmt,deptname)
		defer stmt.Close()

		row, err = stmt.Query()

		if err != nil {
			writeErrorLog("Query Error" + err.Error())
			return
		}
	// fmt.Println(row)
		defer row.Close()

		var message []string=make([]string,0)
		for row.Next() {
			var checktype string
			var checktime string
			var uid string
			var deptname string 
			var badgenumber string 
			if err := row.Scan(&checktype, &checktime, &uid,&deptname,&badgenumber); err == nil {
				if localDebug == true {
					fmt.Println(checktype, checktime, uid,deptname,badgenumber)
				}
				checktype=gbk_utf_8(checktype)
				checktime=gbk_utf_8(checktime)
				uid=gbk_utf_8(uid)
				deptname=gbk_utf_8(deptname)
				badgenumber=gbk_utf_8(badgenumber)

				// fmt.Println("*****************************8")
				var checktypei string
				switch checktype {
				case "I":
					checktypei = "In"
				case "O":
					checktypei = "Out"
				default:
					continue
				}
				var b Info
				b.uid = uid
				b.inout = checktypei
				b.date = checktime[0:4] + checktime[5:7] + checktime[8:10]
				b.time = checktime[11:13] + checktime[14:16]
				b.badgenumber = badgenumber
				b.deptname = deptname

				mess:=b.deptname+","+b.uid+","+b.inout+","+b.date+","+b.time+ "\n"
				message=append(message,mess)
				fmt.Println(mess)
			}

			if localDebug == true {
				fmt.Println(info)
			}
		}
		length:=len(message)
		write(deptname, message, length)
	}
}

//创建部门目录
func createPartDir(deptname string) {
		dir, err := os.Open(FilePath + "/" + deptname)
		if err != nil {
			err := os.MkdirAll(FilePath+"/"+deptname+"/"+"READ", os.ModePerm)
			if err != nil {
				writeErrorLog("Mkdir " + FilePath + deptname + " Error")
			}
		}
		dir.Close()
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
	// var filenamebak string
	filename = QIANZHUI + time.Now().String()[0:4] + time.Now().String()[5:7] + time.Now().String()[8:10] + ".csv"
	// filenamebak = QIANZHUI + time.Now().String()[0:4] + time.Now().String()[5:7] + time.Now().String()[8:10] + "_" + time.Now().String()[11:13] + time.Now().String()[14:16] + "_" + deptment + ".csv"
	if deptment == "" {
		return
	}
	filepath := FilePath + "/" + deptment + "/" + "READ" + "/" + filename
	// filepathbak := FilePathBak + "/" + deptment + "/" + "READ" + "/" + filenamebak

	file, err := os.OpenFile(filepath,os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		writeErrorLog("Create File " + filepath + " Error !")
	}
	// filebak, err := os.OpenFile(filepathbak, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	// if err != nil {
	// 	writeErrorLog("Create File " + filepathbak + " Error !")
	// }

	for i := 0; i < lenth; i++ {
		// fmt.Println(message[i], "*************")
		file.WriteString(message[i])
		// filebak.WriteString(message[i])
	}

	defer file.Close()
	// defer filebak.Close()
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

func gbk_utf_8(gbkstr string) (utfstr string) {
	
     var dec mahonia.Decoder
    dec = mahonia.NewDecoder("gbk")

    if ret, ok := dec.ConvertStringOK(gbkstr); ok {
    	return ret
    }
    return ""
}
