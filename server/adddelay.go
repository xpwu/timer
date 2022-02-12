package server

import (
  "context"
  "github.com/xpwu/go-tinyserver/api"
  "github.com/xpwu/timer/task/delay"
  "time"
)

/**
无论系统中是否存在Id任务，都是直接添加进系统中
 */
type addDelayRequest struct {
  // unit: s
  Duration time.Duration `json:"d"`
  Id       string        `json:"id"`
}

func (s *suite) APIAddDelay(ctx context.Context, request *addDelayRequest) *api.EmptyResponse {

  delay.Add(ctx, request.Id, request.Duration*time.Second)

  return &api.EmptyResponse{}
}
