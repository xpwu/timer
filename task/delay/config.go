package delay

import (
  "github.com/xpwu/go-config/configs"
  "github.com/xpwu/go-httpclient/httpc"
)

type config struct {
  CallbackUrl httpc.RawURL `conf:"url,callback url"`
}

var confValue = &config{}

func init() {
  configs.Unmarshal(confValue)
}
