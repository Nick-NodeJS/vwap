package calculator

import (
	"log"
	"math"
	"sync"
	"time"

	"vwap/pkg/trade"
)

type VWAPCalculator struct {
	Data           map[string][]trade.Trade
	mu             sync.Mutex
	TimeFrame      time.Duration
	PriceSum       map[string]float64
	PriceSquareSum map[string]float64
}

func NewVWAPCalculator(timeFrame time.Duration) *VWAPCalculator {
	return &VWAPCalculator{
		Data:           make(map[string][]trade.Trade),
		TimeFrame:      timeFrame,
		PriceSum:       make(map[string]float64),
		PriceSquareSum: make(map[string]float64),
	}
}

func (v *VWAPCalculator) AddTrade(trade trade.Trade) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if trade.Volume <= 0 || trade.Price <= 0 {
		log.Printf("Invalid trade data ignored: %+v", trade)
		return
	}

	deviation := v.CalculateStandardDeviation(trade.Pair, true)

	if deviation != 0 && trade.Price > 3*deviation {
		log.Printf("Wrong trade price deviation, standart deviation %f, trade price %f", deviation, trade.Price)
		return
	}

	v.PriceSum[trade.Pair] += trade.Price
	v.PriceSquareSum[trade.Pair] += trade.Price * trade.Price
	v.Data[trade.Pair] = append(v.Data[trade.Pair], trade)
	v.cleanup(trade.Pair)
}

func (v *VWAPCalculator) cleanup(pair string) {
	threshold := time.Now().Add(-v.TimeFrame)
	var i int
	for i = 0; i < len(v.Data[pair]) && v.Data[pair][i].Time.Before(threshold); i++ {
		v.PriceSum[pair] -= v.Data[pair][i].Price
		v.PriceSquareSum[pair] -= v.Data[pair][i].Price * v.Data[pair][i].Price
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

func (v *VWAPCalculator) CalculateMean(pair string) float64 {
	v.mu.Lock()
	defer v.mu.Unlock()

	count := len(v.Data[pair])
	if count == 0 {
		return 0
	}
	return v.PriceSum[pair] / float64(count)
}

func (v *VWAPCalculator) CalculateStandardDeviation(pair string, internal bool) float64 {
	if !internal {
		v.mu.Lock()
		defer v.mu.Unlock()
	}

	count := len(v.Data[pair])
	if count == 0 {
		return 0
	}
	mean := v.PriceSum[pair] / float64(count)
	variance := (v.PriceSquareSum[pair] / float64(count)) - (mean * mean)
	return math.Sqrt(variance)
}
