package fixed

import (
  "context"
  "encoding/json"
  "fmt"
  "github.com/robfig/cron/v3"
  "github.com/xpwu/go-log/log"
  "github.com/xpwu/timer/scheduler"
  "github.com/xpwu/timer/task/callback"
  "github.com/xpwu/timer/task/flag"
  "time"
)

type Fixed struct {
  TryCount  uint16                   `json:"try"`
  Id        string                   `json:"id"`
  TimePoint scheduler.UnixTimeSecond `json:"tp"`
  OpFlag    string                   `json:"op"`
}

func NewFixedTask(f *Fixed) scheduler.Task {
  return append([]byte{flag.Fixed}, f.ToBytes()...)
}

func (f *Fixed) Run(ctx context.Context, schedulerTime scheduler.UnixTimeSecond) {
  ctx, logger := log.WithCtx(ctx)
  logger.PushPrefix(fmt.Sprintf("run fixed. id=%s, timepoint=%d", f.Id, f.TimePoint))

  cronTimeB, opF, ok := db.Get(f.Id)
  // 已经删除或者OpFlag不相同的task都不真正的执行
  if !ok || opF != f.OpFlag {
    return
  }

  // 只有非重试的情况下，才添加下一次的scheduler
  if f.TryCount == 0 {
    cronTime := NewCronTimeFromBytes(cronTimeB)
    // 增加一个小的偏移，以防端点处的bug
    next := scheduler.UnixTimeSecond(cronTime.Next(time.Unix(int64(f.TimePoint), 1000)).Unix())
    fixed := &Fixed{
      TryCount:  0,
      Id:        f.Id,
      TimePoint: next,
      OpFlag:    f.OpFlag,
    }
    tk := NewFixedTask(fixed)
    scheduler.AddTask(next, []scheduler.Task{tk})
  }

  req := &callback.Request{
    TimePoint: f.TimePoint,
    Id:        f.Id,
  }

  ok = callback.Callback(ctx, confValue.CallbackUrl, req)
  if ok {
    return
  }

  // retry, 超过最大重试时间，直接放弃
  if int(f.TryCount) >= len(callback.ReTryDuration)-1 {
    return
  }

  tc := f.TryCount + 1
  // 从当前时间计算下次重试的时间
  next := scheduler.UnixTimeSecond(time.Now().Unix()) + callback.ReTryDuration[tc]
  newF := &Fixed{
    TryCount:  tc,
    Id:        f.Id,
    TimePoint: f.TimePoint,
    OpFlag:    f.OpFlag,
  }
  scheduler.AddTask(next, []scheduler.Task{NewFixedTask(newF)})
}

func (f *Fixed) ToBytes() []byte {
  r, err := json.Marshal(f)
  if err != nil {
    panic(err)
  }

  return r
}

func FromBytes(b []byte) *Fixed {
  f := &Fixed{}
  _ = json.Unmarshal(b, f)
  return f
}

type CronTime struct {
  cron.SpecSchedule
  // json:"Location" 覆盖嵌套的 Location 域
  LocationStr string `json:"Location"`
  StartTime   scheduler.UnixTimeSecond
}

func NewCronTimeFromSpec(s *cron.SpecSchedule, start scheduler.UnixTimeSecond) *CronTime {
  return &CronTime{
    SpecSchedule: *s,
    LocationStr:  s.Location.String(),
    StartTime:    start,
  }
}

func NewCronTimeFromBytes(j []byte) *CronTime {
  ret := &CronTime{}
  _ = json.Unmarshal(j, ret)
  ret.Location, _ = time.LoadLocation(ret.LocationStr)
  return ret
}

func (c *CronTime) ToBytes() []byte {
  d, _ := json.Marshal(c)
  return d
}
