package mainf

import (
  "github.com/xpwu/go-cmd/arg"
  "github.com/xpwu/go-cmd/cmd"
  "github.com/xpwu/go-cmd/exe"
  "github.com/xpwu/go-tinyserver/http"
  "github.com/xpwu/timer/db/leveldb"
  "github.com/xpwu/timer/scheduler"
  "github.com/xpwu/timer/server"
  "github.com/xpwu/timer/task"
)

func Main() {
  cmd.RegisterCmd(cmd.DefaultCmdName, "start server", func(args *arg.Arg) {
    arg.ReadConfig(args)
    args.Parse()

    root := exe.Exe.AbsDir
    leveldb.Init(root)

    task.InitRunner()

    scheduler.Start()

    server.AddAPI()
    http.Start()

    // block
    block := make(chan struct{})
    <-block
  })

  cmd.Run()
}
