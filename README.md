[![Build Status](https://travis-ci.org/farwydi/gotify.svg?branch=master)](https://travis-ci.org/farwydi/gotify)
[![codecov](https://codecov.io/gh/farwydi/gotify/branch/master/graph/badge.svg)](https://codecov.io/gh/farwydi/gotify)

# gotify
Send message anywhere

## Install
`go get github.com/farwydi/gotify`

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