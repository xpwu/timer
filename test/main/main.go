package main

import (
  _ "github.com/xpwu/go-cmd/cmd/printconf"
  "github.com/xpwu/timer/main/mainf"
  "github.com/xpwu/timer/test/test"
)

func main() {
  test.AddAPI()

  mainf.Main()
}
