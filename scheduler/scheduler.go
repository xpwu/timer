package scheduler

import (
  "context"
  "encoding/hex"
  "fmt"
  "github.com/xpwu/go-log/log"
  "time"
)

type Task []byte

func MergeTasks(ones, others []Task) []Task {
  // task 内容完全相同的，认为是相同的task
  ret := ones
  m := make(map[string]bool)
  for _,t := range ones {
    m[hex.EncodeToString(t)] = true
  }
  for _,t := range others {
    if !m[hex.EncodeToString(t)] {
      ret = append(ret, t)
    }
  }

  return ret
}

type tsTask struct {
  ts   UnixTimeSecond
  task []Task
}

var (
  newTask = make(chan tsTask, 100)
)

func AddTask(time UnixTimeSecond, tasks []Task) {
  getDB().AppendTasks(time, tasks)

  newTask <- tsTask{
    ts:   time,
    task: tasks,
  }
}

type overD struct {
  ts      UnixTimeSecond
  crashed bool
  tasks   []Task
}

func Start() {
  go func() {
    down := make(chan struct{}, 1)

    for {
      start(context.Background(), down)

      select {
      case <-down:
        time.Sleep(5 * time.Second)
      }
    }
  }()
}

func start(ctx context.Context, down chan<- struct{}) {
  ctx, cancel := context.WithCancel(ctx)
  ctx, logger := log.WithCtx(ctx)
  logger.PushPrefix("scheduler")

  logger.Info("start")

  defer func() {
    if r := recover(); r != nil {
      logger.Fatal(r)
    }
    cancel()

    down <- struct{}{}
  }()

  // start 的加入是为了防止运行中的任务重复被调度(但是不会真正的执行)，
  // 但不能影响到正常的调度的最高原则：宁可重调，也不漏调
  start := UnixTimeSecond(0)

  over := make(chan overD, 1000)
  running := make(map[UnixTimeSecond]bool)

  for {
    logger.Info("find task")

    now := UnixTimeSecond(time.Now().Unix())
    // [start, end) now+1才能找到now
    end := now + 1
    tasks := getDB().AllTasks(ctx, start, end)

    for tasks.Next() {
      logger.Debug(log.LazyMsg(func() string {
        return fmt.Sprintf("find tasks(time=%v, unixts=%d, sum=%d) between %v and %v",
          time.Unix(int64(tasks.TimeStamp()), 0), tasks.TimeStamp(), len(tasks.Tasks()),
          time.Unix(int64(start), 0), time.Unix(int64(end), 0))
      }))

      if len(tasks.Tasks()) == 0 {
        log.Warning(fmt.Sprintf("0 task at unixts(%d)", tasks.TimeStamp()))
        continue
      }

      if running[tasks.TimeStamp()] {
        log.Debug(log.LazyMsg(func() string {
          return fmt.Sprintf("scheduler: find task(time=%v) has run",
            time.Unix(int64(tasks.TimeStamp()), 0))
        }))
        continue
      }

      running[tasks.TimeStamp()] = true
      go func() {
        ctx, logger := log.WithCtx(ctx)
        logger.PushPrefix("runner")
        crashed := false

        defer func() {
          if r := recover(); r != nil {
            logger.Fatal(r)
            crashed = true
          }
          over <- overD{
            ts:      tasks.TimeStamp(),
            crashed: crashed,
            tasks:   tasks.Tasks(),
          }
        }()

        getRunner().Run(ctx, tasks.TimeStamp(), tasks.Tasks())
      }()
    }
    tasks.Release()

    // 设置start, 为了防止未知问题，造成漏调度，当没有运行中的任务时，重置start
    if len(running) == 0 {
      start = UnixTimeSecond(0)
      logger.Debug("find task over, len(running) == 0, move start time to ", start)
    } else {
      start = end
      logger.Debug("find task over, move start time to ", start)
    }

    // 小于等于 now 的都已经取出执行了，所以这里从 end 开始找下一个
    // 否则可能找出的是正在执行的任务，则会空循环，浪费资源
    // 找5分钟范围内的第一个
    next := now + 5*60
    tasks = getDB().AllTasks(ctx, end, next)
    if tasks.First() {
      next = tasks.TimeStamp()
    }
    tasks.Release()
    sleep := time.NewTimer((time.Duration(next - now)) * time.Second)

  forSleep:
    for {
      logger.Info(fmt.Sprintf("will sleep(%ds)", next-now))
      select {
      case <-sleep.C:
        logger.Info("awake")
        break forSleep

      case overData := <-over:
        // drain the channel
      drainOver:
        for {
          logger.Debug(fmt.Sprintf("task(unixts=%d) run over", overData.ts))
          getDB().Delete(overData.ts)
          delete(running, overData.ts)
          if overData.crashed {
            logger.Error("task(unixts=", overData.ts, ") runner crashed!")
            // 延迟5s再试
            AddTask(UnixTimeSecond(time.Now().Unix() + 5), overData.tasks)
          }

          select {
          case overData = <-over:
          default:
            break drainOver
          }
        }

        if len(running) < 10 {
          start = UnixTimeSecond(0)
          logger.Debug("len(running) < 10, move start time to ", start)
        }

      case t := <-newTask:
        // drain the channel, find min
        min := t.ts
      drainNewTask:
        for {
          logger.Info(fmt.Sprintf("new task(exeTs=%ds)", t.ts))
          if t.ts < min {
            min = t.ts
          }
          select {
          case t = <-newTask:
          default:
            break drainNewTask
          }
        }

        expect := min
        // 比上一次求取的next更远或者相等(expect >= next)，则不做任何变动
        if expect >= next {
          break
        }
        // start >= expect, 移动start到 expect 之前
        if start >= expect {
          start = expect - 1
          logger.Debug("move start time to ", start)
        }
        // 更新next
        next = expect

        // 停止之前的定时器
        if !sleep.Stop() {
          // drain the channel
          <-sleep.C
        }

        now = UnixTimeSecond(time.Now().Unix())
        // 应该立即执行的，就直接结束sleep，执行任务，以免多一次无效的sleep
        if next < now+1 {
          break forSleep
        }
        // 更新定时器
        sleep.Reset((time.Duration(next - now)) * time.Second)

      case <-ctx.Done():
        return
      }
    }
  }
}
