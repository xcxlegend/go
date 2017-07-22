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

	beego.Router("/rbac/servers/add", &controllers.ServersController{}, "*:AddServer")
	beego.Router("/rbac/servers/update", &controllers.ServersController{}, "*:UpdateServer")
	beego.Router("/rbac/servers/del", &controllers.ServersController{}, "*:DelServer")
	beego.Router("/rbac/servers/index", &controllers.ServersController{}, "*:Index")

	beego.Router("/rbac/dir/conf", &controllers.DirController{}, "*:Conf")
	beego.Router("/rbac/dir/upload_conf", &controllers.DirController{}, "*:UploadConf")
	beego.Router("/rbac/dir/app", &controllers.DirController{}, "*:App")
	beego.Router("/rbac/dir/file_content", &controllers.DirController{}, "*:FileContent")

	// 列表
	beego.Router("/rbac/sync/index", &controllers.SyncController{}, "*:Index")
	// 发送同步命令
	beego.Router("/rbac/sync/post", &controllers.SyncController{}, "*:Post")
	// 获取进度
	beego.Router("/rbac/sync/process", &controllers.SyncController{}, "*:GetProcess")

	beego.Router("/rbac/sync/upload_app", &controllers.SyncController{}, "*:UploadApp")

	// beego.Router("/rbac/sync/local", &controllers.SyncController{}, "*:Local")

	// 外部接口 下载文件 ?file=
	// beego.Router("/public/download", &controllers.DownController{}, "*:Index")
	ns := beego.NewNamespace("/public/v1",
		beego.NSNamespace("/download",
			beego.NSInclude(
				&controllers.DownController{},
			),
		),
	)
	beego.AddNamespace(ns)

}
