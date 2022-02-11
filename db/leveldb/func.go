package leveldb

func must(v interface{}, err error) interface{} {
  if err != nil {
    panic(err)
  }
  return v
}

func mustOk(err error) {
  if err != nil {
    panic(err)
  }
}

func mustOkOrFunc(err error, f func()) {
  if err != nil {
    f()
    panic(err)
  }
}
