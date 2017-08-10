package controllers

import (
	"time"

	"github.com/beego/admin/src/models"
	m "github.com/xcxlegend/go/lmdgm/models"
	"github.com/xcxlegend/go/lmdgm/pb"
)

const TIME_LAYOUT = "2006-01-02 15:04:05"

//AreaController 分区服务区管理
type AreaController struct {
	BaseController
}

// Index 列表页
func (this *AreaController) Index() {
	if this.IsAjax() {
		areas := m.GetAreaAll()
		this.Data["json"] = &map[string]interface{}{"total": len(areas), "rows": &areas}
		this.ServeJSON()
		return
	} else {
		this.TplName = this.GetTemplatetype() + "/area/index.tpl"
	}
}

func (this *AreaController) parseForm(s *pb.LogicServerInfo) {
	var id, _ = this.GetInt32("id", 0)
	s.Id = &id
	var name = this.GetString("name", "")
	s.Name = &name
	var state = this.GetString("state", "")
	s.State = &state

	var start = this.GetString("start", "")
	if start != "" {
		var thetime, _ = time.ParseInLocation(TIME_LAYOUT, start, time.Local)
		var startstamp = thetime.Unix()
		s.Start = &startstamp
	}

	var end = this.GetString("end", "")
	if end != "" {
		var thetime, _ = time.ParseInLocation(TIME_LAYOUT, end, time.Local)
		var startstamp = thetime.Unix()
		s.End = &startstamp
	}
}

//AddArea 添加
func (this *AreaController) AddArea() {
	s := &pb.LogicServerInfo{}
	this.parseForm(s)
	// beego.Debug(s)
	_, err := m.AddArea(s)
	if err == nil {
		this.DBLogTplData(models.LOGNODE_AREA_ADD, DBLOGNODEREMARK_TPL_AREA_ADD, s)
		this.Rsp(true, "Success")
		return
	}
	this.Rsp(false, err.Error())
	return

}

//DelArea 删除 暂时不开放
func (this *AreaController) DelArea() {
	/* 	Id, _ := this.GetInt64("Id")
	   	var old = m.GetAreaById(Id)
	   	status, err := m.DelAreaById(Id)
	   	if err == nil && status > 0 {
	   		this.DBLogTplData(models.LOGNODE_REDIS_DEL, DBLOGNODEREMARK_TPL_REDIS_DEL, &old)
	   		this.Rsp(true, "Success")
	   		return
	   	}
	   	this.Rsp(false, err.Error())
	   	return */
}

//UpdateArea 更新
func (this *AreaController) UpdateArea() {
	s := &pb.LogicServerInfo{}
	this.parseForm(s)
	var o = m.GetAreaById(s.GetId())
	id, err := m.UpdateArea(s)
	if err != nil {
		this.Rsp(false, err.Error())
		return
	}
	if id > 0 {
		var log = map[string]interface{}{
			"old":    o,
			"update": s,
		}
		this.DBLogTplData(models.LOGNODE_AREA_UPDATE, DBLOGNODEREMARK_TPL_AREA_UPDATE, log)
		this.Rsp(true, "Success")
		return
	}
	this.Rsp(false, "no update")
	return
}
