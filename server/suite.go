package server

import (
  "github.com/xpwu/go-tinyserver/api"
)

type suite struct {
  api.PostJsonSetUpper
  api.PostJsonTearDowner
  api.RootURIMapper
}

func AddAPI() {
  api.Add(func() api.Suite {
    return &suite{}
  })
}

