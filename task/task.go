package task

import (
  "context"
  "github.com/xpwu/timer/scheduler"
  "github.com/xpwu/timer/task/delay"
  "github.com/xpwu/timer/task/fixed"
)

type flag = byte

const (
  invalid flag = iota
  fixedF
  tickerF
  delayF
)

func NewFixedTask(f *fixed.Fixed) scheduler.Task {
  return append([]byte{fixedF}, f.ToBytes()...)
}

//func NewTickerTask(data Ticker) scheduler.Task {
//  return append([]byte{ticker}, data...)
//}

func NewDelayTask(d *delay.Delay) scheduler.Task {
  return append([]byte{delayF}, d.ToBytes()...)
}


type noopTask struct{}

func (n *noopTask) Run(ctx context.Context, runtime scheduler.UnixTimeSecond) {
}

func newTaskRunner(task scheduler.Task) taskRunner {
  switch task[0] {
  case fixedF:
    return fixed.FromBytes(task[1:])
  //case ticker:
  //  r := Ticker(task[1:])
  //  return &r
  case delayF:
    return delay.FromBytes(task[1:])
  default:
    return &noopTask{}
  }
}
