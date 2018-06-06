package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net"
	"net/mail"
	"net/smtp"

	m "github.com/xcxlegend/go/lmdgm/models"
)

// Message 邮件发送数据
type Message struct {
	Subject   string            // 标题
	Content   io.Reader         // 支持html的消息主体
	To        []string          // 邮箱地址
	Extension map[string]string // 发送邮件消息体扩展项
}

// Sender 提供邮件发送接口
type Sender interface {
	// Send 发送邮件
	// msg 邮件发送数据
	// isMass 是否是群发,默认为一对一发送
	Send(msg *Message, isMass bool) error

	// AsyncSend 异步发送邮件
	// msg 邮件发送数据
	// isMass 是否是群发,默认为一对一发送
	// handle 发送完成之后的处理函数，如果发送失败,则返回error
	AsyncSend(msg *Message, isMass bool, handle func(err error)) error
}

type GMailConfig struct {
	GmMailHost        string `json:"gm_mail_server"`
	GmMailAddrName    string `json:"gm_name"`
	GmMailAddrAddress string `json:"gm_addr"`
	GmMailPwd         string `json:"gm_passwd"`
}

func (p *GMailConfig) String() string {
	return fmt.Sprintf("host:%s, name:%s, addr:%s, pwd:%s", p.GmMailHost, p.GmMailAddrName, p.GmMailAddrAddress, p.GmMailPwd)
}

const (
	DBCONFIG_CODE_EMAIL = "email"
)

// CreateMessage 创建邮件内容
// subject 邮件主题
// content 支持html的消息主体
// to 接收者邮箱地址列表
func CreateMessage(subject string, content string, to []string) *Message {
	msg := &Message{
		Subject: subject,
		Content: bytes.NewBufferString(content),
		To:      to,
	}
	return msg
}

// NewSmtpSenderWithGM 使用GM的邮箱配置创建基于smtp的邮件发送实例(使用PlainAuth)
// 从数据库中读取GM邮箱配置的IO慢操作，暂时放到这里， 后续有需要再进行优化
func NewSmtpSenderWithGM() (Sender, error) {
	var c = m.GetConfigByCode(DBCONFIG_CODE_EMAIL)
	var econf = new(GMailConfig)
	if err := json.Unmarshal([]byte(c), econf); err != nil {
		return nil, fmt.Errorf("json unmarshaling failed:%s in %v", err, c)
	}
	if len(econf.GmMailHost) == 0 || len(econf.GmMailAddrAddress) == 0 || len(econf.GmMailPwd) == 0 {
		return nil, fmt.Errorf("check gm mail config")
	}
	from := mail.Address{
		Name:    econf.GmMailAddrName,
		Address: econf.GmMailAddrAddress,
	}
	smtpCli := &SmtpClient{
		addr: econf.GmMailHost,
		from: from,
	}
	host, _, err := net.SplitHostPort(econf.GmMailHost)
	if err != nil {
		return nil, err
	}
	smtpCli.auth = smtp.PlainAuth("", from.Address, econf.GmMailPwd, host)
	return smtpCli, nil
}

// NewSmtpSender 创建基于smtp的邮件发送实例(使用PlainAuth)
// addr 服务器地址
// from 发送者
// authPwd 验证密码
// 如果创建实例发生异常，则返回错误
func NewSmtpSender(addr string, from mail.Address, authPwd string) (Sender, error) {
	smtpCli := &SmtpClient{
		addr: addr,
		from: from,
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	smtpCli.auth = smtp.PlainAuth("", from.Address, authPwd, host)
	return smtpCli, nil
}

// SmtpClient 使用smtp发送邮件
type SmtpClient struct {
	addr string
	from mail.Address
	auth smtp.Auth
}

// Send 发送邮件
func (this *SmtpClient) Send(msg *Message, isMass bool) (err error) {
	if isMass {
		err = this.massSend(msg)
	} else {
		err = this.oneSend(msg)
	}
	return
}

// AsyncSend 异步发送邮件
func (this *SmtpClient) AsyncSend(msg *Message, isMass bool, handle func(err error)) error {
	go func() {
		err := this.Send(msg, isMass)
		handle(err)
	}()
	return nil
}

// oneSend 一对一按顺序发送
func (this *SmtpClient) oneSend(msg *Message) error {
	for _, addr := range msg.To {
		header := this.getHeader(msg.Subject)
		header["To"] = addr
		if msg.Extension != nil {
			for k, v := range msg.Extension {
				header[k] = v
			}
		}
		data := this.getData(header, msg.Content)
		err := smtp.SendMail(this.addr, this.auth, this.from.Address, []string{addr}, data)
		if err != nil {
			return err
		}
	}
	return nil
}

// massSend 群发邮件
func (this *SmtpClient) massSend(msg *Message) error {
	header := this.getHeader(msg.Subject)
	if msg.Extension != nil {
		for k, v := range msg.Extension {
			header[k] = v
		}
	}
	data := this.getData(header, msg.Content)
	return smtp.SendMail(this.addr, this.auth, this.from.Address, msg.To, data)
}

func (this *SmtpClient) getHeader(subject string) map[string]string {
	header := make(map[string]string)
	header["From"] = this.from.String()
	header["Subject"] = mime.QEncoding.Encode("utf-8", subject)
	header["Mime-Version"] = "1.0"
	header["Content-Type"] = "text/html;charset=utf-8"
	header["Content-Transfer-Encoding"] = "Quoted-Printable"
	return header
}

func (this *SmtpClient) getData(header map[string]string, body io.Reader) []byte {
	buf := new(bytes.Buffer)
	for k, v := range header {
		fmt.Fprintf(buf, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(buf, "\r\n")
	io.Copy(buf, body)
	return buf.Bytes()
}
