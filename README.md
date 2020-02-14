# gotify
Send message anywhere

## Adapters

### Telegram
Setup with proxy
```go
package main

import (
    "github.com/farwydi/gotify"
    "github.com/farwydi/gotify/telegram"
    "github.com/go-resty/resty/v2"
)

func main()  {
    client, _ := gotify.NewClient(
        telegram.NewTelegramAdapterWithHttp(
            "token", 1414, resty.New().
                SetProxy("http://proxyserver:8888"),
        ),
    )
    client.Send("title")
}
```