package test

import (
	"net/mail"
	"sync"
	"testing"
	"lmdgame.p3.com/events"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Host    = "smtp.qq.com:25"
	From    = mail.Address{Name: "yuixanchao", Address: "305026363@qq.com"}
	FromPwd = "xxxxxxxx"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:UvZ8im02@tcp(localhost:3306)/lmdgm?charset=utf8")

}

func TestSendWithGM(t *testing.T) {
	sender, err := events.NewSmtpSenderWithGM()
	if err != nil {
		t.Error(err)
		return
	}
	msg := events.CreateMessage(
		"GM同步发送邮件测试",
		"<h1>你好，使用GM账号发送的同步测试邮件内容</h1>",
		[]string{"yuxianchao@lmdgame.com"},
	)

	err = sender.Send(msg, false)
	if err != nil {
		t.Error(err)
	}
}

func TestSend(t *testing.T) {
	sender, err := events.NewSmtpSender(Host, From, FromPwd)
	if err != nil {
		t.Error(err)
		return
	}
	msg := events.CreateMessage(
		"同步发送邮件测试",
		"<h1>你好，同步测试邮件内容</h1>",
		[]string{"yuxianchao@lmdgame.com"},
	)

	err = sender.Send(msg, false)
	if err != nil {
		t.Error(err)
	}
}

func TestAsyncSend(t *testing.T) {
	sender, err := events.NewSmtpSender(Host, From, FromPwd)
	if err != nil {
		t.Error(err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	msg := events.CreateMessage(
		"异步发送邮件测试",
		"<h1>你好，异步测试邮件内容</h1>",
		[]string{"yuxianchao@lmdgame.com"},
	)
	err = sender.AsyncSend(msg, false, func(err error) {
		defer wg.Done()
		if err != nil {
			t.Error(err)
		}
	})
	if err != nil {
		t.Error(err)
	}
	wg.Wait()
}