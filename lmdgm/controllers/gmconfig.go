package controllers

import (
	"github.com/beego/admin/src/models"
	m "github.com/xcxlegend/go/lmdgm/models"
)

//GMConfigController Config数据库控制器
type GMConfigController struct {
	BaseController
}

// Index 列表页
func (this *GMConfigController) Index() {
	page, _ := this.GetInt64("page")
	page_size, _ := this.GetInt64("rows")
	sort := this.GetString("sort")
	order := this.GetString("order")

	if len(order) > 0 {
		if order == "desc" {
			sort = "-" + sort
		}
	} else {
		sort = "Id"
	}

	if this.IsAjax() {
		cfgs, count := m.GetConfigList(page, page_size, sort)
		this.Data["json"] = &map[string]interface{}{"total": count, "rows": &cfgs}
		this.ServeJSON()
		return
	} else {
		this.TplName = this.GetTemplatetype() + "/config/index.tpl"
	}
}

//AddConfig 添加
func (this *GMConfigController) AddConfig() {
	s := m.Config{}
	if err := this.ParseForm(&s); err != nil {
		//handle error
		this.Rsp(false, err.Error())
		return
	}

	id, err := m.AddConfig(&s)
	if err == nil && id > 0 {
		this.DBLogTplData(models.LOGNODE_CONFIG_ADD, DBLOGNODEREMARK_TPL_CONFIG_ADD, &s)
		this.Rsp(true, "Success")
		return
	}
	this.Rsp(false, err.Error())
	return

}

//DelConfig 删除
func (this *GMConfigController) DelConfig() {
	Id, _ := this.GetInt64("Id")
	var old = m.GetConfigById(Id)
	status, err := m.DelConfigById(Id)
	if err == nil && status > 0 {
		this.DBLogTplData(models.LOGNODE_CONFIG_DEL, DBLOGNODEREMARK_TPL_CONFIG_DEL, &old)
		this.Rsp(true, "Success")
		return
	}
	this.Rsp(false, err.Error())
	return
}

//UpdateConfig 更新
func (this *GMConfigController) UpdateConfig() {
	s := m.Config{}
	if err := this.ParseForm(&s); err != nil {
		//handle error
		this.Rsp(false, err.Error())
		return
	}
	var o = m.GetConfigById(s.Id)
	id, err := m.UpdateConfig(&s)
	if err != nil {
		this.Rsp(false, err.Error())
		return
	}
	if id > 0 {
		var log = map[string]interface{}{
			"old":    &o,
			"update": this.Input(),
		}
		this.DBLogTplData(models.LOGNODE_CONFIG_UPDATE, DBLOGNODEREMARK_TPL_CONFIG_UPDATE, log)
		this.Rsp(true, "Success")
		return
	}
	this.Rsp(false, "no update")
	return
}
