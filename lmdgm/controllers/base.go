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
