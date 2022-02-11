package fixed

import (
  "context"
  "github.com/robfig/cron/v3"
  "github.com/xpwu/go-log/log"
  "github.com/xpwu/go-reqid/reqid"
  "github.com/xpwu/timer/scheduler"
  "time"
)

/**

1、多次成功add同一条
Add(a, c, d)    (succeed)
Add(a, c, d)    (succeed)
Add(a, c, d)    (succeed)
Add(a, c, d)    (succeed)
Add(a, c, d)    (succeed)
后面的几次操作不会改变timer的状态，按照第一次成功后的状态执行


2、失败后多次成功add同一条
Add(a, c, d)    (failed)
Add(a, c, d)    (failed)

Add(a, c, d)    (succeed)
Add(a, c, d)    (succeed)
Add(a, c, d)    (succeed)
Add(a, c, d)    (succeed)
失败add 不会改变timer的状态，系统中不会有该timer的任何执行结果
除第一条成功的外，后面的几次成功操作不会改变timer的状态，按照第一次成功后的状态执行


3、成功或者失败添加、再成功删除、再成功添加
Add(a, c, d)    (succeed/failed)
Del(a)          (succeed)
Add(a, c, d)    (succeed)
Del后，之前的状态都清空；重新从时间d开始执行a

*/

func Add(ctx context.Context, id string, cronTimeStr string, start scheduler.UnixTimeSecond) AddState {
  ctx, logger := log.WithCtx(ctx)
  logger.PushPrefix("fixed add")

  sche, err := cron.ParseStandard(cronTimeStr)
  if err != nil {
    logger.Error(err)
    return DataErr
  }
  spec, ok := sche.(*cron.SpecSchedule)
  if !ok {
    logger.Error("cron.ParseStandard().(*cron.SpecSchedule) error")
    return InternalError
  }

  cronTime := NewCronTimeFromSpec(spec, start)
  opFlag := reqid.RandomID()
  next := scheduler.UnixTimeSecond(cronTime.Next(time.Unix(int64(start), 0)).Unix())
  fixed := &Fixed{
    TryCount:  0,
    Id:        id,
    TimePoint: next,
    OpFlag:    opFlag,
  }
  tk := NewFixedTask(fixed)

  // 先添加scheduler, scheduler重复添加也没关系，在runner的时候做最终的检测
  scheduler.AddTask(next, []scheduler.Task{tk})
  // 再写DB，要保证原子性，不能read-do-write
  return db.Add(id, cronTime.ToBytes(), opFlag)
}

func Del(id string) {
  db.Del(id)
}

func Exist(ids []string) (notExist []string) {
  return db.Exist(ids)
}
