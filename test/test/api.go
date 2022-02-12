package test

import (
  "context"
  "fmt"
  "github.com/xpwu/go-log/log"
  "github.com/xpwu/go-tinyserver/api"
  "github.com/xpwu/timer/task/callback"
  "time"
)

type suite struct {
  api.PostJsonSetUpper
  api.PostJsonTearDowner
}

func (s *suite) MappingPreUri() string {
  return "/test"
}

func (s *suite) APIDelay(ctx context.Context, request *callback.Request) *api.EmptyResponse {
  _, logger := log.WithCtx(ctx)

  t := time.Unix(int64(request.TimePoint), 0)

  logger.Info(fmt.Sprintf("timepoint=%v, id=%s", t, request.Id))

  return &api.EmptyResponse{}
}

func AddAPI() {
  api.Add(func() api.Suite {
    return &suite{}
  })
}
