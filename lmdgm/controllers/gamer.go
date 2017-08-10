package controllers

import (
	"github.com/xcxlegend/go/lmdgm/models"
	"github.com/xcxlegend/go/lmdgm/pb"
)

//GamerController 玩家信息查询
type GamerController struct {
	BaseController
}

func (this *GamerController) Index() {
	var AllItems = models.Doc.GamerItemBaseData
	this.Data["items"] = AllItems

	this.Data["players"] = models.Doc.PlayerSdMainData
	//fmt.Println(models.Doc.PlayerSdMainData)
	this.TplName = this.GetTemplatetype() + "/gamer/index.tpl"
}

func (this *GamerController) Search() {
	var id, _ = this.GetInt32("id", 0)
	if id <= 0 {
		//this.Rsp(false, "param error id")
		//return
		this.ResponseJson(nil)
	}
	var c = models.GetRdsClientByRid(0)
	var gamer = models.RDSGetGamerInfo(id, c)
	var players = models.RDSGetPlayers(id, c)
	var goods = models.RDSGetPackGoods(id, c)
	var mails = models.RDSGetGamerMails(id, c)

	var res = map[string]interface{}{
		"gamer":   gamer,
		"players": players,
		"goods":   goods,
		"mails":   mails,
	}
	this.ResponseJson(res)
	return
}

func (this *GamerController) AddPlayer() {
	var gid, _ = this.GetInt32("gid", 0)
	var pid, _ = this.GetInt32("pid", 0)
	if gid == 0 {
		this.Rsp(false, "GID error")
		return
	}
	if _, ok := models.Doc.PlayerSdMainData[pid]; !ok {
		this.Rsp(false, "PID error")
		return
	}
	var c = models.GetRdsClientByRid(0)
	var players = models.RDSGetPlayers(gid, c)
	for _, p := range players {
		if p.Id == pb.Int32(pid) {
			this.Rsp(false, "exist")
			return
		}
	}
	var player = new(pb.Player)
	player.Id = pb.Int32(pid)
	player.Level = pb.Int32(1)
	player.Experience = pb.Int32(0)
	if models.RDSAddPlayers(gid, player, c) {
		this.Rsp(true, "ok")
	} else {
		this.Rsp(false, "exist")
		return
	}
}
