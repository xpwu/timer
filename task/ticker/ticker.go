package ticker

import (
  "context"
  "github.com/xpwu/timer/scheduler"
  "time"
)

type Ticker []byte

func (t *Ticker) Run(ctx context.Context, schedulerTime scheduler.UnixTimeSecond) {

}

func AddTicker(d time.Duration, id string) {

}

func DelTicker(id string) {

}
