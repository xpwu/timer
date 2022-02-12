package scheduler

import (
  "context"
  "github.com/xpwu/go-log/log"
)

var (
	runner Runner = &noopRunner{}
	db DB = &noopDB{}
)

type noopRunner struct {}

func (n *noopRunner) Run(ctx context.Context, schedulerTime UnixTimeSecond, tasks []Task) {
  _,logger := log.WithCtx(ctx)
  logger.Warning("no runner")
}

func SetRunner(r Runner) {
  runner = r
}

func getRunner() Runner {
  return runner
}

type noopDBIter struct {}

func (n *noopDBIter) First() bool {
  return false
}

func (n *noopDBIter) Next() bool {
  return false
}

func (n *noopDBIter) Release() {
}

func (n *noopDBIter) TimeStamp() UnixTimeSecond {
  return 0
}

func (n *noopDBIter) Tasks() []Task {
  return nil
}

type noopDB struct {}

func (n *noopDB) AllTasks(ctx context.Context, start, end UnixTimeSecond) DBIterator {
  return &noopDBIter{}
}

func (n *noopDB) Delete(timestamp UnixTimeSecond) {
}

func (n *noopDB) AppendTasks(timestamp UnixTimeSecond, tasks []Task) {

}

func SetDB(d DB) {
  db = d
}

func getDB() DB {
  return db
}



