package leveldb

import (
  "github.com/syndtr/goleveldb/leveldb"
  "github.com/syndtr/goleveldb/leveldb/opt"
  "github.com/syndtr/goleveldb/leveldb/util"
  "github.com/xpwu/go-log/log"
  "github.com/xpwu/timer/task/fixed"
  "path"
  "reflect"
)

type fixedDB struct {
  db *leveldb.DB
}

func newFixed(root string) *fixedDB {
  p := path.Join(root, "fixed")
  db, err := leveldb.OpenFile(p, nil)
  if err != nil {
    panic(err)
  }
  return &fixedDB{db: db}
}

// ConflictErr : 当id已经存在，但是cronTime不相同时，返回ConflictErr
// InternalError : 其他因为DB等原因造成存储失败时，返回 InternalError
// DataErr : 数据解析等数据相关的错误，返回DataErr
// OK : 其余情况都应是OK，包括id存在且cronTime相同的情况也是OK
//
// opFlag: 只有首次添加id的数据时(包括Del(id)后的首次添加)，才更新DB中opFlag的值
func (f *fixedDB) Add(id string, cronTime []byte, opFlag string) (ret fixed.AddState) {
  ret = fixed.OK

  defer func() {
    if r := recover(); r != nil {
      log.Fatal(r)
      ret = fixed.InternalError
    }
  }()

  key := []byte(id)
  t := must(f.db.OpenTransaction()).(*leveldb.Transaction)
  v, err := t.Get(key, nil)

  if err == leveldb.ErrNotFound {
    err = t.Put(key, fixedToValue(cronTime, opFlag), &opt.WriteOptions{
      Sync: true,
    })
    if err != nil {
      t.Discard()
      ret = fixed.InternalError
      return
    }
  } else {
    oldCron, _ := valueToFixed(v)
    // todo equal
    if !reflect.DeepEqual(cronTime, oldCron) {
      ret = fixed.ConflictErr
    }
    t.Discard()
    return
  }

  err = t.Commit()
  if err != nil {
    // retry once
    err = t.Commit()
  }
  if err != nil {
    t.Discard()
    panic(err)
  }

  return
}

func (f *fixedDB) Exist(ids []string) (notExist []string) {
  ret := make([]string, len(ids))
  for _, s := range ids {
    if !must(f.db.Has([]byte(s), nil)).(bool) {
      ret = append(ret, s)
    }
  }

  return ret
}

func (f *fixedDB) Del(id string) {
  mustOk(f.db.Delete([]byte(id), &opt.WriteOptions{
    Sync: true,
  }))
}

func (f *fixedDB) Get(id string) (cronTime []byte, opFlag string, ok bool) {
  v, err := f.db.Get([]byte(id), nil)
  if err == leveldb.ErrNotFound {
    return nil, "", false
  }

  c, o := valueToFixed(v)

  return c, o, true
}

func (f *fixedDB) Visit(startId string) (
  items []struct {
  Id       string
  CronTime []byte
  OpFlag   string
},
  nextId string) {

  iter := f.db.NewIterator(&util.Range{
    Start: []byte(startId),
    Limit: nil,
  }, nil)

  // [startId, ...), len < 1000
  var sum int = 1000

  items = make([]struct {
    Id       string
    CronTime []byte
    OpFlag   string
  }, 0, sum)

  for iter.Next() && sum >= 0 {
    id := string(iter.Key())
    sum--

    c, o := valueToFixed(iter.Value())
    items = append(items, struct {
      Id       string
      CronTime []byte
      OpFlag   string
    }{Id: id, CronTime: c, OpFlag: o})
  }

  if sum < 0 {
    // 把最后一个作为下一次的起点，本次不返回最后一个
    nextId = items[len(items)-1].Id
    items = items[:len(items)-1]
  } else {
    nextId = ""
  }

  return
}
