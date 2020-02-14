package gotify

import (
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	os.Unsetenv("HTTP_PROXY")
	os.Unsetenv("HTTPS_PROXY")

	client, err := NewClient(
		TelegramAdapterOptions{
			Token:  "695845816:AAEH665gpYS9PBoX29xQ6TsTOOLncjNHlGI",
			ChatId: -369034152,
			// Привет HelpDesk
			HttpClient: &http.Client{
				Transport: &http.Transport{
					Proxy: func(request *http.Request) (*url.URL, error) {
						return url.Parse("http://192.168.7.32:8888/")
					},
				},
			},
		},
		PosixAdapterOptions{
			SmtpAddr: "vps3232.mtu.immo:25",
			To:       []string{"zharikov@immo.ru"},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Send("Hello", []Line{
		C("Hello world"),
		C(B("bold"), " text"),
	})
	if err != nil {
		t.Fatal(err)
	}
}
