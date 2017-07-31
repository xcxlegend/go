package controllers

import (
	"github.com/astaxie/beego"
	"os"
	"path"
	"strings"
)

type DownController struct {
	beego.Controller
	BASE_DIR string
}

func (this *DownController) Prepare() {
	this.BASE_DIR = beego.AppConfig.String("path.sftp.base")
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
}

// @Title Get Resourse files
// @Description 下载服务器文件
// @Success 200 {file}
// @Param   file     query   string true       "file=相对路径地址 完整地址为 http://域名:端口/Base url/download?file=xxx.xx"
// @Failure 404 Not Found
// @router / [get]
func (this *DownController) Index() {
	var filename = this.GetString("file")
	if filename == "" {
		this.Abort("404")
		return
	}
	var path_name = strings.TrimLeft(filename, "/")
	if path_name == "../" {
		path_name = ""
	}
	var fileFullPath = path.Join(this.BASE_DIR, path_name)
	// beego.Debug(fileFullPath)

	file, err := os.Open(fileFullPath)

	if err != nil {
		beego.Error(err)
		this.Abort("404")
		return
	}

	defer file.Close()

	this.Ctx.Output.Download(fileFullPath)
}
