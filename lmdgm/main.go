package main

import (
	"github.com/astaxie/beego"
	_ "github.com/xcxlegend/go/lmdgm/models"
	_ "github.com/xcxlegend/go/lmdgm/routers"
	_ "github.com/xcxlegend/go/lmdgm/servers"
)

func main() {
	const VERSION = "1.1.0"
	// if beego.BConfig.RunMode == "dev"  {
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	// }
	beego.Run()
}
