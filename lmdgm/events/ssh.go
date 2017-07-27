package events

import (
	"github.com/astaxie/beego"
	"github.com/xcxlegend/go/ssh"
)

const (
	CMD_UNTAR = "cd %s && tar zxvf %s"
)

func RunCmdAndRead(c *ssh.Client, cmd string) string {
	sess, err := c.GetSSHClient().NewSession()
	if err != nil {
		beego.Error(err)
		return ""
	}
	defer sess.Close()
	var syncStdout = ssh.NewSyncStdout()
	sess.Stdout = syncStdout
	sess.Stderr = syncStdout
	sess.Run(cmd)

	var ret = syncStdout.Read()
	syncStdout.Close()
	return string(ret)
}

func RunCmd(c *ssh.Client, cmd string) {
	// beego.Debug("run:", cmd)
	sess, err := c.GetSSHClient().NewSession()
	if err != nil {
		beego.Error(err)
		return
	}
	sess.Run(cmd)
	defer sess.Close()
}
