package main

import (
	"log"
	"time"

	"vwap/env"
	"vwap/internal/calculator"
	vwapwebsocket "vwap/internal/websocket"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bitfinexAPIURL := env.GetBitfinexAPIURL()
	pairs := env.GetPairs()

	conn, _, err := websocket.DefaultDialer.Dial(bitfinexAPIURL, nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	// Set VWAP calculator
	calculator := calculator.NewVWAPCalculator(2 * time.Minute)
	channelPairMapping := make(map[int]string)

	// Dynamically manage subscriptions
	go func() {
		vwapwebsocket.ManageSubscriptions(conn, pairs)
	}()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("ReadMessage error:", err)
				continue
			}

			vwapwebsocket.HandleMessage(message, channelPairMapping, calculator)
		}
	}()
	// Update VWAP on pairs
	// go func() {
	for {
		time.Sleep(1 * time.Second)
		for _, pair := range pairs {
			vwap := calculator.CalculateVWAP(pair)
			mean := calculator.CalculateMean(pair)
			stddev := calculator.CalculateStandardDeviation(pair, false)
			log.Printf("VWAP for %s: %.2f, Mean: %.2f, Standard Deviation: %.2f", pair, vwap, mean, stddev)
		}
	}
	// }()
}
