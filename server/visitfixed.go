package server

import (
  "context"
  "github.com/xpwu/timer/task/fixed"
)

type visitFixedRequest struct {
  // "" 表示从头开始
  StartId string `json:"id"`
}

type visitFixedResponse struct {
  Results []addFixedRequest `json:"results"`
  // "" 表示 end
  NextId string `json:"nextId"`
}

func (s *suite) APIVisitFixed(ctx context.Context, request *visitFixedRequest) *visitFixedResponse {
  items, next := fixed.Visit(request.StartId)
  ret := &visitFixedResponse{
    Results: make([]addFixedRequest, 0, len(items)),
    NextId:  next,
  }

  for _, item := range items {
    ret.Results = append(ret.Results, addFixedRequest{
      StartTime: item.Start,
      CronTime:  item.CronTimeStr,
      Id:        item.Id,
    })
  }

  return ret
}
