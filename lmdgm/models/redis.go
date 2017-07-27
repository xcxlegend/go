package models

import (
	"errors"
	"log"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

//Redis 数据库信息模型
type Redis struct {
	Id         int64
	Updatetime time.Time `orm:"auto_now_add;type(datetime)" form:"-"`
	Createtime time.Time `orm:"type(datetime);auto_now_add" `
	Name       string    `orm:"size(32)" form:"Name"`
	Host       string    `orm:"unique;size(60)" form:"Host"`
	Port       int       `orm:"size(6)" form:"Port"`
	Remark     string    `orm:"null;size(200)" form:"Remark" valid:"MaxSize(200)"`
	// IsMain     bool      `orm:"size(1);default(0)" form:"IsMain"`
}

func init() {
	orm.RegisterModel(new(Redis))
}

//TableName 表名
func (u *Redis) TableName() string {
	return "redis"
}

func (u *Redis) Valid(v *validation.Validation) {
	if false {
		v.SetError("Repassword", "两次输入的密码不一样")
	}
}

//验证用户信息
func checkRedis(u *Redis) (err error) {
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

//GetRedisList get Redis list
func GetRedisList(page, page_size int64, sort string) (Rediss []*Redis, count int64) {
	o := orm.NewOrm()
	Redis := new(Redis)
	qs := o.QueryTable(Redis)
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	qs.Limit(page_size, offset).OrderBy(sort).All(&Rediss)
	count, _ = qs.Count()
	return Rediss, count
}

func GetRedisAll() []*Redis {
	o := orm.NewOrm()
	r := new(Redis)
	qs := o.QueryTable(r)
	var Rediss []*Redis
	qs.All(&Rediss)
	return Rediss
}

//AddRedis 添加服务器
func AddRedis(s *Redis) (int64, error) {
	if err := checkRedis(s); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	id, err := o.Insert(s)
	return id, err
}

//UpdateRedis 更新服务器
func UpdateRedis(s *Redis) (int64, error) {
	if err := checkRedis(s); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	r := make(orm.Params)

	if len(s.Host) > 0 {
		r["Host"] = s.Host
	}

	if len(s.Name) > 0 {
		r["Name"] = s.Name
	}

	if s.Port != 0 {
		r["Port"] = s.Port
	}

	// if u.Status != 0 {
	// 	user["Status"] = u.Status
	// }
	if len(r) == 0 {
		return 0, errors.New("update field is empty")
	}
	var table Redis
	num, err := o.QueryTable(table).Filter("Id", s.Id).Update(r)
	return num, err
}

func DelRedisById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Redis{Id: Id})
	return status, err
}

// func GetUserByUsername(username string) (user User) {
// 	user = User{Username: username}
// 	o := orm.NewOrm()
// 	o.Read(&user, "Username")
// 	return user
// }

//GetRedisById 根据ID获取Redis信息
func GetRedisById(id int64) (s Redis) {
	s = Redis{Id: id}
	o := orm.NewOrm()
	o.Read(&s, "Id")
	return s
}

func GetMainRedis() *Redis {
	o := orm.NewOrm()
	var table Redis
	var r = &Redis{}
	o.QueryTable(table).GroupBy("id").One(r)
	return r
}
