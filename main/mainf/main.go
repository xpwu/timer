package mainf

import (
  "github.com/xpwu/go-cmd/arg"
  "github.com/xpwu/go-cmd/cmd"
  "github.com/xpwu/go-cmd/exe"
  "github.com/xpwu/go-config/configs"
  "github.com/xpwu/go-tinyserver/http"
  "github.com/xpwu/timer/db/leveldb"
  "github.com/xpwu/timer/scheduler"
  "github.com/xpwu/timer/server"
  "github.com/xpwu/timer/task"
  "path/filepath"
)

func Main() {
  cmd.RegisterCmd(cmd.DefaultCmdName, "start server", func(args *arg.Arg) {
    config := "config.json"
    args.String(&config, "c", "config file path")
    args.Parse()

    if !filepath.IsAbs(config) {
      config = filepath.Join(exe.Exe.AbsDir, config)
    }

    configs.SetConfigurator(&configs.JsonConfig{ReadFile: config})
    configs.Read()

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
