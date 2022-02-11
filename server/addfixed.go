package server

import (
  "context"
  "github.com/syndtr/goleveldb/leveldb/errors"
  "github.com/xpwu/timer/scheduler"
  "github.com/xpwu/timer/task/fixed"
)

type addFixedRequest struct {
  // unit: s
  StartTime uint64 `json:"start"`
  CronTime  string `json:"cron"`
  Id        string `json:"id"`
}

type fixedResponseStatus = byte

const (
  OK fixedResponseStatus = iota
  IdConflict
)

type addFixedResponse struct {
  Status fixedResponseStatus `json:"status"`
}

func (s *suite) APIAddFixed(ctx context.Context, request *addFixedRequest) *addFixedResponse {
  //ctx, logger := log.WithCtx(ctx)

  res := &addFixedResponse{Status: OK}

  st := fixed.Add(ctx, request.Id, request.CronTime, scheduler.UnixTimeSecond(request.StartTime))

  if st != fixed.OK && st != fixed.ConflictErr {
    s.Req.Terminate(errors.New("db error"))
  }

  if st == fixed.ConflictErr {
    res.Status = IdConflict
  }

  return res
}


