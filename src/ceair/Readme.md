Readme.md

ceair.exe
### 软件简介

使用Go语言开发，使用Go开源库go-odbc (https://gowalker.org/github.com/weigj/go-odbc)提供的数据库驱动。后台数据库为SQLServer
该软件由于数据处理需求而产生，（其实是因为商业软件有点贵，同时也想考验自己的能力而产生的）。

软件的主要作用是从SQLServer数据库中索取需要的字段，并组成csv文档。

### 软件的信息

#### 软件手册

软件ceair.exe的安装目录下为存在一个conf 文件夹，其下存在config.txt配置文件，如果没有特殊需求请不要随意改动。

软件正常运行时不需要管理的，自动运行于后台，尽量少的人为干预。

### 软件错误信息

在软件的安装目录下有一个log文件夹存放日志信息，运行日志和错误日志（默认只运行错误日志），如果需要运行“运行日志信息”请更改配置文件的rtRunLog 的值为true。软件调试测试，可以更改配置文件的 localDebug 的值为 true。软件会实时打印运行信息。

### 软件配置信息

#### 配置文件语法：
以”//”开始的行为注释信息,其他行的内容为配置信息，等号左边为程序内部变量，禁止改动[*改动会造成程序运行异常*]，等号右边为变量的值，可以根据需求改动 [*改动请须知你改动的什么，确保正确*]。
*注：配置文件的等号左右不可以有空格！*

**配置详解:**
- //Windows DSN接口，默认未使用,可改
- DSN=ceair
- //SQLServer服务器地址
- SERVER=127.0.0.1
- //数据库名称
- DATABASE=CEAIR
- //SQLServer 用户名
- USERNAME=******
- //SQLServer 密码
- PASSWORD=*******@!
- //数据文件存放路径
- FilePath=C:/DFS
- //存放文件的前缀
- QIANZHUI=rmsTA_
- //当天所有数据的文件名标识
- ALL_NAME=AllMessge
- //每天更新时间（小时）
- UpdateHour=00
- //每天更新时间（分钟）
- UpdateMin=01
- //程序调试模式开关
- localDebug=false
- //程序运行日志记录开关
- rtRunLog=false

