package main

import (
	"antnet"
	"common"
	"events"
	"flag"
	"fmt"
	"models"
	"routers"
	_ "routers"
	_ "servers"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	//"github.com/astaxie/beego/orm"
	//"github.com/golang/protobuf/proto"
)

const VERSION = "1.8.0"

func main() {

	beego.SetLogger("file", `{"filename":"logs/log.log"}`)

	// if beego.BConfig.RunMode == "dev"  {
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	// }
	beego.SetStaticPath("/static", "static")
	models.InitDocConfig()
	// 初始化统计相关的

	daemon := flag.Int("d", 0, "1 for daemon")
	sql_exc := flag.String("sql", "", "path for exc sql")
	buildJSTable := flag.String("build_js", "", "build table for js")

	flag.Parse()

	if *sql_exc != "" {
		beego.Debug("exc sql:", *sql_exc)
		f := "sql/" + *sql_exc
		if antnet.PathExists(f) {
			sqls, _ := antnet.ReadFile(f)
			for _, sql := range strings.Split(string(sqls), ";") {
				if sql != "" {
					r, err := orm.NewOrm().Raw(sql).Exec()
					beego.Debug("exec sql res:", f, sql, r, err)
				}
			}
		}
		return
	}

	if *buildJSTable == "1" {
		events.BuildJSTable()
		return
	}

	if *daemon == 1 {
		antnet.Daemon("-d")
		antnet.Println("go daemon")
		return
	}

	beego.Debug("系统正在初始化...")
	if beego.AppConfig.String("init_statis") == "1" {
		err := events.InitStatis()
		if err != nil {
			beego.Critical("init statistics err:", err)
			panic(fmt.Sprintf("init statistics err: %v", err))
		}
	}

	// 初始化服务器管理器
	events.InitServerManager()

	events.CronStart()

	// models.InitDBDataSource()
	models.InitFileDataSource()
	routers.InitFileRouter()

	events.Init()
	beego.Debug("系统初始化完成. 正在启动web服务...")

	common.CopyDir("./json/sdk", "/Users/a1234/workspace/sync/conf/20180502144409/conf/sdk")

	beego.Run()

}
