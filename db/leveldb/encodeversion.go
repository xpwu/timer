package leveldb

import (
  "encoding/json"
  "github.com/xpwu/timer/scheduler"
)

type encodeVersion = byte

const (
  jsonV encodeVersion = iota
)

// version | value
// 暂时仅仅支持json方式

func valueToTasks(v []byte) []scheduler.Task {
  if v[0] != jsonV {
    panic("valueToTasks version not support")
  }

  var ret []scheduler.Task
  mustOk(json.Unmarshal(v[1:], &ret))

  return ret
}

func tasksToValue(tasks []scheduler.Task) []byte {
  v, err := json.Marshal(tasks)
  if err != nil {
    panic("tasksToValue error")
  }

  return append([]byte{jsonV}, v...)
}

type fixE struct {
  CronTime []byte `json:"c"`
  OpFlag   string `json:"o"`
}

func fixedToValue(cronTime []byte, opFlag string) []byte {
  f := &fixE{
    CronTime: cronTime,
    OpFlag:   opFlag,
  }

  v := must(json.Marshal(f)).([]byte)

  return append([]byte{jsonV}, v...)
}

func valueToFixed(v []byte) (cronTime []byte, opFlag string) {
  if v[0] != jsonV {
    panic("valueToTasks version not support")
  }

  ret := &fixE{}
  mustOk(json.Unmarshal(v[1:], &ret))

  return ret.CronTime, ret.OpFlag
}
