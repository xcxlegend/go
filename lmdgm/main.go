package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/xcxlegend/go/lmdgm/models"
	_ "github.com/xcxlegend/go/lmdgm/routers"
	_ "github.com/xcxlegend/go/lmdgm/servers"
)

const VERSION = "1.2.0"

func main() {
	// if beego.BConfig.RunMode == "dev"  {
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	// }

	test()
	models.AddDocConfig()
	models.SetConfigData(3, 4, "doc/")
	models.SetFunc()
	beego.Run()

}

func test() {
	fmt.Println("go a test")
}
