package models

import (
	"errors"
	"log"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

//Server 服务器信息模型
type Server struct {
	Id            int64
	Lastlogintime time.Time `orm:"auto_now_add;type(datetime)" form:"-"`
	Createtime    time.Time `orm:"type(datetime);auto_now_add" `
	ServerName    string    `orm:"size(32)" form:"ServerName"`
	OutHost       string    `orm:"unique;size(60)" form:"OutHost"`
	Host          string    `orm:"unique;size(60)" form:"Host"`
	Port          int       `orm:"size(6)" form:"Port"`
	LoginUserName string    `orm:"size(32)" form:"LoginUserName"`
	LoginPassword string    `orm:"size(32)" form:"LoginPassword"`
	Remark        string    `orm:"null;size(200)" form:"Remark" valid:"MaxSize(200)"`
	Status        int       `orm:"size(1)" form:"Status"`
	IsMount       bool      `orm:"-"`
}

func init() {
	orm.RegisterModel(new(Server))
}

//TableName 表名
func (u *Server) TableName() string {
	return "server"
}

func (u *Server) Valid(v *validation.Validation) {
	if false {
		v.SetError("Repassword", "两次输入的密码不一样")
	}
}

//验证用户信息
func checkServer(u *Server) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(&u)
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return errors.New(err.Message)
		}
	}
	return nil
}

//GetServerList get server list
func GetServerList(page, page_size int64, sort string) (servers []*Server, count int64) {
	o := orm.NewOrm()
	server := new(Server)
	qs := o.QueryTable(server)
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	qs.Limit(page_size, offset).OrderBy(sort).All(&servers)
	count, _ = qs.Count()
	return servers, count
}

func GetServerAll() []*Server {
	o := orm.NewOrm()
	server := new(Server)
	qs := o.QueryTable(server)
	var servers []*Server
	qs.All(&servers)
	return servers
}

//AddServer 添加服务器
func AddServer(s *Server) (int64, error) {
	if err := checkServer(s); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	id, err := o.Insert(s)
	return id, err
}

//UpdateServer 更新服务器
func UpdateServer(s *Server) (int64, error) {
	if err := checkServer(s); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	server := make(orm.Params)

	if len(s.Host) > 0 {
		server["Host"] = s.Host
	}

	if len(s.ServerName) > 0 {
		server["ServerName"] = s.ServerName
	}

	if s.Port != 0 {
		server["Port"] = s.Port
	}

	if len(s.LoginUserName) > 0 {
		server["LoginUserName"] = s.LoginUserName
	}

	if len(s.LoginPassword) > 0 {
		server["LoginPassword"] = s.LoginPassword
	}

	if s.Status != 0 {
		server["Status"] = s.Status
	}

	// if u.Status != 0 {
	// 	user["Status"] = u.Status
	// }
	if len(server) == 0 {
		return 0, errors.New("update field is empty")
	}
	var table Server
	num, err := o.QueryTable(table).Filter("Id", s.Id).Update(server)
	return num, err
}

func DelServerById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Server{Id: Id})
	return status, err
}

// func GetUserByUsername(username string) (user User) {
// 	user = User{Username: username}
// 	o := orm.NewOrm()
// 	o.Read(&user, "Username")
// 	return user
// }

//GetServerById 根据ID获取server信息
func GetServerById(id int64) (s Server) {
	s = Server{Id: id}
	o := orm.NewOrm()
	o.Read(&s, "Id")
	return s
}
