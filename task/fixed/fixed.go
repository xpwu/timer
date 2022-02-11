package fixed

import (
  "context"
  "encoding/json"
  "github.com/robfig/cron/v3"
  "github.com/xpwu/timer/scheduler"
  "time"
)

type Fixed struct {
  IsTry     bool                     `json:"try"`
  Id        string                   `json:"id"`
  TimePoint scheduler.UnixTimeSecond `json:"tp"`
  OpFlag    string                   `json:"op"`
}

func (f *Fixed) Run(ctx context.Context, schedulerTime scheduler.UnixTimeSecond) {

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
