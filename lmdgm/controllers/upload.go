package controllers

import (
	"github.com/astaxie/beego"
	"io"
	// "net/url"
	"os"
	"path"
	"strings"
)

const BASE_DIR = "D:\\LegendXie\\Ftp"

type Files struct {
	Path  string `json:"path"`
	IsDir bool   `json:"is_dir"`
	Size  int64  `json:"size"`
}

type UploadController struct {
	BaseController
}

func (this *UploadController) Index() {
	this.Data["base_dir"] = BASE_DIR
	this.TplName = "easyui/upload/index.tpl"
}

func (this *UploadController) Upload() {
	var uploadFile, fh, err = this.GetFile("file")
	// beego.Debug(f, fh, err)
	if err != nil {
		this.ResponseJson(map[string]interface{}{
			"status": -1,
			"info":   "文件不存在",
		})
		return
	}
	// beego.Error(1, err)
	var filepath = this.GetString("path")
	var dir = filepath
	if dir == "" {
		dir = BASE_DIR
	} else {
		dir = path.Join(BASE_DIR, dir)
	}

	f, err := os.OpenFile(dir+"/"+fh.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	// beego.Error(2, err)
	io.Copy(f, uploadFile)
	defer f.Close()
	defer uploadFile.Close()
	this.ResponseJson(map[string]interface{}{
		"status": 1,
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
	path_name = path.Join(BASE_DIR, path_name)
	fi_dir, err := os.Open(path_name)
	if err != nil {
		this.ResponseJson(map[string]interface{}{
			"status": -1,
			"info":   "文件不存在",
		})
		return
	}

	var files = []*Files{}
	fis, err_readdir := fi_dir.Readdir(-1)
	if err_readdir != nil {
		this.ResponseJson(map[string]interface{}{
			"status": -1,
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
		"status": 1,
		"data":   files,
	})
}

func (this *UploadController) Down() {
	var path_name = this.GetString("path", "")
	path_name = strings.TrimLeft(path_name, "/")
	if path_name == "../" {
		path_name = ""
	}
	var fileFullPath = path.Join(BASE_DIR, path_name)
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
