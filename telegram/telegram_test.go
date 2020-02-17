package telegram

import (
	"testing"

	"github.com/farwydi/gotify"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestTelegramAdapter_Format(t *testing.T) {
	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	fixture := `{"ok": true, "result": {"id": 12345678, "first_name": "YourBot", "username": "YourBot"}}`
	responder := httpmock.NewStringResponder(200, fixture)
	fakeUrl := "https://api.telegram.org/bot1/getMe"
	httpmock.RegisterResponder("POST", fakeUrl, responder)

	init := NewTelegramAdapterWithHttp("1", []int64{1}, client)
	adapter, err := init()
	require.NoError(t, err)

	message := []gotify.Line{
		gotify.C("Hello *world*"),
		gotify.C(15),
		gotify.C(gotify.B("bold"), " text"),
		gotify.C(gotify.CODE("code here")),
	}

	require.Equal(t,
		[]byte("Hello \\*world\\*\n15\n*bold* text\n`code here`\n"),
		adapter.Format(message),
	)
}
