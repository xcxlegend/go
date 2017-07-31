package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/xcxlegend/go/compress"

	"github.com/astaxie/beego/httplib"
)

const VERSION = "1.0.0"

var isDebug = false

var CMDS = map[string]bool{
	"upload": true,
	"help":   true,
}

var CmdHandler = map[string]func(){
	"upload": handleUplaod,
	"help":   handleHelp,
}

type JSData struct {
	Info   string `json:"info"`
	Status bool   `json:"status"`
}

// LMD GM CMD TOOL
func main() {

	var args = os.Args

	if len(args) < 2 {
		handleHelp()
		return
	}
	var cmd = args[1]
	if _, ok := CMDS[cmd]; !ok {
		handleHelp()
		return
	}
	if f, ok := CmdHandler[cmd]; ok {
		os.Args = append(os.Args[0:0], os.Args[1:]...)
		// fmt.Println(os.Args)
		f()
		return
	} else {
		fmt.Println(`[E] 程序错误, 命令不存在`)
	}
	return
}

func handleUplaod() {

	var (
		dir      bool
		iszip    bool
		file     string
		rpath    string
		username string
		password string
	)

	// flag.BoolVar(&dir, "dir", false, "")
	flag.BoolVar(&iszip, "zip", false, "")
	flag.BoolVar(&isDebug, "v", false, "")
	flag.StringVar(&file, "file", "", "no file")
	flag.StringVar(&rpath, "path", "", "")
	flag.StringVar(&username, "username", "", "")
	flag.StringVar(&password, "password", "", "")

	flag.Parse()
	// fmt.Println(flag.Args())
	// fmt.Println("dir=", dir)
	// fmt.Println("file=", file)
	// fmt.Println("path=", path)
	// fmt.Println("username=", username)
	// fmt.Println("password=", password)
	const (
		DOMAIN = "http://120.92.146.36:8080"
		// DOMAIN     = "http://localhost:8999"
		URL_LOGIN  = DOMAIN + "/public/login?isajax=1"
		URL_UPLOAD = DOMAIN + "/rbac/upload/post"
	)
	debug("[I] 开始尝试登录...")
	var req = httplib.Post(URL_LOGIN)
	req.Param("username", username)
	req.Param("password", password)
	req.SetEnableCookie(true)
	var d = new(JSData)
	// var res = map[string]interface{}{}
	req.ToJSON(d)
	// fmt.Println(res)
	// if JSData.Status {
	debug("[I]", d.Info)
	// } else {

	// }
	file = strings.Replace(file, "\\", "/", -1)
	// req = httplib.Post("http://localhost:8999/rbac/upload/dir")
	req = httplib.Post(URL_UPLOAD)
	req.Header("X-Requested-With", "XMLHttpRequest")
	req.Param("path", rpath)
	var autozip = ""

	var fs, err = os.Open(file)
	if err != nil {
		fmt.Println("[E] 文件不存在: ", file)
		return
	}
	info, err := fs.Stat()
	if info.IsDir() {
		dir = true
	}

	if dir {
		debug("[I] 尝试打包文件夹: dir:", file)
		// fmt.Println(path.Dir(file))
		autozip = path.Join(path.Dir(file), "upload.zip")
		// fmt.Println(autozip)
		var comp = new(compress.ZipCompress)
		var err = comp.Compress(file, autozip)
		if err != nil {
			fmt.Println("[E] 上传失败: 打包文件夹失败.dir:", autozip)
			return
		}
		debug("[I] 打包完成: ", autozip)
		file = autozip
		iszip = true
	}
	if iszip {
		req.Param("auto_unzip", "on")
	}
	debug("[I] 尝试上传文件: ", file)
	req.SetEnableCookie(true)
	var filename = file
	filename = strings.Replace(filename, "\\", "/", -1)
	// fmt.Println(filename)
	req.PostFile("file", filename)
	req.Header("X-Requested-With", "XMLHttpRequest")
	var d2 = new(JSData)
	// var res, _ = req.Bytes()
	req.ToJSON(d2)
	// if JSData.Status {
	debug("[I]", d2.Info)
	// }
	if autozip != "" {
		os.Remove(autozip)
	}
	fmt.Println("[I] 上传完成..")
}

func handleHelp() {
	fmt.Printf(`
version: %s

cmds	<第二命令>

%10s	%-20s 	

Flags	<参数>

%10s	%-20s
%10s	%-20s
%10s	%-20s
%10s	%-20s
%10s	%-20s
%10s	%-20s

demo	<事例>

./tools upload --file=D:/Res --username= --password=
./tools.exe upload --zip --file=D:/Res.zip --username= --password=

`, VERSION,
		"upload", "上传文件",
		"--help", "帮助",
		"--zip", "表示上传的.zip包需要解压 如果file为目录则不需要设置 会自动设置为解压",
		"--file", "本地文件地址",
		"--path", "目标目录,根目录可以不加",
		"--username", "登录名",
		"--password", "登录密码",
	)
	return
}

func debug(info ...interface{}) {
	if isDebug {
		fmt.Println(info...)
	}
}
