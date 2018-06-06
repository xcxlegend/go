package controllers

import (
	"strings"

	"github.com/astaxie/beego"
)

//InnerBaseController 内网接口
type InnerBaseController struct {
	beego.Controller
}

func (this *InnerBaseController) Prepare() {
	var host = this.Ctx.Input.Host()
	var inner = beego.AppConfig.String("server.innernet.root")
	if inner != "" && !strings.HasPrefix(host, inner) {
		this.Abort("401")
	}
}
func (this *InnerBaseController) ResponseJson(res interface{}) {
	this.Data["json"] = res
	this.ServeJSON()
}
