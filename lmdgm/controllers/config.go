package controllers

import (
	"github.com/beego/admin/src/rbac"
)

type ConfigController struct {
	rbac.CommonController
}

func (c *ConfigController) GetList() {
	c.Ctx.Output.Body([]byte("hello"))
}
