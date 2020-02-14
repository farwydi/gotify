package posix

import (
	"encoding/base64"
	"errors"
	"github.com/farwydi/gotify"
	"net/smtp"
	"strings"
)

func NewPosixAdapter(from string, to []string, smtpAddr string) func() (gotify.Adapter, error) {
	if smtpAddr == "" {
		return func() (gotify.Adapter, error) {
			return nil,
				errors.New("posixAdapter: options 'SmtpAddr' must be not equal empty string")
		}
	}

	if len(to) == 0 {
		return func() (gotify.Adapter, error) {
			return nil,
				errors.New("posixAdapter: options 'To' must be not equal empty array")
		}
	}

	if from == "" {
		return func() (gotify.Adapter, error) {
			return nil,
				errors.New("posixAdapter: options 'From' must be not equal empty string")
		}
	}

	return func() (gotify.Adapter, error) {
		return &posixAdapter{
			conn:     nil,
			smtpAddr: smtpAddr,
			from:     from,
			to:       to,
		}, nil
	}
}

type posixAdapter struct {
	conn     *smtp.Client
	smtpAddr string
	from     string
	to       []string
}

func (ad *posixAdapter) connect() (err error) {
	if ad.conn == nil {
		ad.conn, err = smtp.Dial(ad.smtpAddr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (posixAdapter) Format(text []gotify.Line) []byte {
	size := 0
	for _, tx := range text {
		for _, element := range tx {
			switch t := element.(type) {
			case gotify.CODE:
				size += len(t) + 11
			case gotify.B:
				size += len(t) + 7
			case gotify.TextElement:
				size += len(t)
			case string:
				size += len(t)
			case []byte:
				size += len(t)
			default:
				// TODO: Добавить обработку ошибки типа
				return nil
			}
		}
		size += 5
	}

	complete := make([]byte, size)
	n := 0
	for _, tx := range text {
		for _, element := range tx {
			switch t := element.(type) {
			case gotify.CODE:
				n += copy(complete[n:], "<pre>")
				n += copy(complete[n:], t)
				n += copy(complete[n:], "</pre>")
			case string:
				n += copy(complete[n:], t)
			case gotify.B:
				n += copy(complete[n:], "<b>")
				n += copy(complete[n:], t)
				n += copy(complete[n:], "</b>")
			case gotify.TextElement:
				n += copy(complete[n:], t)
			}
		}
		n += copy(complete[n:], "</br>")
	}

	return complete
}

func (ad posixAdapter) Send(subject string, message ...gotify.Line) error {
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	c, err := smtp.Dial(ad.smtpAddr)
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.Mail(r.Replace(ad.from)); err != nil {
		return err
	}

	for i := range ad.to {
		ad.to[i] = r.Replace(ad.to[i])
		if err = c.Rcpt(ad.to[i]); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	msg := "To: " + strings.Join(ad.to, ",") + "\r\n" +
		"From: " + ad.from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString(ad.Format(message))

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
