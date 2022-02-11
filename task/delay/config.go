package delay

import "github.com/xpwu/go-config/configs"

type config struct {
  CallbackUrl string `conf:"url, callback url"`
}

var confValue = &config{}

func init() {
  configs.Unmarshal(confValue)
}
