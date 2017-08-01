package test

import (
	"bytes"
	"net/mail"
	"sync"
	"testing"
	"lmdgame.com/p3gm/lmdgm/events"
)

var (
	Host    = "smtp.qq.com:25"
	From    = mail.Address{Name: "yuixanchao", Address: "305026363@qq.com"}
	FromPwd = "xxxxxxxx"
)

func TestSend(t *testing.T) {
	sender, err := events.NewSmtpSender(Host, From, FromPwd)
	if err != nil {
		t.Error(err)
		return
	}
	msg := events.CreateMessage(
		"同步发送邮件测试",
		bytes.NewBufferString("<h1>你好，同步测试邮件内容</h1>"),
		[]string{"xiechixin@lmdgame.com"},
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
		bytes.NewBufferString("<h1>你好，异步测试邮件内容</h1>"),
		[]string{"xiechixin@lmdgame.com"},
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