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

	calculator := calculator.NewVWAPCalculator(2 * time.Minute)
	channelPairMapping := make(map[int]string)

	// Dynamically manage subscriptions
	vwapwebsocket.ManageSubscriptions(conn, pairs)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage error:", err)
			continue
		}

		vwapwebsocket.HandleMessage(message, channelPairMapping, calculator)

		for _, pair := range pairs {
			vwap := calculator.CalculateVWAP(pair)
			log.Printf("VWAP for %s: %.2f", pair, vwap)
		}
	}
}
