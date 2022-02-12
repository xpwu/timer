package server

import (
  "context"
  "github.com/xpwu/go-log/log"
  "github.com/xpwu/go-tinyserver/api"
  "github.com/xpwu/timer/task/fixed"
)

type delFixedRequest struct {
  Ids []string `json:"ids"`
}

func (s *suite) APIDelFixed(ctx context.Context, request *delFixedRequest) *api.EmptyResponse {
  ctx, logger := log.WithCtx(ctx)

  for _, id := range request.Ids {
    logger.Debug("del ", id)
    fixed.Del(id)
  }

  return &api.EmptyResponse{}
}
