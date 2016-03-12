package main

/*
import (
	"fmt"
	"os"
	"time"

	// "github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/lunny/godbc"
)

type NxServerState struct {
	ID             int       `xorm:"pk not null 'ID'"`
	GameID         int       `xorm:"not null 'GameID'"`
	IssuerId       int       `xorm:"not null IssuerId"`
	ServerID       int       `xorm:"not null ServerID"`
	ServerName     string    `xorm:"ServerName"`
	OnlineNum      int       `xorm:"not null OnlineNum"`
	MaxOnlineNum   int       `xorm:"not null MaxOnlineNum"`
	ServerIP       string    `xorm:"not null ServerIP"`
	Port           int       `xorm:"not null Port"`
	IsRuning       int       `xorm:"not null IsRuning"`
	ServerStyle    int       `xorm:"ServerStyle"`
	IsStartIPWhile int       `xorm:"not null IsStartIPWhile"`
	LogTime        time.Time `xorm:"IsStartIPWhile"`
	UpdateTime     time.Time `xorm:"UpdateTime"`
	OrderBy        int       `xorm:"not null OrderBy"`
}

func main() {
	File, _ := os.Create("result")
	defer File.Close()
	Engine, err := xorm.NewEngine("odbc", "driver={SQL Server};SERVER=222.24.24.91;Database=renwu2;UID=Admin;PWD=")
	if err != nil {
		fmt.Println("新建引擎", err)
		return
	}
	if err := Engine.Ping(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("|haah")

	Engine.SetTableMapper(core.SameMapper{})

	Engine.ShowSQL = true

	Engine.SetMaxConns(5)

	Engine.SetMaxIdleConns(5)

	result := new(NxServerState)

	lines, _ := Engine.Rows(result)
	defer lines.Close()
	lines.Next()
	r := new(NxServerState)
	for {
		err = lines.Scan(r)
		if err != nil {
			return
		}
		fmt.Println(*r)
		File.WriteString(fmt.Sprintln(*r))
		if !lines.Next() {
			break
		}
	}

}

*/
/*
import (
	"odbc"
)

func main() {
	conn, _ := odbc.Connect("DSN=renwu;UID=user;PWD=password")
	stmt, _ := conn.Prepare("select * from user where username = ?")
	stmt.Execute("admin")
	rows, _ := stmt.FetchAll()
	for i, row := range rows {
		println(i, row)
	}
	stmt.Close()
	conn.Close()
}
*/

import (
	"database/sql"
	"fmt"
	// "go-odbc"
	_ "go-odbc/driver"
)

func main() {

	conn, err := sql.Open("odbc", "driver={SQL Server};DSN=renwu;SERVER=222.24.24.91;Database=renwu2;UID=Admin;PWD=") //

	if err != nil {
		fmt.Println("Connecting Error")
		return
	}
	fmt.Println(conn)
	defer conn.Close()
	stmt, err := conn.Prepare("SELECT checktime from CHECKINOUT")
	if err != nil {
		fmt.Println("Query Error", err, stmt)
		return
	}
	fmt.Println(stmt)
	defer stmt.Close()

	row, err := stmt.Query()
	if err != nil {
		fmt.Println("Query Error", err)
		return
	}
	fmt.Println(row)
	defer row.Close()

	for row.Next() {
		var name string
		if err := row.Scan(&name); err == nil {
			fmt.Println(name, "test")
		}
	}
	fmt.Printf("%s\n", "finish")
	return

}
