package messaging

import (
	"coderhub/conf"
	"fmt"

	"gopkg.in/gomail.v2"
)

type GoMail interface {
	Send(to string, subject string, body string) error
	SendWithHTML(to string, subject string, body string, attachment string) error
}

func NewGoMail() GoMailImpl {
	return GoMailImpl{
		m:        gomail.NewMessage(),
		d:        gomail.NewDialer(conf.GoMail.Host, conf.GoMail.Port, conf.GoMail.Username, conf.GoMail.Password),
		username: conf.GoMail.Username,
	}
}

type GoMailImpl struct {
	// 邮件
	m *gomail.Message
	// 连接
	d *gomail.Dialer
	// 用户名
	username string
}

func (g *GoMailImpl) Send(to string, subject string, body string) error {
	// 设置邮件头, 发件人
	g.m.SetHeader("From", g.username)
	// 设置邮件头, 收件人
	g.m.SetHeader("To", to)
	// 设置邮件头, 主题
	g.m.SetHeader("Subject", subject)
	// 设置邮件正文
	g.m.SetBody("text/html", body)
	// 发送邮件
	return g.d.DialAndSend(g.m)
}

func (g *GoMailImpl) SendWithHTML(to string, subject string, link string) error {
	g.m.SetHeader("From", g.username)
	g.m.SetHeader("To", to)
	g.m.SetHeader("Subject", subject)
	g.m.SetBody("text/html", g.template(link))
	return g.d.DialAndSend(g.m)
}

func (g *GoMailImpl) template(link string) string {
	return fmt.Sprintf(`
		<h1 style="text-align: center;">邮箱密码重置确认</h1>
		<p>尊敬的用户：</p>
		<p>我们收到了您对邮箱密码重置的请求。为了保障您的账户安全，请您通过点击以下链接完成密码重置确认。</p>
		<p>链接：<a href="%s">重置链接</a></p>
		<p>请您注意，该链接仅在10分钟内有效，请尽快完成确认操作。若无法点击链接，请将链接复制到浏览器地址栏中进行访问。</p>
		<p>如果您没有进行密码重置请求，请忽略此封邮件，并及时检查您的账户安全情况。</p>
		<p>再次感谢您对我们的支持与信任！</p>
		<p>祝您使用愉快！</p>
	`, link)
}
