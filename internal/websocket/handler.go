package websocket

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"time"

	"vwap/internal/calculator"
	"vwap/pkg/trade"
)

func HandleMessage(message []byte, channelPairMapping map[int]string, calculator *calculator.VWAPCalculator) {

	var msgMap map[string]interface{}
	if err := json.Unmarshal(message, &msgMap); err == nil {
		if event, ok := msgMap[MESSAGE_EVENT].(string); ok {
			switch event {
			default:
				log.Printf("Unknow event: %s", event)
			case MESSAGE_EVENT_INFO:
				if serverId, ok := msgMap[MESSAGE_SERVER_ID].(string); ok {
					log.Printf("Server connected. serverId: %s", serverId)
					return
				}
			case MESSAGE_EVENT_ERROR:
				log.Println("Error: ", msgMap)
			case MESSAGE_EVENT_SUBSCRIBED:
				handleSubscriptionMessage(msgMap, channelPairMapping)
			}
			return
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
	if channel, ok := msgMap[MESSAGE_CHANNEL_ID].(float64); ok {
		if pair, ok := msgMap[MESSAGE_PAIR].(string); ok {
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
		log.Printf("Received message for unknown channel ID: %f", channelID)
		return
	}

	tradeDetails, err := getTradeDetails(genericMsg)
	if err != nil {
		log.Println("Error to get trade details: ", err)
		return
	}

	processTradeDetails(tradeDetails, pair, calculator)
}

func getTradeDetails(genericMsg []interface{}) (tradeDetails []interface{}, err error) {
	messageType, ok := genericMsg[1].(string)
	if !ok || (messageType != MESSAGE_TRADE_TYPE_TE && messageType != MESSAGE_TRADE_TYPE_TU) {
		return nil, errors.New("wrong trade message type " + messageType)
	}

	tradeDetails = genericMsg[2].([]interface{})
	if len(tradeDetails) < 4 {
		return nil, errors.New("invalid trade details format")
	}
	return tradeDetails, nil
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
