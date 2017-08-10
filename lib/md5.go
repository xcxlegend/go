package lib

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"time"
	"unsafe"
)

// FileMd5 计算文件md5值
func FileMd5(f io.Reader) (string, error) {
	md5Ctx := md5.New()
	io.Copy(md5Ctx, f)
	// beego.Debug("md5:", n, e)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr), nil
}

func Md5ByByte(b []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(b)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

type UUID struct {
	Timestamp int
	ServerId  int16
	Index     int16
}

var uuidIndex int16 = 0

func GetUUID(server_id int16) int64 {
	uuidIndex++
	//fmt.Println(uuidIndex)
	uuid := UUID{int(time.Now().Unix()), server_id, uuidIndex}
	//fmt.Println(uuid)
	return *((*int64)(unsafe.Pointer(&uuid)))
}
