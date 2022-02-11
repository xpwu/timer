package server

import (
  "context"
  "github.com/xpwu/timer/task/delay"
  "time"
)

type addDelayRequest struct {
  // unit: s
  Duration time.Duration `json:"d"`
  Id       string        `json:"id"`
}

func (s *suite) APIAddDelay(ctx context.Context, request *addDelayRequest) *noResponse {

  delay.Add(ctx, request.Id, request.Duration*time.Second)

  return &noResponse{}
}
