package main

import (
	"database/sql"
	"fmt"
	// "go-odbc"
	_ "go-odbc/driver"
	"strconv"
	// "strings"
	"os"
	"time"
)

const (
	DSN        = "renwu"
	SERVER     = "222.24.24.91"
	USERNAME   = "Admin"
	PASSWORD   = ""
	FilePath   = "C:/QianDao"
	Filename   = "QianDao.csv"
	DATABASE   = "renwu2"
	QIANZHUI   = "rmsTA_"
	ALL_NAME   = "AllMessge"
	UpdateHour = 2
	UpdateMin  = 10
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

var Update bool

func main() {
	for true {
		createHomeDir()

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

func createHomeDir() {
	dir, err := os.Open(FilePath)
	if err != nil {
		err := os.Mkdir(FilePath, os.ModePerm)
		if err != nil {
			fmt.Println("Create Home Dir " + FilePath + " Error")
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
			// fmt.Println("=============================================================================================")
			// fmt.Println(user[info[i].uid].deptment, user[info[i].uid].userid, info[i].inout, info[i].date, info[i].time)
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

		mess[dept[j].deptnum] = m
		write(dept[j].deptnum, m, messkey[dept[j].deptnum])
		fmt.Println("----------------------------")
		fmt.Println(m)
		fmt.Println("----------------------------")
		// fmt.Println("===", dept[j].deptnum, m, messkey[dept[j].deptnum])
	}
	writeAll(m_all)
	if Update {
		Update = false
	}
	/*
		if Update {
			for j := 0; j < deptmaplen; j++ {
				m := make(mesg, infolen)
				for i := 0; i < infolen; i++ {
					//fmt.Println(v.uid)
					fmt.Println(user[info[i].uid].deptment, user[info[i].uid].userid, info[i].inout, info[i].date, info[i].time)
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
					//	write(dept[j].deptnum, m, messkey[dept[j].deptnum])

				}

				mess[dept[j].deptnum] = m
				write(dept[j].deptnum, m, messkey[dept[j].deptnum])
				writeAll(m_all)
			}
		}
	*/
}

func writeAll(mess []string) {
	dirname := FilePath + "/" + "AllMessage"
	dir, err := os.Open(dirname)
	if err != nil {
		err := os.Mkdir(dirname, os.ModePerm)
		if err != nil {
			fmt.Println("Create Dir " + dirname + " Error!")
		}
	}
	dir.Close()
	filename := FilePath + "/" + "AllMessage" + "/" + QIANZHUI + ALL_NAME + ".cvs"
	file, err := os.Open(filename)
	if err != nil {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Create File " + filename + " Error!")
		}
		file.Close()

	}
	file.Close()

	file, err = os.OpenFile(filename, os.O_CREATE, os.ModePerm)

	for _, v := range mess {
		file.WriteString(v)
	}
	defer file.Close()
}

//从数据库获取数据，并进行相应的处理
func getData() {
	conn, err := sql.Open("odbc", "driver={SQL Server};DSN="+DSN+";SERVER="+SERVER+";Database="+DATABASE+";UID="+USERNAME+";PWD="+PASSWORD) //

	if err != nil {
		fmt.Println("Connecting Error")
		return
	}
	// fmt.Println(conn)
	defer conn.Close()
	stmt, err := conn.Prepare("SELECT deptname,badgenumber,userid from USERINFO left join DEPARTMENTS ON USERINFO.DEFAULTDEPTID = DEPARTMENTS.DEPTID")
	if err != nil {
		fmt.Println("Query Error", err, stmt)
		return
	}
	// fmt.Println(stmt)
	defer stmt.Close()

	row, err := stmt.Query()
	if err != nil {
		fmt.Println("Query Error", err)
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
			fmt.Println(deptnum, id, uid)

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
	}

	if time.Now().Hour() == UpdateHour && time.Now().Minute() == UpdateMin {
		Update = true
	}
	if Update {
		stmt, err = conn.Prepare("SELECT checktype,checktime,userid from CHECKINOUT where DateDiff(dd,CHECKTIME,getdate())=1")
		if err != nil {
			fmt.Println("Query Error", err, stmt)
			return
		}
		// fmt.Println(stmt)
		defer stmt.Close()

		row, err = stmt.Query()
		if err != nil {
			fmt.Println("Query Error", err)
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
				fmt.Println(checktype, checktime, uid)
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
		fmt.Printf("%s\n", "finish")
		return
	} else {
		// stmt, err = conn.Prepare("SELECT checktype,checktime,userid from CHECKINOUT ")
		stmt, err = conn.Prepare("SELECT checktype,checktime,userid from CHECKINOUT where DateDiff(dd,CHECKTIME,getdate())=0")
		if err != nil {
			fmt.Println("Query Error", err, stmt)
			return
		}
		// fmt.Println(stmt)
		defer stmt.Close()

		row, err = stmt.Query()
		if err != nil {
			fmt.Println("Query Error", err)
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
				fmt.Println(checktype, checktime, uid)
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
		fmt.Printf("%s\n", "finish====")
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
				fmt.Println("Mkdir " + FilePath + dept[i].deptnum + " Error")

			}
			fmt.Println(dept[i], err)
		}

		fmt.Println(dept[i], err)
		dir.Close()

	}
	fmt.Println(dept[1])
}

//创建文件，函数
func createFile(path, name string) {
	filepath := path + name
	file, err := os.Open(filepath)
	if err != nil {
		_, err = os.Create(filepath)
		if err != nil {
			fmt.Println("Create " + filepath + " Error")
		}
	}
	defer file.Close()
}

//写入相应的文件
func write(deptment string, message []string, lenth int) {
	var filename string
	if Update {
		filename = QIANZHUI + info[0].date + ".csv"
		fmt.Println(filename)
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
		fmt.Println(message[i])
	}

	file.Close()
}

func UpdateFile() {

}
