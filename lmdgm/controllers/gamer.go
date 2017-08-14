package controllers

import (
	"encoding/json"
	"fmt"
	bm "github.com/beego/admin/src/models"
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
	var c = models.GetRdsClientByRid(0)
	if id <= 0 {
		//this.Rsp(false, "param error id")
		//return

		id = models.RDSGetGIDByName(this.GetString("id"), c)
		if id <= 0 {
			this.ResponseJson(nil)
		}
	}

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
	var name string
	if p, ok := models.Doc.PlayerSdMainData[pid]; !ok {
		this.Rsp(false, "PID error")
		return
	} else {
		name = *p.Name
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
		//this.DBLogTplData(bm.LOGNODE_GAMER_ADD_PLAYER, DBLOGNODEREMARK_TPL_GAMER_ADD_PLAYER, fmt.Sprintf("%s(%v)", name, pid))
		this.DBLogTpl(bm.LOGNODE_GAMER_ADD_PLAYER, DBLOGNODEREMARK_TPL_GAMER_ADD_PLAYER, gid, fmt.Sprintf("%s(%v)", name, pid))
		this.Rsp(true, "ok")
	} else {
		this.Rsp(false, "exist")
		return
	}
}

func (this *GamerController) UpdatePlayer() {
	var gid, _ = this.GetInt32("gid", 0)
	var pid, _ = this.GetInt32("id", 0)
	if gid == 0 || pid == 0 {
		this.Rsp(false, "error")
		return
	}
	var c = models.GetRdsClientByRid(0)
	var player = models.RDSGetGamerPlayer(gid, pid, c)
	if player == nil {
		this.Rsp(false, "no player")
		return
	}
	var level, _ = this.GetInt32("level")
	var maxPlayerLv int32
	if vMain, ok := models.Doc.PlayerSdMainData[pid]; ok {
		for _, vLv := range models.Doc.PlayerLvUpData {
			if vLv.GetPlayerLvUpId() == vMain.GetPlayerLvUpId() {
				if vLv.GetPlayerLv() > maxPlayerLv {
					maxPlayerLv = vLv.GetPlayerLv()
				}
			}
		}
	}
	if maxPlayerLv < level {
		this.Rsp(false, fmt.Sprintf("max: %v", maxPlayerLv))
		return
	}
	var oplayer = &pb.Player{
		Id:         player.Id,
		Level:      player.Level,
		Experience: player.Experience,
	}
	player.Level = pb.Int32(level)
	if models.RDSUpdatePlayer(gid, player, c) {
		var log, _ = json.MarshalIndent(map[string]interface{}{
			"old":    oplayer,
			"update": player,
		}, "", " ")
		this.DBLogTpl(bm.LOGNODE_GAMER_UPDATE_PLAYER, DBLOGNODEREMARK_TPL_GAMER_UPDATE_PLAYER, gid, string(log))
		this.Rsp(true, "ok")
	} else {
		this.Rsp(false, "fail")
	}
}
