package ssh

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"time"
)

type LoginOption struct {
	User     string
	Password string
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
	fmt.Println(auth)

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", option.Host, option.Port)

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
