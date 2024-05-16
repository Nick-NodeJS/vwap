package calculator

import (
	"log"
	"sync"
	"time"

	"vwap/pkg/trade"
)

type VWAPCalculator struct {
	Data      map[string][]trade.Trade
	mu        sync.Mutex
	TimeFrame time.Duration
}

func NewVWAPCalculator(timeFrame time.Duration) *VWAPCalculator {
	return &VWAPCalculator{
		Data:      make(map[string][]trade.Trade),
		TimeFrame: timeFrame,
	}
}

func (v *VWAPCalculator) AddTrade(trade trade.Trade) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if trade.Volume <= 0 || trade.Price <= 0 {
		log.Printf("Invalid trade data ignored: %+v", trade)
		return
	}

	v.Data[trade.Pair] = append(v.Data[trade.Pair], trade)
	v.cleanup(trade.Pair)
}

func (v *VWAPCalculator) cleanup(pair string) {
	threshold := time.Now().Add(-v.TimeFrame)
	var i int
	for i = 0; i < len(v.Data[pair]) && v.Data[pair][i].Time.Before(threshold); i++ {
	}
	if i > 0 {
		v.Data[pair] = v.Data[pair][i:]
	}
}

func (v *VWAPCalculator) CalculateVWAP(pair string) float64 {
	v.mu.Lock()
	defer v.mu.Unlock()

	var totalVolume, totalPriceVolume float64
	for _, trade := range v.Data[pair] {
		totalVolume += trade.Volume
		totalPriceVolume += trade.Price * trade.Volume
	}

	if totalVolume == 0 {
		return 0
	}
	return totalPriceVolume / totalVolume
}
