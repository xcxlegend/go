package models

import (
	"errors"

	"github.com/astaxie/beego/orm"
)

//Config 数据库信息模型
type Config struct {
	Id     int64
	Code   string `orm:"size(32)" form:"Code"`
	Config string `orm:"null;size(200)" form:"Config"`
}

func init() {
	orm.RegisterModel(new(Config))
}

//TableName 表名
func (u *Config) TableName() string {
	return "config"
}

func GetConfigByCode(code string) string {
	var s = Config{Code: code}
	o := orm.NewOrm()
	o.Read(&s, "Code")
	var conf = ""
	if s.Id > 0 {
		conf = s.Config
	}
	return conf
}

//GetConfigList get Config list
func GetConfigList(page, page_size int64, sort string) (Configs []*Config, count int64) {
	o := orm.NewOrm()
	cfgs := new(Config)
	qs := o.QueryTable(cfgs)
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	qs.Limit(page_size, offset).OrderBy(sort).All(&Configs)
	count, _ = qs.Count()
	return Configs, count
}

func GetConfigAll() []*Config {
	o := orm.NewOrm()
	r := new(Config)
	qs := o.QueryTable(r)
	var Configs []*Config
	qs.All(&Configs)
	return Configs
}

//AddConfig 添加服务器
func AddConfig(s *Config) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(s)
	return id, err
}

//UpdateConfig 更新服务器
func UpdateConfig(s *Config) (int64, error) {

	o := orm.NewOrm()
	r := make(orm.Params)

	if len(s.Config) > 0 {
		r["Config"] = s.Config
	}

	if len(r) == 0 {
		return 0, errors.New("update field is empty")
	}
	var table Config
	num, err := o.QueryTable(table).Filter("id", s.Id).Update(r)
	return num, err
}

func DelConfigById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Config{Id: Id})
	return status, err
}

//GetConfigById 根据ID获取Config信息
func GetConfigById(id int64) (s Config) {
	s = Config{Id: id}
	o := orm.NewOrm()
	o.Read(&s, "Id")
	return s
}
