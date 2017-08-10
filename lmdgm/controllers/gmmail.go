package controllers

import (
	"antnet"
	//"fmt"
	"github.com/xcxlegend/go/lib"
	"github.com/xcxlegend/go/lmdgm/models"
	"github.com/xcxlegend/go/lmdgm/pb"
	"strings"
	"time"
)

//GMMailController 分区服务区管理
type GMMailController struct {
	BaseController
}

func (this *GMMailController) Index() {
	if this.IsAjax() {
		mails := models.RDSGetAllSysMail()
		this.Data["json"] = &map[string]interface{}{"total": len(mails), "rows": &mails}
		this.ServeJSON()
		return
	} else {
		var AllItems = models.Doc.GamerItemBaseData
		this.Data["items"] = AllItems
		//fmt.Println(models.Doc.PlayerSdMainData)
		this.TplName = this.GetTemplatetype() + "/gmmail/index.tpl"
	}
}

func (this *GMMailController) parseAttach(mail *pb.Mail) bool {
	var attaches_str = strings.TrimSpace(this.GetString("attaches", ""))
	var number_str = strings.TrimSpace(this.GetString("attaches_num", ""))
	if attaches_str == "" {
		if number_str == "" {
			return true
		} else {
			return false
		}
	}
	var attaches = strings.Split(attaches_str, ",")
	var numbers = strings.Split(number_str, ",")
	if len(attaches) != len(numbers) {
		return false
	}

	mail.Attachments = []*pb.ItemConfig{}
	for k, att := range attaches {
		var attid = antnet.Atoi(att)
		if attid == 0 {
			return false
		}
		var number = antnet.Atoi(numbers[k])
		if number == 0 {
			return false
		}

		var attachment = &pb.ItemConfig{
			Id:     pb.Int32(int32(attid)),
			Number: pb.Int32(int32(number)),
		}
		mail.Attachments = append(mail.Attachments, attachment)
	}
	return true
}

func (this *GMMailController) SendToGamer() {

	var gid, _ = this.GetInt32("gid", 0)
	if gid == 0 {
		this.Rsp(false, "GID错误")
		return
	}

	var theme = strings.TrimSpace(this.GetString("theme", ""))
	if theme == "" {
		this.Rsp(false, "Theme不能为空")
		return
	}

	var msg = strings.TrimSpace(this.GetString("msg", ""))
	if msg == "" {
		this.Rsp(false, "Msg不能为空")
		return
	}

	var mail = new(pb.Mail)
	if !this.parseAttach(mail) {
		this.Rsp(false, "附件格式错误")
		return
	}

	mail.Id = pb.Int64(lib.GetUUID(0))
	mail.SenderId = pb.Int32(0)
	mail.Msg = pb.String(msg)
	mail.Theme = pb.String(theme)
	mail.Time = pb.Int64(time.Now().Unix())
	mail.State = pb.Int32(0)
	mail.AttachmentState = pb.Int32(0)

	//fmt.Println(mail)

	if models.RDSSendMailToGamer(gid, mail, nil) {
		this.Rsp(true, "ok")
	} else {
		this.Rsp(false, "err")
	}
}
