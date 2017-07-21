package lib

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// FileMd5 计算文件md5值
func FileMd5(f *os.File) (string, error) {
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		// fmt.Println("Copy", err)
		return "", err
	}

	var str = fmt.Sprintf("%x", md5hash.Sum(nil))
	return str, nil

}
