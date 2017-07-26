package ssh

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/astaxie/beego"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type LoginOption struct {
	User     string
	Password string `json:"-"`
	Host     string
	Port     int
}

type Client struct {
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

func (this *Client) GetSftpClient() *sftp.Client {
	return this.sftpClient
}

func (this *Client) GetSSHClient() *ssh.Client {
	return this.sshClient
}

func NewClient(option *LoginOption) (*Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(option.Password))

	clientConfig = &ssh.ClientConfig{
		User:    option.User,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	// fmt.Println(auth)

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", option.Host, option.Port)
	// fmt.Println(addr)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	var client = new(Client)
	client.sftpClient = sftpClient
	client.sshClient = sshClient
	return client, nil
}

//ReadFile 读取文件内容
func (this *Client) ReadFile(filepath string) ([]byte, error) {
	var totalbuf = []byte{}
	f, err := this.GetSftpClient().Open(filepath)
	if err != nil {
		return totalbuf, err
	}
	var buf = make([]byte, 1024)
	for {
		n, _ := f.Read(buf)
		if n == 0 {
			break
		}
		totalbuf = append(totalbuf, buf...)
	}
	return totalbuf, nil
}

//WriteFile 更新文件内容
func (this *Client) WriteFile(filepath, content string) error {
	f, err := this.GetSftpClient().Create(filepath)
	if err != nil {
		beego.Error(f, err)
		return err
	}
	if n, err := f.Write([]byte(content)); err != nil {
		beego.Error(n, err)
		return err
	}
	beego.Debug("update:", filepath, content)
	return nil
}

/**
 * 将本地文件上传到服务器
 * @param  {[type]} this *Client)      Upload(localfile, dest string [description]
 * @return {[type]}      [description]
 */
func (this *Client) Upload(localfile, dest string) error {
	// var localFilePath = "/Users/a1234/Workspaces/Go/project/src/github.com/xcxlegend/go/README.md"
	// var remoteFile = "/root/READMEx.md"
	srcFile, err := os.Open(localfile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := this.sftpClient.Create(dest)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}
	return nil
}

type SyncStdout struct {
	c chan []byte
}

func NewSyncStdout() *SyncStdout {
	// var c = make(chan []byte, 1024)
	return &SyncStdout{make(chan []byte, 1024)}
}
func (s *SyncStdout) Write(p []byte) (n int, err error) {
	go func() {
		s.c <- p
	}()
	return len(p), nil
}

func (s *SyncStdout) Read() []byte {
	// select {
	// case r := <-r.c:
	// 	return r
	// 	break
	// case <-time.After(3 * time.Second):
	// 	break
	// }

	return <-s.c
}

func (s *SyncStdout) Close() {
	close(s.c)
}
