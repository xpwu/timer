package main

import (
  "github.com/xpwu/go-cmd/arg"
  "github.com/xpwu/go-cmd/cmd"
  _ "github.com/xpwu/go-cmd/cmd/printconf"
  "github.com/xpwu/go-cmd/exe"
  "github.com/xpwu/go-config/configs"
  "github.com/xpwu/go-tinyserver/http"
  "github.com/xpwu/timer/db/leveldb"
  "github.com/xpwu/timer/scheduler"
  "github.com/xpwu/timer/server"
  "path/filepath"
  "time"
)

func main() {
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

    down := make(chan struct{}, 1)
    scheduler.Start(down)

    server.AddAPI()
    http.Start()

    for {
      select {
      case <-down:
        time.Sleep(5 * time.Second)
        scheduler.Start(down)
      }
    }
  })

  cmd.Run()
}
