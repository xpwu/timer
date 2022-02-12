package server

import (
  "context"
  "github.com/syndtr/goleveldb/leveldb/errors"
  "github.com/xpwu/timer/scheduler"
  "github.com/xpwu/timer/task/fixed"
  "time"
)

/**
1、所有数据都相同的重复添加，不会对系统有任何的状态改变，并不会从StartTime重新执行此任务，并且返回OK
2、任何新添加的任务(此次添加前系统中不存在)，都是从StartTime开始执行此任务，返回OK
3、Id相同，但是其他数据不同的任务，添加失败，返回IdConflict
 */

type addFixedRequest struct {
  // unit: s; 0:从收到请求开始算
  StartTime scheduler.UnixTimeSecond `json:"start"`
  CronTime  string                   `json:"cron"`
  Id        string                   `json:"id"`
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

  start := request.StartTime
  if start == 0 {
    start = scheduler.UnixTimeSecond(time.Now().Unix())
  }

  st := fixed.Add(ctx, request.Id, request.CronTime, start)

  if st != fixed.OK && st != fixed.ConflictErr {
    s.Req.Terminate(errors.New("db error"))
  }

  if st == fixed.ConflictErr {
    res.Status = IdConflict
  }

  return res
}
