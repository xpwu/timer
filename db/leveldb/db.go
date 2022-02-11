package leveldb

import (
  "github.com/xpwu/timer/scheduler"
  "github.com/xpwu/timer/task/fixed"
)

func Init(root string) {
  fixed.SetDB(newFixed(root))
  scheduler.SetDB(newScheduler(root))
}
