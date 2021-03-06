package callback

import (
  "context"
  "encoding/json"
  "github.com/xpwu/go-httpclient/httpc"
  "github.com/xpwu/go-log/log"
  "github.com/xpwu/timer/scheduler"
  "net/http"
)

/**
1、id为任务添加时的id，透传
2、回调接收服务器需要满足幂等
 */
type Request struct {
  TimePoint scheduler.UnixTimeSecond `json:"time_point"`
  Id        string                   `json:"id"`
}

func Callback(ctx context.Context, rawUrl httpc.RawURL, data *Request) (ok bool) {
  ctx, logger := log.WithCtx(ctx)

  url := rawUrl.String()

  logger.PushPrefix("callback: url=" + url)
  d, err := json.Marshal(data)
  if err != nil {
    logger.Error(err)
    return false
  }

  logger.Debug("data=", string(d))

  var response *http.Response
  err = httpc.Send(ctx, url,
    httpc.WithBytesBody(d), httpc.WithResponse(&response))
  if err != nil {
    logger.Error(err)
    return false
  }

  if response.StatusCode != http.StatusOK {
    logger.Error(response.Status)
    return false
  }

  logger.Info("succeed")
  return true
}

var ReTryDuration = []scheduler.UnixTimeSecond{
  0, 5, 5, 5, 10, 10, 10, 30, 30, 30, 1 * 60, 1 * 60, 1 * 60, // 5min
  5 * 60, 5 * 60, 5 * 60, 5 * 60, 5 * 60, 5 * 60, 5 * 60, 5 * 60, 5 * 60, // 50min
  10 * 60, 10 * 60, 10 * 60, 10 * 60, 10 * 60, 10 * 60, 10 * 60, // 120min
  30 * 60, 30 * 60, 30 * 60, 30 * 60, 30 * 60, 30 * 60, 30 * 60, 30 * 60, // 6h
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 14h
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 1d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 2d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 3d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 4d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 5d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 6d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 7d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 8d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 9d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 10d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 11d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 12d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 13d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 14d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 15d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 16d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 17d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 18d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 19d

  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60,
  60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, 60 * 60, // 20d

}
