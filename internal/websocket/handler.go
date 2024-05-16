package websocket

import (
	"encoding/json"
	"log"
	"math"
	"time"

	"vwap/internal/calculator"
	"vwap/pkg/trade"
)

func HandleMessage(message []byte, channelPairMapping map[int]string, calculator *calculator.VWAPCalculator) {

	var msgMap map[string]interface{}
	if err := json.Unmarshal(message, &msgMap); err == nil {
		if event, ok := msgMap["event"].(string); ok {
			if serverId, ok := msgMap["serverId"].(string); ok && event == "info" {
				log.Println("Server connected: serverId ", serverId)
				return
			}
			if event == "error" {
				log.Println("Error: ", msgMap)
				return
			}
			if event == "subscribed" {
				handleSubscriptionMessage(msgMap, channelPairMapping)
				return
			}
		}
	}

	var genericMsg []interface{}
	if err := json.Unmarshal(message, &genericMsg); err != nil {
		log.Println("Error unmarshalling message:", err)
		return
	}

	if len(genericMsg) > 2 {
		handleTradeMessage(genericMsg, channelPairMapping, calculator)
	}
}

func handleSubscriptionMessage(msgMap map[string]interface{}, channelPairMapping map[int]string) {
	if channel, ok := msgMap["chanId"].(float64); ok {
		if pair, ok := msgMap["pair"].(string); ok {
			channelPairMapping[int(channel)] = pair
			log.Printf("Subscribed to %s with channel ID %d", pair, int(channel))
		}
	}
}

func handleTradeMessage(genericMsg []interface{}, channelPairMapping map[int]string, calculator *calculator.VWAPCalculator) {
	channelID, ok := genericMsg[0].(float64)
	if !ok {
		return
	}

	pair, exists := channelPairMapping[int(channelID)]
	if !exists {
		log.Println("Received message for unknown channel ID:", channelID)
		return
	}

	messageType, ok := genericMsg[1].(string)
	if !ok || (messageType != "te" && messageType != "tu") {
		return
	}

	tradeDetails, ok := genericMsg[2].([]interface{})
	if !ok || len(tradeDetails) < 4 {
		log.Println("Invalid trade details format")
		return
	}

	processTradeDetails(tradeDetails, pair, calculator)
}

func processTradeDetails(tradeDetails []interface{}, pair string, calculator *calculator.VWAPCalculator) {
	timestamp := int64(tradeDetails[1].(float64))
	amount := math.Abs(tradeDetails[2].(float64))
	price := tradeDetails[3].(float64)

	calculator.AddTrade(trade.Trade{
		Pair:   pair,
		Price:  price,
		Volume: amount,
		Time:   time.Unix(timestamp/1000, 0),
	})

	log.Printf("Trade processed for %s: Price: %f, Volume: %f", pair, price, amount)
}
