package controllers

import (
	// "github.com/astaxie/beego"
	"github.com/beego/admin/src/rbac"
)

type BaseController struct {
	rbac.CommonController
}

func (this *BaseController) ResponseJson(res interface{}) {
	this.Data["json"] = res
	this.ServeJSON()
}

func (this *BaseController) Rsp(status bool, info string) {
	var res = map[string]interface{}{
		"status": status,
		"info":   info,
	}
	this.ResponseJson(res)
}
