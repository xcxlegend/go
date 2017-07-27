package controllers

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"time"

	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/xcxlegend/go/compress"
	m "github.com/xcxlegend/go/lmdgm/models"
)

//DirController 获取文件夹文件信息
type DirController struct {
	BaseController
}

//Index 入口
func (this *DirController) Index() {
	this.TplName = this.GetTemplatetype() + "/dir/index.tpl"
}

func (this *DirController) Sftp() {
	var fpath = beego.AppConfig.String("path.sftp.base")
	var files = this.getFiles(fpath)
	this.ResponseJson(files)
}

//Conf 获取配置文件夹
func (this *DirController) Conf() {
	if this.IsAjax() {
		var fpath = beego.AppConfig.String("path.conf.base")
		var files = this.getFiles(fpath)
		this.ResponseJson(files)
		return
	} else {
		this.TplName = this.GetTemplatetype() + "/dir/conf.tpl"
	}
}

//UploadConf 上传配置文件夹
func (this *DirController) UploadConf() {
	var uploadFile, fh, err = this.GetFile("file")
	// beego.Debug(f, fh, err)
	if err != nil {
		this.ResponseJson(map[string]interface{}{
			"status": false,
			"info":   "文件不存在",
		})
		return
	}
	var ext = path.Ext(path.Base(fh.Filename))
	if ext != ".zip" {
		this.ResponseJson(map[string]interface{}{
			"status": false,
			"info":   "文件格式错误, 请上传 .zip文件",
		})
		return
	}
	// beego.Error(1, err)
	var base = beego.AppConfig.String("path.conf.base")
	var dir = path.Join(base, time.Now().Format("20060102150405"))

	os.Mkdir(dir, 0666)

	var filename = path.Join(dir, fh.Filename)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	// beego.Error(2, err)
	io.Copy(f, uploadFile)
	f.Close()
	var comp compress.CompressTool
	switch ext {
	case ".zip", ".rar":
		comp = new(compress.ZipCompress)
		break
	}
	if comp != nil {
		// beego.Debug("comp:", filename, dir)
		if err := comp.Decompress(filename, dir+string(os.PathSeparator)); err == nil {
			var err = os.Remove(filename)
			// beego.Error(err)
			if err != nil {
				beego.Error(err)
			}
		} else {
			beego.Error(err)
		}
	}
	// defer f.Close()
	defer uploadFile.Close()

	// 加入重写
	var mainRedis = m.GetMainRedis()
	var excludeFile = map[string]bool{"auth.json": true, "pvpmatch.json": true}
	var confdir = path.Join(dir, "conf")
	if fs, err := ioutil.ReadDir(confdir); err == nil {
		// beego.Debug(len(fs))
		for _, f := range fs {
			var filep = path.Join(confdir, f.Name())
			if content, err := ioutil.ReadFile(filep); err == nil {
				var js = map[string]interface{}{}
				json.Unmarshal(content, &js)
				// beego.Debug(js)
				if _, ok := excludeFile[f.Name()]; !ok {
					if _, ok := js["Net"]; ok {
						if net, ok := js["Net"].(map[string]interface{}); ok {
							net["ExtIpUrl"] = "http://myip.ksyun.com/get.php"
						}
					}
				}
				js["Redis"] = []string{fmt.Sprintf("@%s:%d", mainRedis.Host, mainRedis.Port)}
				content, err = json.MarshalIndent(js, "", "	")
				ioutil.WriteFile(filep, content, 0666)
			}
		}
	} else {
		beego.Error(dir, err)
	}

	comp = new(compress.TarCompress)
	// #todo conf.tar.gz const
	err = comp.Compress(dir, path.Join(dir, "conf.tar.gz"))

	// beego.Error(err)
	this.ResponseJson(map[string]interface{}{
		"status": true,
		"info":   "上传成功",
	})
}

//App 获取应用文件夹
func (this *DirController) App() {
	if this.IsAjax() {
		var fpath = beego.AppConfig.String("path.app.base")
		var files = this.getFiles(fpath)
		this.ResponseJson(files)
		return
	} else {
		this.TplName = this.GetTemplatetype() + "/dir/app.tpl"
	}
}

//UploadApp 上传应用文件
func (this *DirController) UploadApp() {
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
	var base = beego.AppConfig.String("path.app.base")
	var dir = path.Join(base, time.Now().Format("20060102150405"))

	os.Mkdir(dir, 0666)

	var filename = path.Join(dir, fh.Filename)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	// beego.Error(2, err)
	io.Copy(f, uploadFile)
	f.Close()

	// defer f.Close()
	defer uploadFile.Close()
	this.ResponseJson(map[string]interface{}{
		"status": true,
		"info":   "上传成功",
	})
}

func (this *DirController) getFiles(fpath string) []*Files {
	var f, err = this.GetDirAllFiles(fpath, true)
	if err != nil {
		beego.Error(err)
	}
	return f
}

func (this *DirController) FileContent() {
	var file = this.GetString("file")
	var f, err = ioutil.ReadFile(file)
	if err != nil {
		this.Ctx.Output.Body([]byte("file error"))
		return
	}
	this.Ctx.Output.Body(f)
}
