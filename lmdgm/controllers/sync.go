package controllers

import (
	"fmt"
	"os"
	"time"

	// "github.com/astaxie/beego"
	"io"

	"path"

	"github.com/astaxie/beego"
	"github.com/xcxlegend/go/lib"
	m "github.com/xcxlegend/go/lmdgm/models"
	"github.com/xcxlegend/go/lmdgm/servers"
	"github.com/xcxlegend/go/ssh"
)

const (
	ERR_INFO_SYNC_RUNNING = "同步正在进行"
)

const (
	GROUP_CODE_APP  = "app"
	GROUP_CODE_CONF = "conf"
)

var Servers = map[string]*ssh.LoginOption{
	"P3-DEV-LP": {
		Host:     "47.92.148.191",
		Port:     22,
		User:     "root",
		Password: "lmdSD2017",
	},
	"P3-DEV-DB": {
		Host:     "47.92.150.93",
		Port:     22,
		User:     "root",
		Password: "lmdSD2017",
	},
}

/**
 * 同步分发控制器
 */
type SyncController struct {
	BaseController
	BASE_DIR string
}

type ScheResponse struct {
	*servers.SyncTaskSchedule
	Progress         interface{}
	Server           string
	StartTimeFormat  string
	FinishTimeFormat string
	ServerCode       string
	Host             string
	DestFile         string
}

func (this *SyncController) Index() {
	// var schedules = servers.GetSyncTask().GetSchedules()

	// var groups = this.formatScheduels(schedules)
	// var group = this.GetString("group")
	// if schedule, ok := groups[group]; ok {
	// 	this.Data["schedules"] = schedule
	// }

	// this.Data["group"] = group

	this.TplName = "easyui/sync/index.tpl"
}

func (this *SyncController) Local() {
	// this.Ctx.Output.Body([]byte("http://192.168.1.159:8080/view/P3%E9%A1%B9%E7%9B%AE/job/P3_Android/ws/P3_Android_14/"))
}

func (this *SyncController) UploadApp() {
	var uploadFile, fh, err = this.GetFile("AppInputFile")
	// beego.Debug(f, fh, err)
	if err != nil {
		this.ResponseJson(map[string]interface{}{
			"status": false,
			"info":   "文件不存在",
		})
		return
	}
	// beego.Error(1, err)

	var dir = "D:\\LegendXie\\Ftp"

	var filename = dir + "/" + fh.Filename
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	// beego.Error(2, err)

	io.Copy(f, uploadFile)
	fmd5, err := lib.FileMd5(f)
	if err != nil {
		beego.Error(err)
	}
	f.Close()

	// defer f.Close()
	defer uploadFile.Close()
	this.ResponseJson(map[string]interface{}{
		"status": true,
		"info":   "上传成功",
		"data": map[string]interface{}{
			"file": fh.Filename,
			"path": filename,
			"md5":  fmd5,
		},
	})
	return
}

// 开始
func (this *SyncController) Post() {
	// go this.postServer()

	// const LOCAL_FILE = "D:\\LegendXie\\Ftp\\lmdgm"
	// var _, err = os.Stat(LOCAL_FILE)

	// if err != nil {
	// 	this.ResponseJson(map[string]interface{}{
	// 		"status": false,
	// 		"info":   err.Error(),
	// 	})
	// 	return
	// }

	// for code, option := range Servers {
	// 	var option = &servers.SyncTaskOption{
	// 		LoginOption: option,
	// 		Code:        code,
	// 		Localfile:   LOCAL_FILE,
	// 		Destfile:    "/data/test/lmdgm",
	// 	}
	// 	err = servers.GetSyncTask().Run(option)
	// 	if err != nil {
	// 		this.ResponseJson(map[string]interface{}{
	// 			"status": false,
	// 			"info":   err.Error(),
	// 		})
	// 		return
	// 	}
	// }

	// this.ResponseJson(map[string]interface{}{
	// 	"status": true,
	// })
	// // this.Ctx.Output.Body([]byte("uploading..."))
	// return
}

// 获取进度
func (this *SyncController) GetProcess() {
	var schedules = servers.GetSyncTask().GetSchedules()
	var groups = this.formatScheduels(schedules)
	var group = this.GetString("group")
	schedule, ok := groups[group]
	if ok {
		this.ResponseJson(schedule)
	} else {
		this.ResponseJson([]map[string]interface{}{})
	}
}

func (this *SyncController) formatScheduels(groups servers.SchedulesGroup) map[string][]*ScheResponse {
	var res = map[string][]*ScheResponse{}
	for g, group := range groups {
		res[g] = []*ScheResponse{}
		for k, v := range group.Schedules {
			var scheres = &ScheResponse{
				v,
				0,
				"",
				"",
				"",
				k,
				v.Option.LoginOption.Host,
				v.Option.Destfile,
			}
			if v.Total > 0 {
				if !v.Done && v.Present == v.Total {
					v.Present -= 1
				}
				scheres.Progress = fmt.Sprintf("%.3f", float64(v.Present)/float64(v.Total)*100)
			}
			if server, ok := Servers[k]; ok {
				scheres.Server = server.Host
			}
			scheres.StartTimeFormat = time.Unix(v.StartTime, 0).Format("2006-01-02 15:04:05")
			if v.FinishTime > 0 {
				scheres.FinishTimeFormat = time.Unix(v.FinishTime, 0).Format("2006-01-02 15:04:05")
			}
			res[g] = append(res[g], scheres)
		}
	}
	return res
}

//SyncApp 同步应用
func (this *SyncController) SyncApp() {

	if servers.GetSyncTask().CheckGroupCodeRun(GROUP_CODE_APP) {
		this.ResponseJson(map[string]interface{}{
			"status": false,
			"info":   "running",
		})
		return
	}

	var dir = this.GetString("dir")
	if dir == "" {
		this.ResponseJson(map[string]interface{}{
			"status": false,
			"info":   "param error",
		})
		return
	}
	var base = beego.AppConfig.String("path.app.base")
	var fs, err = this.GetDirAllFiles(path.Join(base, dir), true)
	if err != nil {
		this.ResponseJson(map[string]interface{}{
			"status": false,
			"info":   "path error",
		})
		// beego.Error(err)
		return
	}
	var allServer = m.GetServerAll()
	var DestDir = beego.AppConfig.String("server.path.app.base")
	var groupOp = servers.SyncTaskGroupOption{
		Options:   []*servers.SyncTaskOption{},
		GroupCode: GROUP_CODE_APP,
	}
	for _, server := range allServer {
		var serverOption = &ssh.LoginOption{
			Host:     server.Host,
			Port:     server.Port,
			User:     server.LoginUserName,
			Password: server.LoginPassword,
		}
		// beego.Debug("server: ", server, serverOption)
		for _, f := range fs {
			var option = &servers.SyncTaskOption{
				LoginOption: serverOption,
				Code:        fmt.Sprintf("%s_%s_%s_%s", GROUP_CODE_APP, serverOption.Host, dir, f.Path),
				Localfile:   f.FullPath,
				Destfile:    path.Join(DestDir, f.Path),
				Group:       GROUP_CODE_APP,
			}
			// beego.Debug(option)
			groupOp.Options = append(groupOp.Options, option)
		}
	}

	err = servers.GetSyncTask().RunGroup(&groupOp)
	if err != nil {
		beego.Error("run err: ", err)
		this.ResponseJson(map[string]interface{}{
			"status": false,
			"info":   err.Error(),
		})
		return
	}

	// beego.Debug(fs)
	this.ResponseJson(map[string]interface{}{
		"status": true,
		"info":   "同步开始",
	})
	return
}

//SyncConf 同步配置
func (this *SyncController) SyncConf() {
	if servers.GetSyncTask().CheckGroupCodeRun(GROUP_CODE_CONF) {
		this.ResponseJson(map[string]interface{}{
			"status": false,
			"info":   "running",
		})
		return
	}

	var dir = this.GetString("dir")
	if dir == "" {
		this.ResponseJson(map[string]interface{}{
			"status": false,
			"info":   "param error",
		})
		return
	}
	var base = beego.AppConfig.String("path.conf.base")
	var tarpath = path.Join(base, dir, "conf.tar.gz")

	// if err != nil {
	// 	this.ResponseJson(map[string]interface{}{
	// 		"status": false,
	// 		"info":   "path error",
	// 	})
	// 	// beego.Error(err)
	// 	return
	// }
	var allServer = m.GetServerAll()
	var DestDir = beego.AppConfig.String("server.path.app.base")
	var groupOp = servers.SyncTaskGroupOption{
		Options:   []*servers.SyncTaskOption{},
		GroupCode: GROUP_CODE_CONF,
	}
	for _, server := range allServer {
		var serverOption = &ssh.LoginOption{
			Host:     server.Host,
			Port:     server.Port,
			User:     server.LoginUserName,
			Password: server.LoginPassword,
		}
		// beego.Debug("server: ", server, serverOption)
		// for _, f := range fs {
		// 	if f.IsDir {

		// 		for _, c := range f.Childs {
		// 			if c.IsDir {
		// 				continue
		// 			}
		// 			var option = &servers.SyncTaskOption{
		// 				LoginOption: serverOption,
		// 				Code:        fmt.Sprintf("%s_%s_%s_%s/%s", GROUP_CODE_CONF, serverOption.Host, dir, f.Path, c.Path),
		// 				Localfile:   c.FullPath,
		// 				Destfile:    path.Join(DestDir, f.Path, c.Path),
		// 				Group:       GROUP_CODE_CONF,
		// 			}
		// 			// beego.Debug(option)
		// 			groupOp.Options = append(groupOp.Options, option)
		// 		}

		// 	} else {
		var option = &servers.SyncTaskOption{
			LoginOption: serverOption,
			Code:        fmt.Sprintf("%s_%s_%s_%s", GROUP_CODE_CONF, serverOption.Host, dir, "conf.tar.gz"),
			Localfile:   tarpath,
			Destfile:    path.Join(DestDir, "conf.tar.gz"),
			Group:       GROUP_CODE_CONF,
			CallBack:    servers.SYNCTASKCALLBACKTYPE_UNTAR,
		}

		// beego.Debug(option)
		groupOp.Options = append(groupOp.Options, option)
		// 	}
		// }
	}

	var err = servers.GetSyncTask().RunGroup(&groupOp)
	if err != nil {
		beego.Error("run err: ", err)
		this.ResponseJson(map[string]interface{}{
			"status": false,
			"info":   err.Error(),
		})
		return
	}

	// beego.Debug(fs)
	this.ResponseJson(map[string]interface{}{
		"status": true,
		"info":   "同步开始",
	})
	return
}
