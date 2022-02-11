package flag

type Flag = byte

const (
  Invalid Flag = iota
  Fixed
  Ticker
  Delay
)
