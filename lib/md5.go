package lib

import (
	"crypto/md5"
	"encoding/hex"
	"io"
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
