package controllers

import (
	// "github.com/astaxie/beego"
	"os"
	"path"

	"encoding/base64"

	"strings"

	"sort"

	"github.com/beego/admin/src/rbac"
)

const (
	DBLOGNODEREMARK_TPL_UPLOAD_POST = `上传了文件: %s To: %s/`
	DBLOGNODEREMARK_TPL_SERVER_ADD  = `添加了服务器: 
	%s`
	DBLOGNODEREMARK_TPL_SERVER_UPDATE = `更新了服务器: 
	%s`
	DBLOGNODEREMARK_TPL_SERVER_DEL = `删除了服务器: 
	%s`
	DBLOGNODEREMARK_TPL_SERVER_SSH_MOUNT = `对服务器运行命令: 硬盘挂载: 
	%s`
	DBLOGNODEREMARK_TPL_SERVER_SSH_START = `对服务器运行命令: 开启应用: 
	%s`
	DBLOGNODEREMARK_TPL_SERVER_SSH_CLOSE = `对服务器运行命令: 关闭应用: 
	%s`
	DBLOGNODEREMARK_TPL_SERVER_SSH_EDIT_JSON = `对服务器运行命令: 编辑配置文件: \r\n
	id: %v,
	host: "%s",
	file: "%s",
	content: "%s",
	update: "%s"
	`
	DBLOGNODEREMARK_TPL_DIR_UPLOAD_CONF = `上传了配置文件到: dir: "%s"`
	DBLOGNODEREMARK_TPL_DIR_UPLOAD_APP  = `上传了应用文件到: dir: "%s"`
	DBLOGNODEREMARK_TPL_SYNC_CONF       = `分发了配置文件IN DIR: "%s" 到服务器: %v`
	DBLOGNODEREMARK_TPL_SYNC_APP        = `分发了应用文件IN DIR: "%s" 到服务器: %v`
	DBLOGNODEREMARK_TPL_REDIS_ADD       = `添加了REDIS: 
	%s`
	DBLOGNODEREMARK_TPL_REDIS_UPDATE = `更新了REDIS: 
	%s`
	DBLOGNODEREMARK_TPL_REDIS_DEL = `删除了REDIS: 
	%s`
	DBLOGNODEREMARK_TPL_CONFIG_ADD = `添加了CONFIG: 
	%s`
	DBLOGNODEREMARK_TPL_CONFIG_UPDATE = `更新了CONFIG: 
	%s`
	DBLOGNODEREMARK_TPL_CONFIG_DEL = `删除了CONFIG: 
	%s`
	DBLOGNODEREMARK_TPL_AREA_ADD = `添加了区服: 
	%s`
	DBLOGNODEREMARK_TPL_AREA_UPDATE = `更新了区服: 
	%s`
	DBLOGNODEREMARK_TPL_MAIL_SEND_TO_GAMER = `给玩家 (GId: %v) 发送了邮件:
	%s`
	DBLOGNODEREMARK_TPL_GAMER_ADD_PLAYER = `给玩家 (GId: %v) 添加了球员:
	%s`
	DBLOGNODEREMARK_TPL_GAMER_UPDATE_PLAYER = `给玩家 (GId: %v) 更新了球员:
	%s`
)

type BaseController struct {
	rbac.CommonController
}

type Files struct {
	FullPath       string   `json:"full_path"`
	FullPathBase64 string   `json:"full_path_64"`
	Path           string   `json:"path"`
	IsDir          bool     `json:"is_dir"`
	Size           int64    `json:"size"`
	ModTime        int64    `json:"mod_time"`
	Base           bool     `json:"base"`
	Childs         []*Files `json:"children"`
}

/**
Len() int
	// Less reports whether the element with
	// index i should sort before the element with index j.
	Less(i, j int) bool
	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
*/

type FilesList []*Files

func (fs FilesList) Len() int {
	return len(fs)
}

func (fs FilesList) Less(i, j int) bool {
	return fs[i].ModTime < fs[j].ModTime
}

func (fs FilesList) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func NewFilesFromf(fi os.FileInfo, fpath string) *Files {
	var f = new(Files)
	f.IsDir = fi.IsDir()
	f.Path = fi.Name()
	f.Size = fi.Size()
	f.FullPath = path.Join(fpath, f.Path)
	f.ModTime = fi.ModTime().Unix()
	f.FullPathBase64 = strings.Replace(base64.StdEncoding.EncodeToString([]byte(f.FullPath)), "=", "*", -1)
	if f.IsDir && f.Childs == nil {
		f.Childs = []*Files{}
	}
	return f
}

func (this *BaseController) ResponseJson(res interface{}) {
	this.Data["json"] = res
	this.ServeJSON()
}

func (this *BaseController) Rsp(status bool, info string) {
	var res = map[string]interface{}{
		"status": status,
		"info":   info,
	}
	this.ResponseJson(res)
}

func (this *BaseController) GetDirFiles(fpath string, revers_sort bool) ([]*Files, error) {
	fi_dir, err := os.Open(fpath)
	var files = []*Files{}
	if err != nil {
		return files, err
	}

	fis, err_readdir := fi_dir.Readdir(-1)
	if err_readdir != nil {
		return files, err
	}
	for _, fi := range fis {
		var f = NewFilesFromf(fi, fpath)
		files = append(files, f)
	}

	if revers_sort && len(files) > 0 {
		var revers_files = []*Files{}
		for index := len(files) - 1; index >= 0; index-- {
			revers_files = append(revers_files, files[index])
		}
		files = revers_files
	}

	return files, nil
}

func (this *BaseController) GetDirAllFiles(fpath string, revers_sort bool) ([]*Files, error) {
	var files = FilesList{}
	fi_dir, err := os.Open(fpath)
	if err != nil {
		return files, err
	}
	fis, err_readdir := fi_dir.Readdir(-1)
	if err_readdir != nil {
		return files, err
	}
	for _, fi := range fis {
		var f = NewFilesFromf(fi, fpath)
		f.Base = true
		this.getDirAllFiles(f)
		files = append(files, f)
	}
	sort.Sort(files)
	if revers_sort && len(files) > 0 {
		var revers_files = []*Files{}
		for index := len(files) - 1; index >= 0; index-- {
			revers_files = append(revers_files, files[index])
		}
		files = revers_files
	}

	return files, nil
}

func (this *BaseController) getDirAllFiles(file *Files) {
	if file.IsDir && len(file.Childs) == 0 {
		fi_dir, err := os.Open(file.FullPath)
		if err != nil {
			return
		}
		fis, err_readdir := fi_dir.Readdir(-1)
		if err_readdir != nil {
			return
		}
		for _, fi := range fis {
			var f = NewFilesFromf(fi, file.FullPath)
			this.getDirAllFiles(f)
			file.Childs = append(file.Childs, f)
		}
	}
}
