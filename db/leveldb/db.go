package leveldb

import (
  "github.com/xpwu/timer/scheduler"
  "github.com/xpwu/timer/task/fixed"
  "path/filepath"
)

func Init(root string) {
  root = filepath.Join(root, "db")
  fixed.SetDB(newFixed(root))
  scheduler.SetDB(newScheduler(root))
}
