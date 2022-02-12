package fixed

type AddState byte

const (
  OK AddState = iota
  InternalError
  ConflictErr
  DataErr
)

type DB interface {
  // ConflictErr : 当id已经存在，但是cronTime不相同时，返回ConflictErr
  // InternalError : 其他因为DB等原因造成存储失败时，返回 InternalError
  // DataErr : 数据解析等数据相关的错误，返回DataErr
  // OK : 其余情况都应是OK，包括id存在且cronTime相同的情况也是OK
  //
  // opFlag: 只有首次添加id的数据时(包括Del(id)后的首次添加)，才更新DB中opFlag的值
  Add(id string, cronTime []byte, opFlag string) AddState

  Exist(ids []string) (notExist []string)
  Del(id string)
  Get(id string) (cronTime []byte, opFlag string, ok bool)

  Visit(startId string) (items []struct {
    Id       string
    CronTime []byte
    OpFlag   string
  },
    nextId string)
}

type noopDB struct {
}

func (n *noopDB) Add(id string, cronTime []byte, opFlag string) AddState {
  return InternalError
}

func (n *noopDB) Exist(ids []string) (notExist []string) {
  return ids
}

func (n *noopDB) Del(id string) {

}

func (n *noopDB) Get(id string) (cronTime []byte, opFlag string, ok bool) {
  return nil, "", false
}

func (n *noopDB) Visit(startId string) (
  items []struct {
  Id       string
  CronTime []byte
  OpFlag   string
},
  nextId string) {

  return []struct {
    Id       string
    CronTime []byte
    OpFlag   string
  }{}, ""
}

var db DB = &noopDB{}

func SetDB(d DB) {
  db = d
}
