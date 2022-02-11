package delay

import (
  "context"
  "encoding/json"
  "github.com/xpwu/timer/scheduler"
  "time"
)

type Delay []byte

func (d *Delay) Run(ctx context.Context, schedulerTime scheduler.UnixTimeSecond) {

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

func AddDelay(d time.Duration, id string) {

}
