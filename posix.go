package gotify

import (
	"encoding/base64"
	"errors"
	"net/smtp"
	"strings"
)

type PosixAdapterOptions struct {
	SmtpAddr string
	From     string
	To       []string
}

func (o PosixAdapterOptions) init() (ad adapter, err error) {
	if o.SmtpAddr == "" {
		return nil,
			errors.New("posixAdapter: options 'SmtpAddr' must be not equal empty string")
	}

	if len(o.To) == 0 {
		return nil,
			errors.New("posixAdapter: options 'To' must be not equal empty array")
	}

	if o.From == "" {
		o.From = "notify@immo.ru"
	}

	return &posixAdapter{
		PosixAdapterOptions: o,
	}, nil
}

type posixAdapter struct {
	PosixAdapterOptions
	conn *smtp.Client
}

func (ad *posixAdapter) connect() (err error) {
	if ad.conn == nil {
		ad.conn, err = smtp.Dial(ad.SmtpAddr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (posixAdapter) Format(text []Line) []byte {
	size := 0
	for _, tx := range text {
		for _, element := range tx {
			switch t := element.(type) {
			case CODE:
				size += len(t) + 11
			case B:
				size += len(t) + 7
			case TextElement:
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
			case CODE:
				n += copy(complete[n:], "<pre>")
				n += copy(complete[n:], t)
				n += copy(complete[n:], "</pre>")
			case string:
				n += copy(complete[n:], t)
			case B:
				n += copy(complete[n:], "<b>")
				n += copy(complete[n:], t)
				n += copy(complete[n:], "</b>")
			case TextElement:
				n += copy(complete[n:], t)
			}
		}
		n += copy(complete[n:], "</br>")
	}

	return complete
}

func (ad *posixAdapter) send(subject string, message ...Line) error {
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	c, err := smtp.Dial(ad.SmtpAddr)
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.Mail(r.Replace(ad.From)); err != nil {
		return err
	}

	for i := range ad.To {
		ad.To[i] = r.Replace(ad.To[i])
		if err = c.Rcpt(ad.To[i]); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	msg := "To: " + strings.Join(ad.To, ",") + "\r\n" +
		"From: " + ad.From + "\r\n" +
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
