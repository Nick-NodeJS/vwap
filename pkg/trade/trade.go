package trade

import "time"

type Trade struct {
	Pair   string
	Price  float64
	Volume float64
	Time   time.Time
}
