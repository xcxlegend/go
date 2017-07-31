package controllers

import (
	"github.com/beego/admin/src/models"
	m "github.com/xcxlegend/go/lmdgm/models"
)

//GMRedisController Redis数据库控制器
type GMRedisController struct {
	BaseController
}

// Index 列表页
func (this *GMRedisController) Index() {
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
		rediss, count := m.GetRedisList(page, page_size, sort)
		this.Data["json"] = &map[string]interface{}{"total": count, "rows": &rediss}
		this.ServeJSON()
		return
	} else {
		// tree := this.GetTree()
		// this.Data["tree"] = &tree
		// this.Data["Rediss"] = &Rediss
		// if this.GetTemplatetype() != "easyui" {
		// 	this.Layout = this.GetTemplatetype() + "/public/layout.tpl"
		// }
		this.TplName = this.GetTemplatetype() + "/redis/index.tpl"
	}
}

//AddRedis 添加
func (this *GMRedisController) AddRedis() {
	s := m.Redis{}
	if err := this.ParseForm(&s); err != nil {
		//handle error
		this.Rsp(false, err.Error())
		return
	}
	if s.Port == 0 {
		s.Port = 6379
	}
	id, err := m.AddRedis(&s)
	if err == nil && id > 0 {
		this.DBLogTplData(models.LOGNODE_REDIS_ADD, DBLOGNODEREMARK_TPL_REDIS_ADD, &s)
		this.Rsp(true, "Success")
		return
	}
	this.Rsp(false, err.Error())
	return

}

//DelRedis 删除
func (this *GMRedisController) DelRedis() {
	Id, _ := this.GetInt64("Id")
	var old = m.GetRedisById(Id)
	status, err := m.DelRedisById(Id)
	if err == nil && status > 0 {
		this.DBLogTplData(models.LOGNODE_REDIS_DEL, DBLOGNODEREMARK_TPL_REDIS_DEL, &old)
		this.Rsp(true, "Success")
		return
	}
	this.Rsp(false, err.Error())
	return
}

//UpdateRedis 更新
func (this *GMRedisController) UpdateRedis() {
	s := m.Redis{}
	if err := this.ParseForm(&s); err != nil {
		//handle error
		this.Rsp(false, err.Error())
		return
	}
	var o = m.GetRedisById(s.Id)
	id, err := m.UpdateRedis(&s)
	if err != nil {
		this.Rsp(false, err.Error())
		return
	}
	if id > 0 {
		var log = map[string]interface{}{
			"old":    &o,
			"update": this.Input(),
		}
		this.DBLogTplData(models.LOGNODE_REDIS_DEL, DBLOGNODEREMARK_TPL_REDIS_DEL, log)
		this.Rsp(true, "Success")
		return
	}
	this.Rsp(false, "no update")
	return
}
