package main

import (
	"testing"
	"time"

	"vwap/internal/calculator"
	"vwap/pkg/trade"
)

func TestVWAPCalculator_CalculateVWAP(t *testing.T) {
	calculator := calculator.NewVWAPCalculator(7 * time.Minute)

	// Add some test trades
	timestamp1 := time.Now().Add(-5 * time.Minute)
	timestamp2 := time.Now().Add(-4 * time.Minute)
	timestamp3 := time.Now().Add(-3 * time.Minute)
	timestamp4 := time.Now().Add(-2 * time.Minute)
	calculator.AddTrade(trade.Trade{Pair: "BTCUSD", Price: 100, Volume: 10, Time: timestamp1})
	calculator.AddTrade(trade.Trade{Pair: "BTCUSD", Price: 110, Volume: 15, Time: timestamp2})

	calculator.AddTrade(trade.Trade{Pair: "BTCUSD", Price: 910, Volume: 5, Time: timestamp3})

	calculator.AddTrade(trade.Trade{Pair: "ETHUSD", Price: 200, Volume: 5, Time: timestamp4})

	// Test VWAP calculation
	expectedBTCUSDVWAP := 106.0 // ((100 * 10) + (110 * 15)) / (10 + 15)
	actualBTCUSDVWAP := calculator.CalculateVWAP("BTCUSD")
	if actualBTCUSDVWAP != expectedBTCUSDVWAP {
		t.Errorf("VWAP calculation for BTCUSD incorrect. Expected: %.2f, Got: %.2f", expectedBTCUSDVWAP, actualBTCUSDVWAP)
	}

	expectedETHUSDVWAP := 200.0
	actualETHUSDVWAP := calculator.CalculateVWAP("ETHUSD")
	if actualETHUSDVWAP != expectedETHUSDVWAP {
		t.Errorf("VWAP calculation for ETHUSD incorrect. Expected: %.2f, Got: %.2f", expectedETHUSDVWAP, actualETHUSDVWAP)
	}
}
