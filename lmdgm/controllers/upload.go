package controllers

import (
	"fmt"
	"io"

	"github.com/astaxie/beego"
	"github.com/beego/admin/src/models"
	"github.com/xcxlegend/go/compress"
	// "net/url"
	"os"
	"path"
	"strings"
)

var ZIPEXT = map[string]bool{
	".zip": true,
}

type UploadController struct {
	BaseController
	BASE_DIR string
}

func (this *UploadController) Prepare() {
	this.BASE_DIR = beego.AppConfig.String("path.sftp.base")
}

func (this *UploadController) Index() {
	this.Data["base_dir"] = this.BASE_DIR
	this.TplName = "easyui/upload/index.tpl"
}

func (this *UploadController) Upload() {
	var uploadFile, fh, err = this.GetFile("file")
	beego.Debug(uploadFile, fh.Filename, err)
	if err != nil {
		this.ResponseJson(map[string]interface{}{
			"status": false,
			"info":   "文件不存在",
		})
		return
	}
	// beego.Error(1, err)
	var filepath = this.GetString("path")
	var dir = filepath
	if dir == "" {
		dir = this.BASE_DIR
	} else {
		dir = path.Join(this.BASE_DIR, dir)
	}
	var fhfilename = strings.Replace(fh.Filename, "\\", "/", -1)
	var filename = dir + "/" + path.Base(fhfilename)
	f, err := os.Create(filename) //  OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		this.ResponseJson(map[string]interface{}{
			"status": false,
			"info":   "服务器错误:路径错误",
		})
		return
	}
	io.Copy(f, uploadFile)
	// beego.Error(n, err, filename, path.Base(fh.Filename))
	var ext = path.Ext(path.Base(f.Name()))
	beego.Debug("ext", ext)
	var auto_unzip = this.GetString("auto_unzip")
	var _, zipok = ZIPEXT[ext]
	beego.Debug(auto_unzip, zipok)
	f.Close()
	if auto_unzip == "on" && zipok {
		var comp compress.CompressTool
		switch ext {
		case ".zip", ".rar":
			comp = new(compress.ZipCompress)
			break
		}
		if comp != nil {
			beego.Debug("comp:", filename, dir)
			if err := comp.Decompress(filename, dir+string(os.PathSeparator)); err == nil {
				var err = os.Remove(filename)
				beego.Error(err)
			} else {
				beego.Error(err)
			}
		}
	}
	// defer f.Close()
	defer uploadFile.Close()

	this.DBLog(models.LOGNODE_UPLOAD_POST, fmt.Sprintf(DBLOGNODEREMARK_TPL_UPLOAD_POST, fh.Filename, filepath))
	this.ResponseJson(map[string]interface{}{
		"status": true,
		"info":   "上传成功",
		"data": map[string]interface{}{
			"file": fh.Filename,
			"path": filepath,
		},
	})
}

func (this *UploadController) Dir() {
	var path_name = this.GetString("path", "")
	path_name = strings.TrimLeft(path_name, "/")
	if path_name == "../" {
		path_name = ""
	}
	path_name = path.Join(this.BASE_DIR, path_name)
	fi_dir, err := os.Open(path_name)
	if err != nil {
		this.ResponseJson(map[string]interface{}{
			"status": false,
			"info":   "文件不存在",
		})
		return
	}

	var files = []*Files{}
	fis, err_readdir := fi_dir.Readdir(-1)
	if err_readdir != nil {
		this.ResponseJson(map[string]interface{}{
			"status": false,
			"info":   "非文件夹",
		})
		return
	}
	for _, fi := range fis {
		var f = new(Files)
		f.IsDir = fi.IsDir()
		f.Path = fi.Name()
		f.Size = fi.Size()
		files = append(files, f)
	}

	this.ResponseJson(map[string]interface{}{
		"status": true,
		"data":   files,
	})
}

func (this *UploadController) Down() {
	var path_name = this.GetString("path", "")
	path_name = strings.TrimLeft(path_name, "/")
	if path_name == "../" {
		path_name = ""
	}
	var fileFullPath = path.Join(this.BASE_DIR, path_name)
	// beego.Debug(fileFullPath)

	file, err := os.Open(fileFullPath)

	if err != nil {
		beego.Error(err)
		this.Abort("404 NOT FOUND")
		return
	}

	defer file.Close()

	this.Ctx.Output.Download(fileFullPath)

}
