package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/xcxlegend/go/lmdgm/controllers:DownController"] = append(beego.GlobalControllerRouter["github.com/xcxlegend/go/lmdgm/controllers:DownController"],
		beego.ControllerComments{
			Method: "Index",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/xcxlegend/go/lmdgm/controllers:InnerLogController"] = append(beego.GlobalControllerRouter["github.com/xcxlegend/go/lmdgm/controllers:InnerLogController"],
		beego.ControllerComments{
			Method: "Upload",
			Router: `/upload`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

}
