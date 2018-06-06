package models

import (
	"antnet"
	"errors"
	"fmt"
	"sort"

	"github.com/astaxie/beego"

	"github.com/xcxlegend/go/lmdgm/pb"

	"gopkg.in/redis.v5"

	"github.com/astaxie/beego/validation"
)

type Area struct {
	c *redis.Client
	antnet.RedisModel
}

const (
	AREA_REDIS_KEY = "server.area"
)

type LogicServerInfoSlice []*pb.LogicServerInfo

func (fs LogicServerInfoSlice) Len() int {
	return len(fs)
}

func (fs LogicServerInfoSlice) Less(i, j int) bool {
	return fs[i].GetId() < fs[j].GetId()
}

func (fs LogicServerInfoSlice) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (a *Area) Valid(v *validation.Validation) {
	if false {
		v.SetError("Repassword", "两次输入的密码不一样")
	}
}

func (a *Area) SetRedisClient() bool {
	var redisconf = GetMainRedis()
	if redisconf.Host == "" {
		return false
	}
	a.c = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisconf.Host, redisconf.Port),
		Password: "",
		PoolSize: 5,
	})
	return true
}

//验证用户信息
func checkArea(u *pb.LogicServerInfo) (err error) {
	if u.GetId() <= 0 {
		return errors.New("error data")
	}
	return nil
}

//GetAreaList get Area list
// func GetAreaList(page, page_size int64, sort string) (Areas []*pb.LogicServerInfo, count int64) {

// 	return Areas, count
// }

func GetAreaAll() []*pb.LogicServerInfo {
	var areamod = new(Area)
	var list = LogicServerInfoSlice{}
	if !areamod.SetRedisClient() {
		return list
	}
	var all, err = areamod.c.HGetAll(AREA_REDIS_KEY).Result()
	if err != nil {
		beego.Error(err)
		return list
	}
	for _, v := range all {
		var data = new(pb.LogicServerInfo)
		if areamod.ParseDBStr(v, data) {
			list = append(list, data)
		}
	}
	sort.Sort(list)
	return list
}

//AddArea 添加服务器
func AddArea(s *pb.LogicServerInfo) (int32, error) {
	if err := checkArea(s); err != nil {
		return 0, err
	}
	var areamod = new(Area)
	if !areamod.SetRedisClient() {
		return 0, errors.New("main redis uncollect")
	}
	var ok, err = areamod.c.HSetNX(AREA_REDIS_KEY, antnet.Sprintf("%v", s.GetId()), areamod.DBStr(s)).Result()
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("hsetnx fail")
	}
	return s.GetId(), nil
}

//UpdateArea 更新服务器
func UpdateArea(s *pb.LogicServerInfo) (int32, error) {
	if err := checkArea(s); err != nil {
		return 0, err
	}
	var areamod = new(Area)
	if !areamod.SetRedisClient() {
		return 0, errors.New("main redis uncollect")
	}
	var _, err = areamod.c.HSet(AREA_REDIS_KEY, antnet.Sprintf("%v", s.GetId()), areamod.DBStr(s)).Result()
	if err != nil {
		return 0, err
	}

	return 1, nil
}

/* func DelAreaById(Id int64) (int32, error) {

}
*/

//GetAreaById 根据ID获取pb.LogicServerInfo信息
func GetAreaById(id int32) *pb.LogicServerInfo {
	var areamod = new(Area)
	if !areamod.SetRedisClient() {
		return nil
	}
	var s = new(pb.LogicServerInfo)
	var data = areamod.c.HGet(AREA_REDIS_KEY, antnet.Sprintf("%v", id)).Val()
	areamod.ParseDBStr(data, s)
	return s
}
