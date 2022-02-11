package delay

import (
  "context"
  "encoding/json"
  "fmt"
  "github.com/xpwu/go-log/log"
  "github.com/xpwu/timer/scheduler"
  "github.com/xpwu/timer/task"
  "time"
)

type Delay struct {
  TryCount  uint16                   `json:"try"`
  Id        string                   `json:"id"`
  TimePoint scheduler.UnixTimeSecond `json:"tp"`
}

type callback struct {
  TimePoint scheduler.UnixTimeSecond `json:"time_point"`
  Id        string                   `json:"id"`
}

func (d *Delay) Run(ctx context.Context, schedulerTime scheduler.UnixTimeSecond) {
  ctx, logger := log.WithCtx(ctx)
  logger.PushPrefix(fmt.Sprintf("run delay. id=%s, timepoint=%d", d.Id, d.TimePoint))

  req := &callback{
    TimePoint: d.TimePoint,
    Id:        d.Id,
  }

  ok := task.Callback(ctx, confValue.CallbackUrl, req)
  if ok {
    return
  }

  // retry, 超过最大重试时间，直接放弃
  if int(d.TryCount) >= len(task.ReTryDuration)-1 {
    return
  }

  tc := d.TryCount + 1
  // 从当前时间计算下次重试的时间
  next := scheduler.UnixTimeSecond(time.Now().Unix()) + task.ReTryDuration[tc]
  newD := &Delay{
    TryCount:  tc,
    Id:        d.Id,
    TimePoint: d.TimePoint,
  }
  scheduler.AddTask(next, []scheduler.Task{task.NewDelayTask(newD)})
}

func (d *Delay) ToBytes() []byte {
  r,err := json.Marshal(d)
  if err != nil {
    panic(err)
  }

  return r
}

func FromBytes(b []byte) *Delay {
  d := &Delay{}
  _ = json.Unmarshal(b, d)
  return d
}

