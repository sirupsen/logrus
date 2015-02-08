package logrus_mail

import (
	"bytes"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	format = "20060102 15:04:05"
)

// MailHook to sends logs by email without authentication.
type MailHook struct {
	AppName string
	c       *smtp.Client
}

// MailAuthHook to sends logs by email with authentication.
type MailAuthHook struct {
	AppName  string
	Host     string
	Port     int
	From     *mail.Address
	To       *mail.Address
	Username string
	Password string
}

// NewMailHook creates a hook to be added to an instance of logger.
func NewMailHook(appname string, host string, port int, from string, to string) (*MailHook, error) {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial(host + ":" + strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	// Validate sender and recipient
	sender, err := mail.ParseAddress(from)
	if err != nil {
		return nil, err
	}
	recipient, err := mail.ParseAddress(to)
	if err != nil {
		return nil, err
	}

	// Set the sender and recipient.
	c.Mail(sender.String())
	c.Rcpt(recipient.String())

	return &MailHook{
		AppName: appname,
		c:       c,
	}, nil

}

// NewMailAuthHook creates a hook to be added to an instance of logger.
func NewMailAuthHook(appname string, host string, port int, from string, to string, username string, password string) (*MailAuthHook, error) {
	// Check if server listens on that port.
	conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port), 3*time.Second)
	if err != nil {
		return nil, err
	} else {
		defer conn.Close()
	}

	// Validate sender and recipient
	sender, err := mail.ParseAddress(from)
	if err != nil {
		return nil, err
	}
	receiver, err := mail.ParseAddress(to)
	if err != nil {
		return nil, err
	}

	return &MailAuthHook{
		AppName:  appname,
		Host:     host,
		Port:     port,
		From:     sender,
		To:       receiver,
		Username: username,
		Password: password}, nil
}

// Fire is called when a log event is fired.
func (hook *MailHook) Fire(entry *logrus.Entry) error {
	wc, err := hook.c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()
	message := createMessage(entry, hook.AppName)
	if _, err = message.WriteTo(wc); err != nil {
		return err
	}
	return nil
}

// Fire is called when a log event is fired.
func (hook *MailAuthHook) Fire(entry *logrus.Entry) error {
	auth := smtp.PlainAuth("", hook.Username, hook.Password, hook.Host)

	message := createMessage(entry, hook.AppName)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		hook.Host+":"+strconv.Itoa(hook.Port),
		auth,
		hook.From.Address,
		[]string{hook.To.Address},
		message.Bytes(),
	)
	if err != nil {
		return err
	}
	return nil
}

// Levels returns the available logging levels.
func (hook *MailAuthHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}
}

// Levels returns the available logging levels.
func (hook *MailHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}
}

func createMessage(entry *logrus.Entry, appname string) *bytes.Buffer {
	body := entry.Time.Format(format) + " - " + entry.Message
	subject := appname + " - " + entry.Level.String()
	message := bytes.NewBufferString(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))
	return message
}
