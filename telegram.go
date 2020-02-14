package gotify

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"net/http"
	"strings"
)

type TelegramAdapterOptions struct {
	ChatId     int64
	Token      string
	HttpClient *http.Client
}

func (o TelegramAdapterOptions) init() (ad adapter, err error) {
	if o.HttpClient == nil {
		o.HttpClient = http.DefaultClient
	}

	if o.Token == "" {
		return nil,
			errors.New("telegramAdapter: options 'Token' must be not equal empty string")
	}

	if o.ChatId == 0 {
		return nil,
			errors.New("telegramAdapter: options 'ChatId' must be not equal 0")
	}

	bot, err := tgbotapi.NewBotAPIWithClient(
		o.Token,
		o.HttpClient,
	)

	return &telegramAdapter{
		BotAPI:                 bot,
		TelegramAdapterOptions: o,
	}, err
}

type telegramAdapter struct {
	*tgbotapi.BotAPI
	TelegramAdapterOptions
}

func (telegramAdapter) escaped(t []byte) []byte {
	t = bytes.Replace(t, []byte("_"), []byte("\\_"), -1)
	t = bytes.Replace(t, []byte("*"), []byte("\\*"), -1)
	t = bytes.Replace(t, []byte("["), []byte("\\["), -1)
	t = bytes.Replace(t, []byte("`"), []byte("\\`"), -1)
	return t
}

func (telegramAdapter) escapedStr(t string) string {
	t = strings.Replace(t, "_", "\\_", -1)
	t = strings.Replace(t, "*", "\\*", -1)
	t = strings.Replace(t, "[", "\\[", -1)
	t = strings.Replace(t, "`", "\\`", -1)
	return t
}

func (tg *telegramAdapter) Format(text []Line) []byte {
	for i, tx := range text {
		for j, element := range tx {
			switch t := element.(type) {
			case B:
				text[i][j] = B(tg.escaped(t))
			case CODE:
				text[i][j] = CODE(bytes.Replace(t, []byte("`"), []byte("\\`"), -1))
			case TextElement:
				text[i][j] = TextElement(tg.escaped(t))
			case string:
				text[i][j] = tg.escapedStr(t)
			case []byte:
				text[i][j] = tg.escaped(t)
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				text[i][j] = fmt.Sprintf("%d", t)
			default:
				// TODO: Добавить обработку ошибки типа
				return nil
			}
		}
	}

	size := 0
	for _, tx := range text {
		for _, element := range tx {
			switch t := element.(type) {
			case B:
				size += len(t) + 2
			case CODE:
				size += len(t) + 2
			case TextElement:
				size += len(t)
			case string:
				size += len(t)
			case []byte:
				size += len(t)
			}
		}
		size += 1
	}

	complete := make([]byte, size)
	n := 0
	for _, tx := range text {
		for _, element := range tx {
			switch t := element.(type) {
			case string:
				n += copy(complete[n:], t)
			case B:
				n += copy(complete[n:], "*")
				n += copy(complete[n:], t)
				n += copy(complete[n:], "*")
			case CODE:
				n += copy(complete[n:], "`")
				n += copy(complete[n:], t)
				n += copy(complete[n:], "`")
			case TextElement:
				n += copy(complete[n:], t)
			}
		}
		complete[n] = byte('\n')
		n += 1
	}

	return complete
}

func (ad *telegramAdapter) send(subject string, message ...Line) error {
	// Компиляция
	msg := fmt.Sprintf("*%s*\n%s", subject, ad.Format(message))

	_, err := ad.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           ad.ChatId,
			ReplyToMessageID: 0,
		},
		Text:                  msg,
		ParseMode:             "markdown",
		DisableWebPagePreview: false,
	})
	return err
}
