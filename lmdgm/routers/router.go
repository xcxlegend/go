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
	beego.Router("/rbac/upload", &controllers.UploadController{}, "*:Index")
	beego.Router("/rbac/upload/dir", &controllers.UploadController{}, "*:Dir")
	beego.Router("/rbac/upload/post", &controllers.UploadController{}, "*:Upload")
	beego.Router("/rbac/upload/down", &controllers.UploadController{}, "*:Down")
}
