/* package controllers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"time"





















	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/xcxlegend/go/ssh"
	cssh "golang.org/x/crypto/ssh"
	"fmt"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type TerminalIO struct {
	buf chan []byte
}

func NewTerminalIO() *TerminalIO {
	return &TerminalIO{make(chan []byte, 1024)}
}

func (t *TerminalIO) Write(p []byte) (n int, err error) {
	t.buf <- p
	return len(p), nil
}

func (t *TerminalIO) Read(p []byte) (n int, err error) {
	p = <-t.buf
	return len(p), nil
}

func (t *TerminalIO) Close() error {
	close(t.buf)
	return nil
}

type WSController struct {
	BaseController
	ws   *websocket.Conn
	send chan []byte
	tio  *TerminalIO
}

const (
	writeWait  = 10 * time.Second
	readWait   = 5 * 60 * time.Second
	pingPeriod = (60 * time.Second * 9) / 10
)

func (this *WSController) Get() {
	id, _ := this.GetInt64("id", 0)
	if id <= 0 {
		this.Rsp(false, "param error")
		return
	}

	var c, err = getSSHClientByServerId(id)
	if err != nil {
		this.Rsp(false, err.Error())
		return
	}
	this.ws, err = upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	this.send = make(chan []byte, 256)
	if err != nil {
		beego.Error(err)
		return
	}
	this.tio = NewTerminalIO()

	this.getTerminal(c)
	defer func() {
		if this.ws != nil {
			this.ws.Close()
		}
	}()

	go this.wsWrite()
	go this.cmdRead()
	this.wsRead()
	// return
}

func (this *WSController) cmdRead() {
	defer func() {
		close(this.send)
	}()
	buf := make([]byte, 1024)
	// var logs = []string{}
	for {
		size, err := this.tio.Read(buf)
		if err != nil {
			beego.Error(err)
			return
		}
		beego.Debug("get:", string(buf))
		safeMessage := base64.StdEncoding.EncodeToString([]byte(buf[:size]))
		this.send <- []byte(string(safeMessage))
	}
}

func (this *WSController) cmdWrite(b []byte) error {
	_, err := this.tio.Write(b)
	return err
}

func (this *WSController) wsRead() {

	for {
		this.ws.SetReadDeadline(time.Now().Add(readWait))
		op, r, err := this.ws.NextReader()
		if err != nil {
			beego.Error(err, op, r)
			return
		}
		beego.Error(err, op, r)
		switch op {
		case websocket.PongMessage:
			this.ws.SetReadDeadline(time.Now().Add(readWait))
		case websocket.TextMessage:
			message, err := ioutil.ReadAll(r)
			beego.Debug("ws msg:", message, err)
			err = this.cmdWrite(message)
			if err != nil {
				beego.Error(err)
				// return
			}

		}
	}
}

func (this *WSController) wsWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-this.send:
			if !ok {
				this.w(websocket.CloseMessage, []byte{})
				return
			}
			if err := this.w(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := this.w(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
func (this *WSController) w(opCode int, data []byte) error {
	this.ws.SetWriteDeadline(time.Now().Add(writeWait))
	beego.Debug("cmd:", string(data))
	return this.ws.WriteMessage(opCode, data)
}

func (this *WSController) getTerminal(c *ssh.Client) {
	sess, err := c.GetSSHClient().NewSession()
	if err != nil {
		beego.Error(err)
		return
	}
	defer sess.Close()
	this.tio = new(TerminalIO)
	sess.Stdout = this.tio
	sess.Stderr = this.tio
	sess.Stdin = this.tio

	modes := cssh.TerminalModes{
		cssh.ECHO:          1,
		cssh.TTY_OP_ISPEED: 14400,
		cssh.TTY_OP_OSPEED: 14400,
	}

	if err := sess.RequestPty("xterm-256color", 25, 100, modes); err != nil {
		fmt.Println("创建终端出错: ", err)
		return
	}
	err = sess.Shell()
	if err != nil {
		fmt.Println("执行Shell出错: ", err)
		return
	}
	go sess.Wait()
	return
}
*/
package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/kr/pty"
	m "github.com/xcxlegend/go/lmdgm/models"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type SSH_Config struct {
	Ip       string
	Username string
	Password string
	Port     int
	Command  *exec.Cmd
	Pty      *os.File
}

func (this *SSH_Config) Close() {
	if this.Command != nil {
		this.Command.Process.Signal(syscall.SIGHUP)
	}
	if this.Pty != nil {
		this.Pty.Close()
	}
}
func (this *SSH_Config) Login() error {
	var err error
	this.Command = exec.Command("ssh", "-o", "StrictHostKeyChecking=no", fmt.Sprintf("-p %d", this.Port), this.Username+"@"+this.Ip)
	this.Pty, err = pty.Start(this.Command)
	if err != nil {
		beego.Error(this.Command, this.Pty, err)
		// beego.Error(err)
		return err
	}
	i := 0
	for {
		if i >= 10 {
			return errors.New("login error")
		}
		buf := make([]byte, 1024)
		size, err2 := this.Pty.Read(buf)
		// beego.Debug(string(buf))
		if strings.Contains(string(buf), "refused") {
			return errors.New(string(buf))
		}
		if err2 != nil {
			return err
		}
		var rcmd = string([]byte(buf[:size]))
		if !strings.Contains(rcmd, "password") && !strings.Contains(rcmd, "Password") {
			i++
			continue
		}
		this.Pty.Write([]byte(this.Password + "\n"))
		if err != nil {
			panic(err)
		}
		if err != nil {
			return errors.New("login error")
		}
		return nil
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WSController struct {
	beego.Controller
	ws   *websocket.Conn
	send chan []byte
	SSH_Config
}

const (
	writeWait  = 10 * time.Second
	readWait   = 5 * 60 * time.Second
	pingPeriod = (60 * time.Second * 9) / 10
)

func (this *WSController) Get() {
	id, err := this.GetInt64("id", 0)
	var serv = m.GetServerById(id)
	if err != nil {
		// beego.Error(err)
		return
	}
	// beego.Debug("ser:", serv)
	this.Ip = serv.Host
	this.Username = serv.LoginUserName
	this.Password = serv.LoginPassword
	this.Port = serv.Port
	this.ws, err = upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	this.send = make(chan []byte, 256)
	if err != nil {
		beego.Error(err)
		return
	}
	defer func() {
		if this.ws != nil {
			this.ws.Close()
		}
		this.Close()
	}()
	if err = this.Login(); err != nil {
		beego.Error(err)
		http.Error(this.Ctx.ResponseWriter, err.Error(), 400)
		return
	}

	go this.wsWrite()
	go this.cmdRead()
	this.wsRead()
	// return
}

func (this *WSController) cmdRead() {
	defer func() {
		close(this.send)
	}()
	buf := make([]byte, 1024)
	// var logs = []string{}
	for {
		size, err := this.Pty.Read(buf)
		if err != nil {
			beego.Error(err)
			return
		}
		safeMessage := base64.StdEncoding.EncodeToString([]byte(buf[:size]))

		this.send <- []byte(string(safeMessage))

	}
}

func (this *WSController) cmdWrite(b []byte) error {
	_, err := this.Pty.Write(b)
	return err
}

func (this *WSController) wsRead() {

	for {
		this.ws.SetReadDeadline(time.Now().Add(readWait))
		op, r, err := this.ws.NextReader()
		if err != nil {
			beego.Error(err, op, r)
			return
		}
		beego.Error(err, op, r)
		switch op {
		case websocket.PongMessage:
			this.ws.SetReadDeadline(time.Now().Add(readWait))
		case websocket.TextMessage:
			message, err := ioutil.ReadAll(r)
			beego.Debug("ws msg:", message, err)
			err = this.cmdWrite(message)
			if err != nil {
				beego.Error(err)
				// return
			}

		}
	}
}

func (this *WSController) wsWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-this.send:
			if !ok {
				this.w(websocket.CloseMessage, []byte{})
				return
			}
			if err := this.w(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := this.w(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
func (this *WSController) w(opCode int, data []byte) error {
	this.ws.SetWriteDeadline(time.Now().Add(writeWait))
	// beego.Debug("cmd:", string(data))
	return this.ws.WriteMessage(opCode, data)
}
