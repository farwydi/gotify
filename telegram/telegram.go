package telegram

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/Syfaro/telegram-bot-api"
	"github.com/farwydi/gotify"
	"github.com/go-resty/resty/v2"
)

func NewTelegramAdapter(token string, chatIds []int64) func() (gotify.Adapter, error) {
	return NewTelegramAdapterWithHttp(token, chatIds, resty.New())
}

func NewTelegramAdapterWithHttp(token string, chatIds []int64, client *resty.Client) func() (gotify.Adapter, error) {
	if token == "" {
		return func() (gotify.Adapter, error) {
			return nil,
				errors.New("telegramAdapter: options 'Token' must be not equal empty string")
		}
	}

	if len(chatIds) == 0 {
		return func() (gotify.Adapter, error) {
			return nil,
				errors.New("telegramAdapter: options 'ChatId' must be not empty")
		}
	}

	bot, err := tgbotapi.NewBotAPIWithClient(
		token,
		client.GetClient(),
	)
	if err != nil {
		return func() (gotify.Adapter, error) {
			return nil, err
		}
	}

	return func() (gotify.Adapter, error) {
		return &telegramAdapter{
			tgApi:   bot,
			chatIds: chatIds,
			token:   token,
		}, nil
	}
}

type telegramAdapter struct {
	tgApi   *tgbotapi.BotAPI
	chatIds []int64
	token   string
}

func escaped(t []byte) []byte {
	t = bytes.ReplaceAll(t, []byte("_"), []byte("\\_"))
	t = bytes.ReplaceAll(t, []byte("*"), []byte("\\*"))
	t = bytes.ReplaceAll(t, []byte("["), []byte("\\["))
	t = bytes.ReplaceAll(t, []byte("`"), []byte("\\`"))
	return t
}

func escapedStr(t string) string {
	t = strings.ReplaceAll(t, "_", "\\_")
	t = strings.ReplaceAll(t, "*", "\\*")
	t = strings.ReplaceAll(t, "[", "\\[")
	t = strings.ReplaceAll(t, "`", "\\`")
	return t
}

func (telegramAdapter) Format(text []gotify.Line) []byte {
	for i, tx := range text {
		for j, element := range tx {
			switch t := element.(type) {
			case gotify.B:
				text[i][j] = gotify.B(escaped(t))
			case gotify.CODE:
				text[i][j] = gotify.CODE(bytes.ReplaceAll(t, []byte("`"), []byte("\\`")))
			case gotify.TextElement:
				text[i][j] = gotify.TextElement(escaped(t))
			case string:
				text[i][j] = escapedStr(t)
			case []byte:
				text[i][j] = escaped(t)
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
			case gotify.B:
				size += len(t) + 2
			case gotify.CODE:
				size += len(t) + 2
			case gotify.TextElement:
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
			case gotify.B:
				n += copy(complete[n:], "*")
				n += copy(complete[n:], t)
				n += copy(complete[n:], "*")
			case gotify.CODE:
				n += copy(complete[n:], "`")
				n += copy(complete[n:], t)
				n += copy(complete[n:], "`")
			case gotify.TextElement:
				n += copy(complete[n:], t)
			}
		}
		complete[n] = byte('\n')
		n += 1
	}

	return complete
}

func (ad telegramAdapter) Send(subject string, message ...gotify.Line) error {
	// Компиляция
	msg := fmt.Sprintf("*%s*\n%s", subject, ad.Format(message))

	for _, chatId := range ad.chatIds {
		_, err := ad.tgApi.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: chatId,
			},
			Text:      msg,
			ParseMode: "markdown",
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (ad telegramAdapter) SendRaw(p []byte) error {
	_, err := ad.tgApi.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: ad.chatId,
		},
		Text:                  string(p),
		ParseMode:             "markdown",
		DisableWebPagePreview: true,
	})
	if err != nil {
		return err
	}

	return nil
}
