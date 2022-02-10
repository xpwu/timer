package scheduler

import "context"

// unix 时间戳  unit: s
type UnixTimeSecond uint64

type DBIterator interface {
  First() bool
  Next() bool

  Release()

  TimeStamp() UnixTimeSecond
  Tasks() []Task
}

type DB interface {
  // [start, end)
  AllTasks(ctx context.Context, start, end UnixTimeSecond) DBIterator
  Delete(timestamp UnixTimeSecond)
  AppendTasks(timestamp UnixTimeSecond, tasks []Task)
}
