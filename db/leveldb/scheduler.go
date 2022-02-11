package leveldb

import (
  "context"
  "encoding/binary"
  "github.com/syndtr/goleveldb/leveldb"
  "github.com/syndtr/goleveldb/leveldb/iterator"
  "github.com/syndtr/goleveldb/leveldb/opt"
  "github.com/syndtr/goleveldb/leveldb/util"
  "github.com/xpwu/timer/scheduler"
  "path"
)

// 根据leveldb的按照byte比较的方式，这里使用大端方式编码即可满足需求
func unixTimeSecondToKey(t scheduler.UnixTimeSecond) []byte {
  r := make([]byte, 8)
  binary.BigEndian.PutUint64(r, uint64(t))
  return r
}

func keyToUnixTimeSecond(key []byte) scheduler.UnixTimeSecond {
  r := binary.BigEndian.Uint64(key)
  return scheduler.UnixTimeSecond(r)
}

type schIter struct {
  iterator.Iterator
}

func (s *schIter) TimeStamp() scheduler.UnixTimeSecond {
  return keyToUnixTimeSecond(s.Key())
}

func (s *schIter) Tasks() []scheduler.Task {
  return valueToTasks(s.Value())
}

type schedulerDB struct {
  db *leveldb.DB
}

func newScheduler(root string) *schedulerDB {
  p := path.Join(root, "scheduler")
  db,err := leveldb.OpenFile(p, nil)
  if err != nil {
    panic(err)
  }
  return &schedulerDB{db: db}
}

func (s *schedulerDB) AllTasks(ctx context.Context, start, end scheduler.UnixTimeSecond) scheduler.DBIterator {
  i := s.db.NewIterator(&util.Range{
    Start: unixTimeSecondToKey(start),
    Limit: unixTimeSecondToKey(end),
  }, nil)

  return &schIter{i}
}

func (s *schedulerDB) Delete(timestamp scheduler.UnixTimeSecond) {
  _ = s.db.Delete(unixTimeSecondToKey(timestamp), nil)
}

func (s *schedulerDB) AppendTasks(timestamp scheduler.UnixTimeSecond, tasks []scheduler.Task) {
  key := unixTimeSecondToKey(timestamp)
  t := must(s.db.OpenTransaction()).(*leveldb.Transaction)
  old := make([]scheduler.Task, 0)

  ok, err := t.Has(key, nil)
  if err != nil {
    t.Discard()
    panic(err)
  }
  if ok {
    v, err := t.Get(key, nil)
    if err != nil {
      t.Discard()
      panic(err)
    }
    old = valueToTasks(v)
  }

  newV := tasksToValue(scheduler.MergeTasks(old, tasks))
  mustOkOrFunc(t.Put(key, newV, &opt.WriteOptions{
    Sync: true,
  }), func() {
    t.Discard()
  })

  err = t.Commit()
  if err != nil {
    // retry once
    err = t.Commit()
  }
  if err != nil {
    t.Discard()
    panic(err)
  }
}

