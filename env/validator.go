// validator.go
package env

import (
	"log"
	"os"
	"strings"
)

// GetBitfinexAPIURL returns the Bitfinex API URL from the environment variables.
// If the variable is not set or fails validation, it returns a default value.
func GetBitfinexAPIURL() string {
	bitfinexAPIURL := os.Getenv("BITFINEX_API_URL")
	if bitfinexAPIURL == "" {
		log.Println("BITFINEX_API_URL is not set. Using default value.")
		return "wss://api-pub.bitfinex.com/ws/2" // Default value
	}
	return bitfinexAPIURL
}

// GetPairs returns the list of pairs from the environment variables.
// If the variable is not set or fails validation, it returns a default value.
func GetPairs() []string {
	pairsString := os.Getenv("PAIRS")
	if pairsString == "" {
		log.Println("PAIRS is not set. Using default value.")
		return []string{"BTCUSD", "ETHUSD"} // Default value
	}
	return strings.Split(pairsString, ",")
}
