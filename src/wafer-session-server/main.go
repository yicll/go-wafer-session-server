package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"wafer-session-server/redis_helper"
	_ "wafer-session-server/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	if beego.BConfig.RunMode == "prod" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	// 根据配置的数据库类型初始化相应的连接池
	if beego.AppConfig.String("dbdriver") == "mysql" {
		RegisterMysqlConnection()
	} else if beego.AppConfig.String("dbdriver") == "redis" {
		host := beego.AppConfig.String("redis.host")
		port, _ := beego.AppConfig.Int("redis.port")
		pwd := beego.AppConfig.String("redis.password")
		db, _ := beego.AppConfig.Int("redis.db")
		mi, _ := beego.AppConfig.Int("redis.maxidle")
		ma, _ := beego.AppConfig.Int("redis.maxactive")
		timeout, _ := beego.AppConfig.Int("redis.timeout")

		redis_helper.NewReidsPool(host, port, pwd, db, mi, ma, time.Duration(timeout)*time.Second)
	} else {
		fmt.Println("server config dbdriver error, [redis/mysql] required")
		os.Exit(1)
	}

	// 设置日志配置信息
	pwd, _ := os.Getwd()
	execDir := flag.String("d", pwd, "execute directory")
	fmt.Println("current execute directory:", *execDir)

	beego.SetLogger("file", fmt.Sprintf(`{"filename":"%s/logs/server.log"}`, *execDir))

	if beego.BConfig.RunMode == "prod" {
		beego.SetLevel(beego.LevelInformational)
		beego.BeeLogger.DelLogger("console")
	}
	beego.SetLogFuncCall(true)

	beego.Run()
}

// 注册mysql连接
func RegisterMysqlConnection() {
	host := beego.AppConfig.String("mysql.host")
	port := beego.AppConfig.String("mysql.port")
	user := beego.AppConfig.String("mysql.user")
	pwd := beego.AppConfig.String("mysql.password")
	db := beego.AppConfig.String("mysql.db")
	debug, _ := beego.AppConfig.Bool("mysql.debug")

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&autocommit=true",
		user, pwd, host, port, db)

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.Debug = debug
	orm.RegisterDataBase("default", "mysql", conn, 30, 1000)
}
