package server

import (
  "context"
  "github.com/xpwu/timer/task/fixed"
)

type existFixedRequest struct {
  Ids []string `json:"ids"`
}

type existFixedResponse struct {
  Ids []string `json:"ids"`
}

func (s *suite) APIExistFixed(ctx context.Context, request *existFixedRequest) *existFixedResponse {
  res := &existFixedResponse{}

  res.Ids = fixed.Exist(request.Ids)

  return res
}
