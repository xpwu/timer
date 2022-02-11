package fixed

import (
  "encoding/json"
  "testing"
)

func TestCronTime(t *testing.T) {
  l := &CronTime{}
  d, err := json.Marshal(l)
  if err != nil {
    t.Error(err)
    return
  }

  t.Log(string(d))
}
