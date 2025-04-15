package tools

import (
	"log"

	"gopkg.in/gomail.v2"
)

type MailerConfig struct {
	Host         string
	Port         int
	AuthEmail    string
	AuthPassword string
}

func NewMailer(sender string, to []string, cc []mailCC, subject string, attachments []string, config mailConfig) Mailer {
	return &mailer{
		Sender:      sender,
		To:          to,
		CC:          cc,
		Subject:     subject,
		Attachments: attachments,
		Config:      config,
	}
}

type Mailer interface {
	SendMail(doLog bool) error
}

type mailer struct {
	Sender      string
	To          []string
	CC          []mailCC
	Subject     string
	Attachments []string
	Config      mailConfig
	Body        mailBody
}

type mailCC struct {
	Email string
	Name  string
}

type mailBody struct {
	Content     string
	ContentType string
}

type mailConfig struct {
	HOST          string
	PORT          int
	AUTH_EMAIL    string
	AUTH_PASSWORD string
}

func (m *mailer) doMailLog(doLog bool, err error) {
	if !doLog {
		return
	}

	if err != nil {
		log.Println("====Email Error====")
	}

	log.Println("From :", m.Sender)
	log.Println("To :", m.To)
	log.Println("CC :", m.CC)
	log.Println("Subject :", m.Subject)
	log.Println("Attachment :", m.Attachments)

	if err != nil {
		log.Println("Error :", err)
	}
}

func (m *mailer) SendMail(doLog bool) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", m.Sender)

	for _, to := range m.To {
		mailer.SetHeader("To", to)
	}

	for _, cc := range m.CC {
		mailer.SetAddressHeader("Cc", cc.Email, cc.Name)
	}

	mailer.SetHeader("Subject", m.Subject)
	mailer.SetBody(m.Body.ContentType, m.Body.Content)

	for _, attachment := range m.Attachments {
		mailer.Attach(attachment)
	}

	dialer := gomail.NewDialer(
		m.Config.HOST,
		m.Config.PORT,
		m.Config.AUTH_EMAIL,
		m.Config.AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		m.doMailLog(doLog, err)
		return err
	}

	m.doMailLog(doLog, nil)
	return nil
}
