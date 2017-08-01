package servers

// 同步任务管理
import (
	"errors"
	"os"
	"sync"
	"time"

	"path"

	"fmt"

	"github.com/astaxie/beego"
	"github.com/xcxlegend/go/lib"
	"github.com/xcxlegend/go/lmdgm/events"
	"github.com/xcxlegend/go/ssh"
)

var (
	ERROR_RUNNING             = errors.New("task is running")
	ERROR_GROUP_OPTION_LENGTH = errors.New("group none options")
)

//SyncTaskCode 同步标识
// type SyncTaskCode string

// const (
// 	//SYNCTASKCODE_APP 同步app的标识
// 	SYNCTASKCODE_APP SyncTaskCode = "app"
// 	//SYNCTASKCODE_CONF 同步conf的标识
// 	SYNCTASKCODE_CONF = "conf"
// )

/*



//SyncTaskManage 同步任务管理
type SyncTaskManage struct {
	list map[SyncTaskCode]SyncTask
}

func (this *SyncTaskManage) Init() {
	this.list = map[SyncTaskCode]SyncTask{}
}

func (this *SyncTaskManage) Run(task SyncTask) error {
	var code = task.GetCode()
	if _, ok := this.list[code]; ok {
		return ERROR_RUNNING
	}
	if task.Run() != nil {
		this.list[code] = task
	}
}

func (this *SyncTaskManage) GetTaskSchedule(code SyncTaskCode) interface{} {

}

type SyncTask interface {
	Run() error
	GetCode() SyncTaskCode
	GetSchedule() interface{}
}

type SyncTaskApp struct {
	LoginOption *ssh.LoginOption
	Localfile   string
	Destfile    string
}

func (this *SyncTaskApp) Run() error {

}

func (this *SyncTaskApp) GetCode() SyncTaskCode {

}

func (this *SyncTaskApp) GetSchedule() interface{} {

} */

var syncTask *SyncTask
var _ = beego.Debug

func init() {
	syncTask = new(SyncTask)
	syncTask.Init()
}

func GetSyncTask() *SyncTask {
	return syncTask
}

type SyncTaskCallBackType string

const (
	SYNCTASKCALLBACKTYPE_UNTAR SyncTaskCallBackType = "untar"
)

type SyncTaskSchedule struct {
	Option     *SyncTaskOption `json:"option"`
	StartTime  int64           `json:"start_time"`
	FinishTime int64           `json:"finish_time"`
	Total      int64           `json:"total"`
	Present    int64           `json:"present"`
	Error      string          `json:"error"`
	Done       bool            `json:"done"`
	Md5        string          `json:"md5"`
}

type SyncTaskScheduleGroup struct {
	GroupCode string                       `json:"group_code"`
	Schedules map[string]*SyncTaskSchedule `json:"schedules"`
	Total     int                          `json:"total"`   // total schedules
	Present   int                          `json:"present"` // finish schedules
	Done      bool                         `json:"done"`
}

type SyncTaskGroupOption struct {
	Options   []*SyncTaskOption
	GroupCode string
}

type SyncTaskOption struct {
	LoginOption *ssh.LoginOption     `json:"server"`
	Code        string               `json:"code"` // 单文件一致标志 同一个文件只能发起一次
	Localfile   string               `json:"local_file"`
	Destfile    string               `json:"dest_file"`
	Group       string               `json:"group"` // 组一致标识 同一个组只能发起一次
	CallBack    SyncTaskCallBackType `json:"-"`
}

type SchedulesGroup map[string]*SyncTaskScheduleGroup

//SyncTask 同步任务
type SyncTask struct {
	Schedules     SchedulesGroup // group -> code
	SchedulesLock *sync.RWMutex
}

func (this *SyncTask) Init() {
	this.Schedules = make(SchedulesGroup)
	this.SchedulesLock = new(sync.RWMutex)
}

func (this *SyncTask) IsRun(option *SyncTaskOption) bool {
	// this.SchedulesLock.RLock()
	// defer this.SchedulesLock.RUnlock()
	// if group, ok := this.Schedules[option.Group]; ok {
	// 	if v, ok := group[option.Code]; ok && !v.Done {
	// 		return true
	// 	}
	// }
	return false
}

func (this *SyncTask) GetSchedules() SchedulesGroup {
	return this.Schedules
}

func (this *SyncTask) IsGroupRun(group *SyncTaskGroupOption) bool {
	return this.CheckGroupCodeRun(group.GroupCode)
}

//CheckGroupCodeRun 检查是否在运行
func (this *SyncTask) CheckGroupCodeRun(code string) bool {
	this.SchedulesLock.RLock()
	defer this.SchedulesLock.RUnlock()
	if group, ok := this.Schedules[code]; ok && !group.Done {
		return true
	}
	return false
}

func (this *SyncTask) RunGroup(group *SyncTaskGroupOption) error {
	if len(group.GroupCode) == 0 {
		return ERROR_GROUP_OPTION_LENGTH
	}

	if this.IsGroupRun(group) {
		return ERROR_RUNNING
	}

	var scheduleGroup = new(SyncTaskScheduleGroup)
	scheduleGroup.GroupCode = group.GroupCode
	scheduleGroup.Total = len(group.Options)
	scheduleGroup.Schedules = make(map[string]*SyncTaskSchedule)
	this.SchedulesLock.Lock()
	for _, option := range group.Options {
		var schedule = new(SyncTaskSchedule)
		schedule.StartTime = time.Now().Unix()
		schedule.Option = option
		scheduleGroup.Schedules[option.Code] = schedule
		go this.handle(option, schedule, scheduleGroup)
	}
	this.Schedules[group.GroupCode] = scheduleGroup
	this.SchedulesLock.Unlock()
	return nil
}

func (this *SyncTask) Run(option *SyncTaskOption) error {
	// if this.IsRun(option) {
	// 	return ERROR_RUNNING
	// }
	// var schedule = new(SyncTaskSchedule)
	// schedule.StartTime = time.Now().Unix()
	// schedule.Option = option
	// this.SchedulesLock.Lock()
	// if sches, ok := this.Schedules[option.Group]; !ok {
	// 	this.Schedules[option.Group] = make(map[string]*SyncTaskSchedule)
	// } else {
	// 	sches[option.Code] = schedule
	// }
	// this.SchedulesLock.Unlock()
	// go this.handle(option, schedule)
	return nil
}

func (this *SyncTask) finish(option *SyncTaskOption, group *SyncTaskScheduleGroup) {
	this.SchedulesLock.Lock()
	defer this.SchedulesLock.Unlock()
	if sche, ok := group.Schedules[option.Code]; ok {
		sche.Done = true
		sche.FinishTime = time.Now().Unix()
		group.Present++
		if group.Present == group.Total {
			group.Done = true
		}
	}
}

// func (this *SyncTask) closeOnError(code string) {
// 	this.SchedulesLock.Lock()
// 	defer this.SchedulesLock.Unlock()
// 	if schedule, ok := this.Schedules[code]; ok {
// 		schedule.Done = true
// 		schedule.FinishTime = time.Now().Unix()
// 	}
// }

func (this *SyncTask) handle(option *SyncTaskOption, schedule *SyncTaskSchedule, group *SyncTaskScheduleGroup) {
	var c, err = ssh.NewClient(option.LoginOption)
	defer this.finish(option, group)
	if err != nil {
		schedule.Error = err.Error()
		return
	}
	srcFile, err := os.Open(option.Localfile)
	if err != nil {
		schedule.Error = err.Error()
		return
	}
	defer srcFile.Close()
	var dir = path.Dir(option.Destfile)
	err = c.GetSftpClient().Mkdir(dir)
	if err != nil {
		// beego.Error(err)
	}
	err = c.GetSftpClient().Remove(option.Destfile)
	if err != nil {
		beego.Error(err)
	}
	dstFile, err := c.GetSftpClient().Create(option.Destfile)

	if err != nil {
		schedule.Error = err.Error()
		return
	}
	// defer dstFile.Close()
	var srcStat, _ = srcFile.Stat()
	schedule.Total = srcStat.Size()
	buf := make([]byte, 1024*1024)
	// go func() {
	// 	for {
	// 		if schedule.Done {
	// 			break
	// 		}
	// 		<-time.Tick(1000 * time.Microsecond)
	// 		beego.Debug(schedule.Present, "/", schedule.Total, float64(schedule.Present)/float64(schedule.Total)*100, "%")
	// 	}
	// }()
	// var totalBuf = []byte{}
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf[0:n])
		schedule.Present += int64(len(buf[0:n]))
		// totalBuf = append(totalBuf, buf[0:n]...)
	}
	dstFile.Close()
	schedule.Md5 = "calculating.."
	err = c.GetSftpClient().Chmod(option.Destfile, 0777)
	if err != nil {
		beego.Error(err)
	}
	dstFile, err = c.GetSftpClient().Open(option.Destfile)
	if err != nil {
		beego.Error(err)
	}

	schedule.Md5, _ = lib.FileMd5(dstFile)

	switch option.CallBack {
	case SYNCTASKCALLBACKTYPE_UNTAR:
		events.RunCmd(c, fmt.Sprintf(events.CMD_UNTAR, path.Dir(option.Destfile), path.Base(option.Destfile)))
		c.GetSftpClient().Remove(option.Destfile)
		break

	}

	// schedule.Md5 = lib.Md5ByByte(dstFile)
	schedule.Present = schedule.Total
}
