package server

import (
  "context"
  "github.com/xpwu/go-log/log"
  "github.com/xpwu/timer/task/fixed"
)

type delFixedRequest struct {
  Ids []string `json:"ids"`
}

func (s *suite) APIDelFixed(ctx context.Context, request *delFixedRequest) *noResponse {
  ctx, logger := log.WithCtx(ctx)

  for _, id := range request.Ids {
    logger.Debug("del ", id)
    fixed.Del(id)
  }

  return &noResponse{}
}
