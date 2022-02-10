package scheduler

import "context"

type Runner interface {
  Run(ctx context.Context, schedulerTime UnixTimeSecond, tasks []Task)
}
