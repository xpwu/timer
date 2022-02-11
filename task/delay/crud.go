package delay

import (
  "context"
  "github.com/xpwu/timer/scheduler"
  "time"
)

func Add(ctx context.Context, id string, d time.Duration) {
  t := scheduler.UnixTimeSecond(time.Now().Add(d).Unix())
  delay := &Delay{
    TryCount:  0,
    Id:        id,
    TimePoint: t,
  }

  scheduler.AddTask(t, []scheduler.Task{NewDelayTask(delay)})
}
