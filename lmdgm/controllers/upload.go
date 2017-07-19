package controllers

import (
	"github.com/astaxie/beego"
	"github.com/xcxlegend/go/compress"
	"io"
	// "net/url"
	"os"
	"path"
	"strings"
)

var ZIPEXT = map[string]bool{
	".zip": true,
	".rar": true,
}

type Files struct {
	Path  string `json:"path"`
	IsDir bool   `json:"is_dir"`
	Size  int64  `json:"size"`
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
	// beego.Debug(f, fh, err)
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
	var filename = dir + "/" + fh.Filename
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	// beego.Error(2, err)
	io.Copy(f, uploadFile)
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
