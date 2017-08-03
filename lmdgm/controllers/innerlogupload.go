package controllers

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/astaxie/beego"
)

//InnerLogController 日志上传
type InnerLogController struct {
	InnerBaseController
}

// @Title Post Upload logs
// @Description 上传
// @Success 200 {"status": true}
// @Param   file     formData   file true    "以文件类型上传"
// @Failure 401 Unauthorized
// @router /upload [post]
func (this *InnerLogController) Upload() {
	var uploadFile, fh, err = this.GetFile("file")
	// beego.Debug(f, fh, err)
	if err != nil {
		this.ResponseJson(map[string]interface{}{
			"status": false,
			"info":   "文件不存在",
		})
		return
	}

	// beego.Error(1, err)
	var dir = beego.AppConfig.String("path.log.base")
	os.Mkdir(dir, 0666)
	var file = strings.Replace(fh.Filename, "\\", "/", -1)
	var filename = path.Join(dir, path.Base(file))
	beego.Debug("file:", filename)
	f, err := os.Create(filename) //  OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		beego.Error("create: ", filename, err)
	}
	// beego.Error(2, err)
	_, err = io.Copy(f, uploadFile)
	if err != nil {
		beego.Error("copy: ", err)
	}
	f.Close()
	// defer f.Close()
	defer uploadFile.Close()
	this.ResponseJson(map[string]interface{}{
		"status": true,
	})
}
