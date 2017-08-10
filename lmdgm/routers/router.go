// @APIVersion 1.0.0
// @Title mobile API
// @Description lmd game gm background
// @Contact xcx_legender@qq.com
package routers

import (
	"github.com/astaxie/beego"
	"github.com/beego/admin"
	"github.com/xcxlegend/go/lmdgm/controllers"
)

func init() {
	admin.Run()
	beego.Router("/", &controllers.MainController{})
	beego.Router("/rbac/config", &controllers.ConfigController{}, "*:GetList")
	beego.Router("/rbac/upload/index", &controllers.UploadController{}, "*:Index")
	beego.Router("/rbac/upload/dir", &controllers.UploadController{}, "*:Dir")
	beego.Router("/rbac/upload/post", &controllers.UploadController{}, "*:Upload")
	beego.Router("/rbac/upload/down", &controllers.UploadController{}, "*:Down")

	// 服务器管理
	beego.Router("/rbac/servers/add", &controllers.ServersController{}, "*:AddServer")
	beego.Router("/rbac/servers/update", &controllers.ServersController{}, "*:UpdateServer")
	beego.Router("/rbac/servers/del", &controllers.ServersController{}, "*:DelServer")
	beego.Router("/rbac/servers/index", &controllers.ServersController{}, "*:Index")
	beego.Router("/rbac/servers/ajax_get_stat", &controllers.ServersController{}, "*:GetStat")
	beego.Router("/rbac/servers/ssh/close", &controllers.ServersController{}, "*:SSHClosePid")
	beego.Router("/rbac/servers/ssh/start", &controllers.ServersController{}, "*:SSHStartApp")
	beego.Router("/rbac/servers/ssh/mount", &controllers.ServersController{}, "*:SSHMount")
	beego.Router("/rbac/servers/conf_content", &controllers.ServersController{}, "*:GetConfContent")
	beego.Router("/rbac/servers/update_conf", &controllers.ServersController{}, "*:UpdateConfContent")

	beego.Router("/rbac/servers/terminal", &controllers.ServersController{}, "*:Terminal")

	beego.Router("/rbac/ssh/ws", &controllers.WSController{}, "*:Get")

	beego.Router("/rbac/dir/index", &controllers.DirController{}, "*:Index")
	beego.Router("/rbac/dir/conf", &controllers.DirController{}, "*:Conf")
	beego.Router("/rbac/dir/upload_conf", &controllers.DirController{}, "*:UploadConf")
	beego.Router("/rbac/dir/app", &controllers.DirController{}, "*:App")
	beego.Router("/rbac/dir/upload_app", &controllers.DirController{}, "*:UploadApp")
	beego.Router("/rbac/dir/file_content", &controllers.DirController{}, "*:FileContent")

	// 数据库管理
	beego.Router("/rbac/redis/add", &controllers.GMRedisController{}, "*:AddRedis")
	beego.Router("/rbac/redis/update", &controllers.GMRedisController{}, "*:UpdateRedis")
	beego.Router("/rbac/redis/del", &controllers.GMRedisController{}, "*:DelRedis")
	beego.Router("/rbac/redis/index", &controllers.GMRedisController{}, "*:Index")

	// 区服管理
	beego.Router("/rbac/area/add", &controllers.AreaController{}, "*:AddArea")
	beego.Router("/rbac/area/update", &controllers.AreaController{}, "*:UpdateArea")
	// beego.Router("/rbac/area/del", &controllers.AreaController{}, "*:DelArea")
	beego.Router("/rbac/area/index", &controllers.AreaController{}, "*:Index")

	// 配置管理
	beego.Router("/rbac/config/add", &controllers.GMConfigController{}, "*:AddConfig")
	beego.Router("/rbac/config/update", &controllers.GMConfigController{}, "*:UpdateConfig")
	beego.Router("/rbac/config/del", &controllers.GMConfigController{}, "*:DelConfig")
	beego.Router("/rbac/config/index", &controllers.GMConfigController{}, "*:Index")

	beego.Router("/rbac/log/index", &controllers.LogController{}, "*:Index")

	// 列表
	beego.Router("/rbac/sync/index", &controllers.SyncController{}, "*:Index")
	// 发送同步命令
	beego.Router("/rbac/sync/post", &controllers.SyncController{}, "*:Post")
	// 获取进度
	beego.Router("/rbac/sync/process", &controllers.SyncController{}, "*:GetProcess")

	beego.Router("/rbac/sync/upload_app", &controllers.SyncController{}, "*:UploadApp")

	beego.Router("/rbac/sync/sync_app", &controllers.SyncController{}, "*:SyncApp")
	beego.Router("/rbac/sync/sync_conf", &controllers.SyncController{}, "*:SyncConf")

	// beego.Router("/rbac/sync/local", &controllers.SyncController{}, "*:Local")

	// =========== GM =========
	beego.Router("/gm/gamer/index", &controllers.GamerController{}, "*:Index")
	beego.Router("/gm/gamer/search", &controllers.GamerController{}, "*:Search")
	beego.Router("/gm/gamer/add_player", &controllers.GamerController{}, "*:AddPlayer")

	beego.Router("/gm/mail/send_to_gamer", &controllers.GMMailController{}, "*:SendToGamer")
	beego.Router("/gm/mail/index", &controllers.GMMailController{}, "*:Index")

	// 外部接口 下载文件 ?file=
	// beego.Router("/public/download", &controllers.DownController{}, "*:Index")
	ns := beego.NewNamespace("/public/v1",
		beego.NSNamespace("/download",
			beego.NSInclude(
				&controllers.DownController{},
			),
		),
		beego.NSNamespace("/log",
			beego.NSInclude(
				&controllers.InnerLogController{},
			),
		),
	)
	beego.AddNamespace(ns)

}
