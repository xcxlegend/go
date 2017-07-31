package controllers

import (
	"regexp"

	"github.com/beego/admin/src/models"

	"fmt"

	"strings"

	"path"

	"errors"

	"sort"

	"github.com/astaxie/beego"
	m "github.com/xcxlegend/go/lmdgm/models"
	"github.com/xcxlegend/go/ssh"
)

//ServersController 服务器管理控制器
type ServersController struct {
	BaseController
}

//AppStatResponse 应用配置运行状态
type AppStatResponse struct {
	PID      string `json:"pid"`
	ServerId int64  `json:"server_id"`
	Stime    string `json:"stime"`
	Conf     string `json:"confile"`
	IsRun    bool   `json:"is_run"`
}

type AppStatResponseList []*AppStatResponse

func (fs AppStatResponseList) Len() int {
	return len(fs)
}

func (fs AppStatResponseList) Less(i, j int) bool {
	return fs[i].Conf < fs[j].Conf
}

func (fs AppStatResponseList) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

// func (this *ServersController) AddShow() {
// 	m := models.Server{}
// 	this.Data["m"] = m
// 	this.Data["op"] = "a"
// 	this.TplName = "easyui/servers/form.html"
// }

// Index 列表页
func (this *ServersController) Index() {

	page, _ := this.GetInt64("page")
	page_size, _ := this.GetInt64("rows")
	sort := this.GetString("sort")
	order := this.GetString("order")
	if len(order) > 0 {
		if order == "desc" {
			sort = "-" + sort
		}
	} else {
		sort = "Id"
	}

	if this.IsAjax() {
		servers, count := m.GetServerList(page, page_size, sort)
		for _, s := range servers {
			this.getServerMount(s)
		}
		this.Data["json"] = &map[string]interface{}{"total": count, "rows": &servers}

		this.ServeJSON()
		return
	} else {
		// tree := this.GetTree()
		// this.Data["tree"] = &tree
		// this.Data["servers"] = &servers
		// if this.GetTemplatetype() != "easyui" {
		// 	this.Layout = this.GetTemplatetype() + "/public/layout.tpl"
		// }
		this.TplName = this.GetTemplatetype() + "/servers/index.tpl"
	}
}

//AddServer 添加
func (this *ServersController) AddServer() {
	s := m.Server{}
	if err := this.ParseForm(&s); err != nil {
		//handle error
		this.Rsp(false, err.Error())
		return
	}

	if s.Host == "" {
		// 获取内网地址
		this.getServerHost(&s)
	}

	id, err := m.AddServer(&s)
	if err == nil && id > 0 {

		s.Id = id
		this.DBLogTplData(models.LOGNODE_SERVER_ADD, DBLOGNODEREMARK_TPL_SERVER_ADD, &s)
		// this.DBLog(models.LOGNODE_SERVER_ADD, fmt.Sprintf(DBLOGNODEREMARK_TPL_SERVER_ADD, log))
		this.Rsp(true, "Success")
		return
	}
	this.Rsp(false, err.Error())
	return

}

//DelServer 删除
func (this *ServersController) DelServer() {
	Id, _ := this.GetInt64("Id")
	var old = m.GetServerById(Id)
	status, err := m.DelServerById(Id)
	if err == nil && status > 0 {
		this.DBLogTplData(models.LOGNODE_SERVER_DEL, DBLOGNODEREMARK_TPL_SERVER_DEL, &old)
		this.Rsp(true, "Success")
		return
	}
	this.Rsp(false, err.Error())
	return
}

// func (this *ServersController) ModifyShow() {
// 	id, err := this.GetInt("id", -1)
// 	m := models.Server{Id: id}
// 	err = orm.NewOrm().Read(&m, "id")
// 	if err != nil {
// 		beego.Error(err)
// 		this.ToJsonFail(err.Error())
// 		return
// 	}
// 	this.Data["m"] = m
// 	this.Data["op"] = "m"
// 	this.TplName = "easyui/servers/form.html"
// }

//UpdateServer 更新
func (this *ServersController) UpdateServer() {
	s := m.Server{}
	if err := this.ParseForm(&s); err != nil {
		//handle error
		this.Rsp(false, err.Error())
		return
	}
	var o = m.GetServerById(s.Id)
	if o.Host == "" || s.Host == "" {
		this.getServerHost(&o)
		s.Host = o.Host
	}
	id, err := m.UpdateServer(&s)
	if err != nil {
		this.Rsp(false, err.Error())
		return
	}
	if id > 0 {
		var log = map[string]interface{}{
			"old":    &o,
			"update": this.Input(),
		}
		this.DBLogTplData(models.LOGNODE_SERVER_UPDATE, DBLOGNODEREMARK_TPL_SERVER_UPDATE, log)
		this.Rsp(true, "Success")
		return
	}
	this.Rsp(false, "no update")
	return
}

//getServerHost 获取内网ip
func (this *ServersController) getServerHost(s *m.Server) {
	// beego.Debug(s)
	var c, err = ssh.NewClient(&ssh.LoginOption{
		User:     s.LoginUserName,
		Password: s.LoginPassword,
		Host:     s.OutHost,
		Port:     s.Port,
	})
	if err != nil {
		beego.Error(err)
		this.Rsp(false, "无法连接服务器,添加失败")
		return
	}
	var ipstr = this.RunCmdAndRead(c, "ifconfig")
	fmt.Println(ipstr)
	var iplines = strings.Split(ipstr, `
`)
	fmt.Println(iplines, len(iplines))
	if len(iplines) > 2 {
		var ip = strings.Split(strings.TrimSpace(iplines[1]), " ")
		fmt.Println(ip)
		match, _ := regexp.MatchString(`\d+\.\d+\.\d+\.\d+`, ip[1])
		if match {
			s.Host = ip[1]
		}
	}
}

//GetStat Ajax获取服务器情况 ping, ssh, run
func (this *ServersController) GetStat() {
	var id, _ = this.GetInt64("id", 0)
	if id <= 0 {
		this.Rsp(false, "param error")
		return
	}

	var serv = m.GetServerById(id)
	if serv.Id == 0 {
		this.Rsp(false, "server error")
		return
	}

	var c, err = ssh.NewClient(&ssh.LoginOption{
		User:     serv.LoginUserName,
		Password: serv.LoginPassword,
		Host:     serv.Host,
		Port:     serv.Port,
	})
	if err != nil {
		beego.Error(err)
		this.Rsp(false, "server connect error")
		return
	}
	// 所有的配置文件
	var cfs = this.getServerConfiles(c)

	var runcommands = this.getServerRunCmd(id, c)

	for f, _ := range cfs {
		if _, ok := runcommands[f]; !ok {
			runcommands[f] = &AppStatResponse{
				Conf:     f,
				ServerId: id,
				IsRun:    false,
			}
		}
	}

	// beego.Debug(c, err, string(<-ch))
	// this.Rsp(true, )
	var res = AppStatResponseList{}
	for _, ru := range runcommands {
		res = append(res, ru)
	}
	sort.Sort(res)
	this.Data["json"] = &map[string]interface{}{"total": len(res), "rows": &res}
	this.ServeJSON()
}

//SSHClosePid 关闭应用
func (this *ServersController) SSHClosePid() {
	var id, _ = this.GetInt64("id", 0)
	var pid = this.GetString("pid")
	if id <= 0 || pid == "" {
		this.Rsp(false, "param error")
		return
	}

	var serv = m.GetServerById(id)
	if serv.Id == 0 {
		this.Rsp(false, "server error")
		return
	}

	var c, err = ssh.NewClient(&ssh.LoginOption{
		User:     serv.LoginUserName,
		Password: serv.LoginPassword,
		Host:     serv.Host,
		Port:     serv.Port,
	})
	sess, err := c.GetSSHClient().NewSession()
	if err != nil {
		beego.Error(err)
		return
	}
	defer sess.Close()
	sess.Run(fmt.Sprintf("kill %s", pid))
	var log = map[string]interface{}{
		"id":      serv.Id,
		"host":    serv.OutHost,
		"confile": this.GetString("confile"),
	}
	this.DBLogTplData(models.LOGNODE_SERVER_SSH_CLOSE, DBLOGNODEREMARK_TPL_SERVER_SSH_CLOSE, log)
	this.Rsp(true, "run over")
}

//SSHStartApp 启动应用
func (this *ServersController) SSHStartApp() {
	var id, _ = this.GetInt64("id", 0)
	var confile = this.GetString("confile")
	if id <= 0 || confile == "" {
		this.Rsp(false, "param error")
		return
	}

	var serv = m.GetServerById(id)
	if serv.Id == 0 {
		this.Rsp(false, "server error")
		return
	}

	var c, err = ssh.NewClient(&ssh.LoginOption{
		User:     serv.LoginUserName,
		Password: serv.LoginPassword,
		Host:     serv.Host,
		Port:     serv.Port,
	})

	var cfs = this.getServerConfiles(c)
	if _, ok := cfs[confile]; !ok {
		this.Rsp(false, "配置文件错误")
		return
	}

	var runcommands = this.getServerRunCmd(id, c)
	if runing, ok := runcommands[confile]; ok && runing.IsRun {
		this.Rsp(false, "运行中")
		return
	}
	sess, err := c.GetSSHClient().NewSession()
	if err != nil {
		beego.Error(err)
		return
	}
	defer sess.Close()
	// go func() {
	err = sess.Run(fmt.Sprintf(beego.AppConfig.String("server.app.runcmd"),
		beego.AppConfig.String("server.path.app.base"),
		beego.AppConfig.String("server.app.name"),
		confile,
	))
	// beego.Error(err)
	// }()
	var log = map[string]interface{}{
		"id":      serv.Id,
		"host":    serv.OutHost,
		"confile": this.GetString("confile"),
	}
	this.DBLogTplData(models.LOGNODE_SERVER_SSH_START, DBLOGNODEREMARK_TPL_SERVER_SSH_START, log)
	this.Rsp(true, "run over")
}

//SSHMount 挂载服务器硬盘
func (this *ServersController) SSHMount() {
	var id, _ = this.GetInt64("id", 0)
	if id <= 0 {
		this.Rsp(false, "param error")
		return
	}

	var serv = m.GetServerById(id)
	if serv.Id == 0 {
		this.Rsp(false, "server error")
		return
	}

	var c, err = ssh.NewClient(&ssh.LoginOption{
		User:     serv.LoginUserName,
		Password: serv.LoginPassword,
		Host:     serv.Host,
		Port:     serv.Port,
	})
	if err != nil {
		this.Rsp(false, "server error")
		return
	}
	this.RunCmdAndRead(c, `mkfs.ext4 /dev/vdb && mkdir -p /data && mount /dev/vdb /data && echo "/dev/vdb /data ext4 defaults 0 0">>/etc/fstab`)
	// beego.Debug("ret:", ret)
	this.DBLogTplData(models.LOGNODE_SERVER_SSH_MOUNT, DBLOGNODEREMARK_TPL_SERVER_SSH_MOUNT, &serv)
	this.Rsp(true, "run over")
}

//GetConfContent 获取配置文件内容
func (this *ServersController) GetConfContent() {
	var id, _ = this.GetInt64("id", 0)
	var filename = strings.TrimSpace(this.GetString("file"))
	if id <= 0 || filename == "" {
		this.Rsp(false, "param error")
		return
	}

	var c, err = getSSHClientByServerId(id)
	if err != nil {
		this.Rsp(false, err.Error())
		return
	}

	var filepath = path.Join(beego.AppConfig.String("server.path.conf.base"), filename)

	content, err := c.ReadFile(filepath)
	if err != nil {
		this.Rsp(false, err.Error())
		return
	}
	this.Ctx.Output.Body(content)
}

func (this *ServersController) Terminal() {
	var id, _ = this.GetInt64("id", 0)
	this.Data["id"] = id
	this.Data["host"] = this.Ctx.Request.Host
	this.TplName = this.GetTemplatetype() + "/servers/terminal.tpl"
}

//UpdateConfContent 更新配置文件内容
func (this *ServersController) UpdateConfContent() {
	var id, _ = this.GetInt64("id", 0)
	var filename = strings.TrimSpace(this.GetString("file"))
	var content = this.GetString("content")
	if id <= 0 || filename == "" {
		this.Rsp(false, "param error")
		return
	}

	var c, err = getSSHClientByServerId(id)
	if err != nil {
		this.Rsp(false, err.Error())
		return
	}

	var filepath = path.Join(beego.AppConfig.String("server.path.conf.base"), filename)
	// beego.Debug("file:", filepath)
	var f, _ = c.ReadFile(filepath)
	if err := c.WriteFile(filepath, content); err != nil {
		this.Rsp(false, err.Error())
		return
	}
	// var log = map[string]interface{}{
	// 	"id":      id,
	// 	"host":    c.Option.Host,
	// 	"content": strings.TrimSpace(string(f)),
	// 	"update":  content,
	// }
	this.DBLogTpl(models.LOGNODE_SERVER_SSH_EDIT_JSON, DBLOGNODEREMARK_TPL_SERVER_SSH_EDIT_JSON,
		id,
		c.Option.Host,
		filename,
		string(f),
		content,
	)
	// this.DBLogTplData(models.LOGNODE_SERVER_SSH_EDIT_JSON, DBLOGNODEREMARK_TPL_SERVER_SSH_EDIT_JSON, log)
	this.Rsp(true, "")
}

func getSSHClientByServerId(id int64) (*ssh.Client, error) {
	var serv = m.GetServerById(id)
	if serv.Id == 0 {
		// this.Rsp(false, "server error")
		return nil, errors.New("server error")
	}
	// beego.Debug("serv:", serv)
	var c, err = ssh.NewClient(&ssh.LoginOption{
		User:     serv.LoginUserName,
		Password: serv.LoginPassword,
		Host:     serv.Host,
		Port:     serv.Port,
	})
	return c, err
}

//getServerConfiles 获取配置文件
func (this *ServersController) getServerConfiles(c *ssh.Client) map[string]bool {
	var dir = beego.AppConfig.String("server.path.conf.base")
	var cfs = map[string]bool{}
	fis, err := c.GetSftpClient().ReadDir(dir)
	// beego.Debug("confiles:", fis, err)
	if err != nil {
		beego.Error(err)
		return cfs
	}
	var ext = beego.AppConfig.String("server.conf.ext")
	for _, f := range fis {
		if strings.HasSuffix(f.Name(), ext) {
			cfs[f.Name()] = true
		}
	}
	return cfs
}

func (this *ServersController) getServerRunCmd(serverId int64, c *ssh.Client) map[string]*AppStatResponse {
	var runcommands = map[string]*AppStatResponse{}
	sess, err := c.GetSSHClient().NewSession()
	if err != nil {
		beego.Error(err)
		return runcommands
	}
	defer sess.Close()
	var syncStdout = ssh.NewSyncStdout()
	sess.Stdout = syncStdout
	sess.Stderr = syncStdout
	sess.Run(fmt.Sprintf("ps x -e -o 'pid,stime,args' | grep '%s' | grep %s",
		beego.AppConfig.String("server.app.name"),
		beego.AppConfig.String("server.conf.ext"),
	))

	var ret = syncStdout.Read()
	syncStdout.Close()

	fmt.Println(string(ret))

	var lines = []string{}
	lines = strings.Split(string(ret), `
`)
	// beego.Debug("lines:", lines, len(lines))

	for _, l := range lines {
		l = strings.TrimSpace(l)
		var i = strings.Split(l, " ")
		// beego.Debug("line:", i, len(i))
		if len(i) >= 5 && path.Base(i[2]) == beego.AppConfig.String("server.app.name") {
			var cf = path.Base(i[4])
			runcommands[cf] = &AppStatResponse{
				PID:      i[0],
				Conf:     cf,
				ServerId: serverId,
				Stime:    i[1],
				IsRun:    true,
			}
		}
	}
	return runcommands
}

func (this *ServersController) RunCmdAndRead(c *ssh.Client, cmd string) string {
	sess, err := c.GetSSHClient().NewSession()
	if err != nil {
		beego.Error(err)
		return ""
	}
	defer sess.Close()
	var syncStdout = ssh.NewSyncStdout()
	sess.Stdout = syncStdout
	sess.Stderr = syncStdout
	sess.Run(cmd)

	var ret = syncStdout.Read()
	syncStdout.Close()
	return string(ret)
}

func (this *ServersController) getServerMount(server *m.Server) {
	var c, err = ssh.NewClient(&ssh.LoginOption{
		User:     server.LoginUserName,
		Password: server.LoginPassword,
		Host:     server.OutHost,
		Port:     server.Port,
	})
	if err != nil {
		return
	}
	server.IsMount = this.getMountStat(c)
}

func (this *ServersController) getMountStat(c *ssh.Client) bool {
	var ret = this.RunCmdAndRead(c, "df -lh")
	var lines = strings.Split(ret, `
`)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, "/data") {
			return true
		}
	}
	return false
}
