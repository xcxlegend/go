package controllers

import (
	m "github.com/beego/admin/src/models"
)

//LogController Redis数据库控制器
type LogController struct {
	BaseController
}

// Index 列表页
func (this *LogController) Index() {
	page, _ := this.GetInt64("page")
	page_size, _ := this.GetInt64("rows")
	sort := this.GetString("sort")
	order := this.GetString("order")
	if len(order) > 0 {
		if order != "asc" {
			sort = "-" + sort
		}
	} else {
		sort = "-Id"
	}

	if this.IsAjax() {
		rediss, count := m.Getloglist(page, page_size, sort, this.GetString("type"))
		this.Data["json"] = &map[string]interface{}{"total": count, "rows": &rediss}
		this.ServeJSON()
		return
	} else {
		this.Data["types"] = m.LogNodeTitle
		this.TplName = this.GetTemplatetype() + "/log/index.tpl"
	}
}
