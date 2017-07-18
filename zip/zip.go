package main

import (
	"archive/zip"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"os"
	"path"
	"strings"
)

func DeCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			beego.Error(1, err)
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		fmt.Println(filename)
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			beego.Error(2, err)
			continue
		}
		w, err := os.Create(filename)
		if err != nil {

			beego.Error(3, err)
			continue
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			beego.Error(4, err)
			continue
		}
		w.Close()
		rc.Close()
	}
	return nil
}

func DeCompressZip() {
	const File = "D:\\LegendXie\\测试文档\\test2.zip"
	const dir = "D:\\LegendXie\\测试文档\\z\\"
	os.Mkdir(dir, 0777) //创建一个目录

	cf, err := zip.OpenReader(File) //读取zip文件
	if err != nil {
		fmt.Println(err)
	}
	defer cf.Close()
	for _, file := range cf.File {
		rc, err := file.Open()
		if err != nil {
			fmt.Println(err)
		}

		f, err := os.Create(dir + file.Name)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		n, err := io.Copy(f, rc)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(n)
	}

}

func getDir(fpath string) string {
	return path.Dir(fpath)
	if strings.HasSuffix(fpath, "\\") || strings.HasSuffix(fpath, "/") {
		fpath = strings.TrimRight(fpath, "\\")
		fpath = strings.TrimRight(fpath, "/")
	}
	// fmt.Println(fpath)
	return fpath
}

func main() {
	DeCompress("D:\\LegendXie\\测试文档\\test.zip", "D:\\LegendXie\\测试文档\\z\\")
	// DeCompressZip()
}
