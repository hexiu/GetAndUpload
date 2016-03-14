package main
import (
  "fmt"
  ftp "github.com/jlaffaye/ftp"
  "log"
  "net"
  "os"
  "path/filepath"
  "strconv"
  "strings"
  "time"
)
func main() {
  fmt.Println("Hello World!")
  if len(os.Args) != 2 {
    log.Fatal("Usage:" + filepath.Base(os.Args[0]) + " log_dir ")
    os.Exit(1)
  }
  //logFileName是将要分析的日志
  logFileName, _, _ := getLogFileName()
  serverIp := getLocalIpAddr()
  //serverName, _ := os.Hostname()
  time.Sleep(time.Duration(90) * time.Second)
  dir := os.Args[1]
  filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
    if f == nil {
      return err
    }
    if f.IsDir() {
      return nil
    }
    if f.Name() == logFileName {
      fmt.Println(path)
      //pathFields的作用是将日志path解析成一个数据，从而可以得到日志的域名,注意，如果是linux系统，要用“/”
      pathFields := strings.Split(path, "\")
      var domainName string
      if len(pathFields) > 3 {
        domainName = pathFields[len(pathFields)-3]
      }
      fmt.Println(time.Now())
      ftpUploadFile("ftp-server-ip:21", "logftpuser", "ftp-password", path, domainName, serverIp+"_"+logFileName)
      fmt.Println(time.Now())
    }
    return nil
  })
}
func getLogFileName() (string, string, string) {
  MonthTOstr := map[string]string{"January": "01",
    "February": "02",
    "March":   "03",
    "April":   "04",
    "May":    "05",
    "June":   "06",
    "July":   "07",
    "August":  "08",
    "September": "09",
    "October":  "10",
    "November": "11",
    "December": "12"}
  timenow := time.Now()
  year, month, day := timenow.Date()
  //monthStr := month.String()
  hour, _, _ := timenow.Clock()
  yearStr := strings.TrimLeft(strconv.Itoa(year), "20") //去掉前面的四位数年份如"2014"年的“20”
  dayStr, hourStr := strconv.Itoa(day), strconv.Itoa(hour)
  if day < 10 {
    dayStr = "0" + dayStr
  }
  if hour < 10 {
    hourStr = "0" + hourStr
  }
  fileName := "ex" + yearStr + MonthTOstr[month.String()] + dayStr + hourStr + ".log"
  logDay := yearStr + MonthTOstr[month.String()] + dayStr
  logMonth := yearStr + MonthTOstr[month.String()]
  //monthSrt := strconv.Itoa(timenow.Month())
  //fmt.Println(fileName, logDay)
  return fileName, logDay, logMonth
  //fmt.Println(fileName)
}
func getLocalIpAddr() string {
  //这里使用一个合法的IP就行了，端口随便，即使没有打开也行，也许因为使用UDP，如果用TCP的话，对端不打开就会有问题
  conn, err := net.Dial("udp", "192.168.8.51:80")
  if err != nil {
    //fmt.Println(err.Error())
    return "127.0.0.1"
  }
  defer conn.Close()
  //fmt.Println(conn.LocalAddr().String())
  //conn.
  //fmt.Println(strings.Split(conn.LocalAddr().String(), ":")[0])
  return strings.Split(conn.LocalAddr().String(), ":")[0]
}
func ftpUploadFile(ftpserver, ftpuser, pw, localFile, remoteSavePath, saveName string) {
  ftp, err := ftp.Connect(ftpserver)
  if err != nil {
    fmt.Println(err)
  }
  err = ftp.Login(ftpuser, pw)
  if err != nil {
    fmt.Println(err)
  }
  //注意是 pub/log，不能带“/”开头
  ftp.ChangeDir("pub/log")
  dir, err := ftp.CurrentDir()
  fmt.Println(dir)
  ftp.MakeDir(remoteSavePath)
  ftp.ChangeDir(remoteSavePath)
  dir, _ = ftp.CurrentDir()
  fmt.Println(dir)
  file, err := os.Open(localFile)
  if err != nil {
    fmt.Println(err)
  }
  defer file.Close()
  err = ftp.Stor(saveName, file)
  if err != nil {
    fmt.Println(err)
  }
  ftp.Logout()
  ftp.Quit()
  fmt.Println("success upload file:", localFile)
}