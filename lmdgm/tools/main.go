package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/xcxlegend/go/compress"

	"encoding/json"
	"github.com/astaxie/beego/httplib"
	"github.com/pkg/errors"
	"github.com/xcxlegend/go/ssh"
)

const VERSION = "1.0.0"

var isDebug = false

var CMDS = map[string]bool{
	"upload":       true,
	"help":         true,
	"upgrade":      true,
	"gm_reflector": true,
	"tp_builder":   true,
}

var CmdHandler = map[string]func(){
	"upload":       handleUplaod,
	"help":         handleHelp,
	"upgrade":      hanldUpgrade,
	"gm_reflector": handlGmReflector,
	"tp_builder":   handleTemplateBuilder,
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
		host     string
	)

	// flag.BoolVar(&dir, "dir", false, "")
	flag.BoolVar(&iszip, "zip", false, "")
	flag.BoolVar(&isDebug, "v", false, "")
	flag.StringVar(&file, "file", "", "no file")
	flag.StringVar(&rpath, "path", "", "")
	flag.StringVar(&username, "username", "", "")
	flag.StringVar(&password, "password", "", "")
	flag.StringVar(&host, "host", "", "")

	flag.Parse()
	// fmt.Println(flag.Args())
	// fmt.Println("dir=", dir)
	// fmt.Println("file=", file)
	// fmt.Println("path=", path)
	// fmt.Println("username=", username)
	// fmt.Println("password=", password)
	var DOMAIN = "http://120.92.146.36:8080"
	if host != "" {
		DOMAIN = host
	}
	var (
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
%10s	%-20s

服务器列表:
	程序服: http://111.231.74.30:19527
	策划服: http://120.92.146.36:8081
	稳定服: http://120.92.146.36:8080

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
		"--host", "GM域名,具体域名查看域名列表, 必须以http开头,如果有端口, 必须带端口, 结尾不带/. 如: http://120.92.146.36:8080",
	)
	return
}

func hanldUpgrade() {
	const TARNAME = "temp.tar"
	type GmServerConf struct {
		Id       int      `json:"id"`
		Host     string   `json:"host"`
		Port     int      `json:"port"`
		User     string   `json:"user"`
		Pass     string   `json:"pass"`
		NeedKey  bool     `json:"need_key"`
		BasePath string   `json:"base_path"`
		Run      bool     `json:"run"`
		Scripts  []string `json: "scripts"`
	}
	var gmServerLists = []*GmServerConf{}
	var json, err = ReadFile("gmserver.json")
	if err != nil {
		beego.Error("gmserver.json not exist")
		return
	}
	if err := JsonUnPack(json, &gmServerLists); err != nil {
		beego.Error("gmserver.json format error:", err)
		return
	}

	var local string
	var skey string
	flag.StringVar(&local, "local_dir", "", "")
	flag.StringVar(&skey, "skey", "", "")
	flag.Parse()
	ssh.RSA_KEY, _ = ReadFile(skey)

	// 本地打包
	beego.Info("start pack temp")

	vfile, _ := os.Create(filepath.Join(local, "version.txt"))
	var version = time.Now().Format("2006-01-02 15:04:05")
	vfile.Write([]byte(version))
	os.Chdir(local)
	var c = exec.Command("sh", "pack.sh")
	o, err := c.Output()
	beego.Debug("sh pack.sh:", string(o), err)

	//var tarc = new(compress.TarCompress)
	var tfile = filepath.Join(local, TARNAME)
	/*if err := tarc.Compress(filepath.Join(local, "temp"), tfile); err != nil {
		beego.Error("pack tar fail err: ", err)
	}*/
	if _, err := os.Stat(tfile); err != nil {
		beego.Error("pack fail err:", err)
	}
	sy := &sync.WaitGroup{}
	var tfpLock = sync.Mutex{}
	var serVersions = map[int]string{}
	var serVersionsLock = new(sync.Mutex)
	for _, c := range gmServerLists {
		if !c.Run {
			continue
		}
		sy.Add(1)
		go func(c *GmServerConf, sy *sync.WaitGroup) {
			beego.Info("start server:", c)
			var client, err = ssh.NewClient(&ssh.LoginOption{
				User:     c.User,
				Password: c.Pass,
				Host:     c.Host,
				Port:     c.Port,
			})
			defer func() {
				sy.Done()
				client.GetSSHClient().Close()
				client.GetSftpClient().Close()
			}()
			if err != nil {
				beego.Error("server connect err: ", c.Id, err)
				return
			}

			var sftp = client.GetSftpClient()
			defer sftp.Close()
			fp := filepath.Join(c.BasePath, TARNAME)
			sftp.Remove(fp)
			f, err := sftp.Create(fp)
			defer f.Close()
			if err != nil {
				beego.Error("sftp create tar err:", c.Id, fp, err)
				return
			}
			tfpLock.Lock()
			tfp, err := os.Open(tfile)
			if err != nil {
				beego.Error("open tar file err:", tfile, err)
			}
			defer tfp.Close()
			n, err := io.Copy(f, tfp)
			tfpLock.Unlock()
			beego.Info("copy tar file over:", c.Id, fp, n, err)
			if err != nil {
				beego.Error("sftp copy tar err:", c.Id, err)
				return
			}

			var sc, _ = client.GetSSHClient().NewSession()
			defer client.GetSftpClient().Close()
			var rcmd = fmt.Sprintf("cd %v && tar -zxf %v && \\cp temp/* ./ -rf && rm -rf temp && rm %v -rf", // && rm temp -rf && rm %v -rf
				c.BasePath,
				TARNAME,
				TARNAME,
			)

			o, err := sc.Output(rcmd)
			sc.Close()

			beego.Info("server run cmd:", c.Id, rcmd, string(o), err)

			if len(c.Scripts) > 0 {
				for _, s := range c.Scripts {
					sc, _ = client.GetSSHClient().NewSession()
					rcmd = fmt.Sprintf("cd %v && %v", c.BasePath, s)
					o, err = sc.Output(rcmd)
					beego.Info("run server scripts over:", c.Id, rcmd, string(o), err)
					sc.Close()
				}
			}
			beego.Info("handle server over:", c.Id)
			var verFile = filepath.Join(c.BasePath, "version.txt")
			//verFileh, err := sftp.Open(verFile)
			//var v = make([]byte, 32)
			//var n2 int
			//if err == nil {
			//	n2, err = verFileh.Read(v)
			//} else {
			//	beego.Error("open version file err:", verFile, err)
			//}

			sc, _ = client.GetSSHClient().NewSession()
			v, err := sc.Output(fmt.Sprintf("cat %v", filepath.Join(c.BasePath, "version.txt")))
			sc.Close()

			if err == nil {
				serVersionsLock.Lock()
				serVersions[c.Id] = string(v)
				serVersionsLock.Unlock()
			} else {
				beego.Error("version file error:", verFile, err)
			}
		}(c, sy)
	}
	sy.Wait()
	beego.Info("upgrade version:", version)
	beego.Info("server version:", serVersions)
}

func handlGmReflector() {
	const (
		InitFile     = "pb_reflect_init.go"
		InitFilePath = "src/models"
		PbPath       = "src/pb"
	)
	var baseDir string
	flag.StringVar(&baseDir, "dir", "", "")
	flag.Parse()

	var pathPb = filepath.Join(baseDir, PbPath)
	var pathFile = filepath.Join(baseDir, InitFilePath, InitFile)
	beego.Debug(pathFile)

	fs, err := ioutil.ReadDir(pathPb)
	if err != nil {
		beego.Error("read pb dir err:", pathPb, err)
		return
	}

	var pbstructs = []string{}

	for _, fi := range fs {
		if !strings.HasSuffix(fi.Name(), "pb.go") {
			continue
		}
		var p = filepath.Join(pathPb, fi.Name())
		c, err := ioutil.ReadFile(p)
		if err != nil {
			beego.Error("read pb file err:", p, err)
		}
		var s = parsePbstructs(string(c))
		if len(s) > 0 {
			pbstructs = append(pbstructs, s...)
		}
	}
	var content = `
package models

import (
	"github.com/golang/protobuf/proto"
	"pb"
)

func init() {
	PBReflector = new(PbReflector)
	PBReflector.m = make(map[string]func() proto.Message)
`

	for _, pb := range pbstructs {
		if strings.HasSuffix(pb, "C2S") || strings.HasSuffix(pb, "S2C") {
			continue
		}
		content += fmt.Sprintf(`
PBReflector.Register("pb.%v", func() proto.Message {
		return new(pb.%v)
	})
`, pb, pb)
	}
	content += `
}`

	err = ioutil.WriteFile(pathFile, []byte(content), 0666)

	//var f, err = os.Create(pathFile)
	//if err != nil {
	beego.Error("write init file error:", pathFile, err)
	//}
	//defer f.Close()

}

func parsePbstructs(content string) []string {
	var matched = []string{}
	var regx, err = regexp.Compile("type (.+) struct")
	if err != nil {
		return matched
	}
	var all = regx.FindAllStringSubmatch(content, -1)
	//beego.Debug("parse pb:", all)
	for _, i := range all {
		if len(i) > 1 {
			matched = append(matched, i[1])
		}
	}
	return matched
}

func debug(info ...interface{}) {
	if isDebug {
		fmt.Println(info...)
	}
}

func ReadFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("文件读取错误")
	}
	return data, nil
}

func JsonUnPack(data []byte, msg interface{}) error {
	if data == nil || msg == nil {
		return errors.New("json unpack error")
	}

	err := json.Unmarshal(data, msg)
	if err != nil {
		return errors.New("json unpack error")
	}
	return nil
}
