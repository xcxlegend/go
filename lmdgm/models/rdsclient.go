package models

import (
	"antnet"
	"fmt"
	"github.com/xcxlegend/go/lmdgm/pb"
	"gopkg.in/redis.v5"
	"sort"
)

const (
	REDIS_KEY_SYSTEM_MAIL   = "sysmail"
	REDIS_KEY_GAMER_MAIN    = "gamer.%v"
	REDIS_KEY_GAMER_PLAYERS = "gamer.%v.player"
	REDIS_KEY_GAMER_PACK    = "gamer.%v.pack"
	REDIS_KEY_GAMER_MAIL    = "gamer.%v.mail"
)

const (
	CHANNL_NOTIFY = "0"
	GM_SERVER_ID  = 0
)

const (
	RedisNotifyNone = iota
	RedisNotifyNewGamer
	RedisNotifyGamerOnline
	RedisNotifyGamerUpdate
	RedisNotifyGamerOffline
	RedisNotifyFriendRequest
	RedisNotifyFriendRequestAccept
	RedisNotifyNewPVPResult
	RedisNotifyGamerPlayerLvUpdate         //恭喜xx玩家达到xx等级
	RedisNotifyGamerPlayerInitSpeciality   //恭喜xx玩家解锁槽位获得xx特质
	RedisNotifyGamerPlayerUpdateSpeciality //恭喜xx玩家更新特质获得xx特质
	RedisNotifyGamerNewMail
)

// 根据redisID获取client
func GetRdsClientByRid(id int64) *redis.Client {
	var conf *Redis
	if id > 0 {
		conf = GetRedisById(id)
	} else {
		conf = GetMainRedis()
	}
	var c = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: "",
		PoolSize: 5,
	})
	return c
}

// ==========  PBreader ==========

// 获取redis里gamer信息
func RDSGetGamerInfo(gid int32, client *redis.Client) *pb.Gamer {
	if client == nil {
		client = GetRdsClientByRid(0)
	}
	var data, _ = client.HGet(fmt.Sprintf(REDIS_KEY_GAMER_MAIN, gid), "main").Result()
	var gamer = new(pb.Gamer)
	antnet.ParseDBStr(data, gamer)
	return gamer
}

// 获取redis里玩家球员信息
func RDSGetPlayers(gid int32, client *redis.Client) []*pb.Player {
	if client == nil {
		client = GetRdsClientByRid(0)
	}
	var datas, _ = client.HGetAll(fmt.Sprintf(REDIS_KEY_GAMER_PLAYERS, gid)).Result()
	var slice = []*pb.Player{}
	for _, data := range datas {
		var p = new(pb.Player)
		if antnet.ParseDBStr(data, p) {
			slice = append(slice, p)
		}
	}
	return slice
}

// 给玩家添加球员
func RDSAddPlayers(gid int32, player *pb.Player, client *redis.Client) bool {
	if client == nil {
		client = GetRdsClientByRid(0)
	}
	return client.HSetNX(fmt.Sprintf(REDIS_KEY_GAMER_PLAYERS, gid), fmt.Sprintf("%v", *player.Id), antnet.DBStr(player)).Val()
}

// 获取redis里玩家背包信息
func RDSGetPackGoods(gid int32, client *redis.Client) []*pb.ItemConfig {
	if client == nil {
		client = GetRdsClientByRid(0)
	}
	var datas, _ = client.HGetAll(fmt.Sprintf(REDIS_KEY_GAMER_PACK, gid)).Result()
	var slice = []*pb.ItemConfig{}
	for _, data := range datas {
		var p = new(pb.ItemConfig)
		if antnet.ParseDBStr(data, p) {
			slice = append(slice, p)
		}
	}
	return slice
}

type MailSlice []*pb.Mail

func (fs MailSlice) Len() int {
	return len(fs)
}

func (fs MailSlice) Less(i, j int) bool {
	return *fs[i].Time < *fs[j].Time
}

func (fs MailSlice) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

// 读取玩家所有的邮件
func RDSGetGamerMails(gid int32, client *redis.Client) MailSlice {
	if client == nil {
		client = GetRdsClientByRid(0)
	}
	var datas, _ = client.HGetAll(fmt.Sprintf(REDIS_KEY_GAMER_MAIL, gid)).Result()
	var slice = MailSlice{}
	for _, data := range datas {
		var p = new(pb.Mail)
		if antnet.ParseDBStr(data, p) {
			slice = append(slice, p)
		}
	}
	sort.Sort(slice)
	return slice
}

// 发送邮件给玩家
func RDSSendMailToGamer(gid int32, mail *pb.Mail, client *redis.Client) bool {
	if client == nil {
		client = GetRdsClientByRid(0)
	}

	var data = antnet.DBStr(mail)
	if client.HSet(fmt.Sprintf(REDIS_KEY_GAMER_MAIL, gid), fmt.Sprintf("%v", mail.GetId()), data).Val() {
		if err := client.Publish(CHANNL_NOTIFY,
			fmt.Sprintf("%d.%d.%d.%v",
				GM_SERVER_ID,
				RedisNotifyGamerNewMail,
				gid,
				data,
			)).Err(); err == nil {
			return true
		}
	}
	return false
}

func RDSGetAllSysMail() []*pb.Mail {
	var client = GetRdsClientByRid(0)
	var datas, _ = client.HGetAll(REDIS_KEY_SYSTEM_MAIL).Result()
	var slice = MailSlice{}
	for _, data := range datas {
		var p = new(pb.Mail)
		if antnet.ParseDBStr(data, p) {
			slice = append(slice, p)
		}
	}
	sort.Sort(slice)
	return slice
}
