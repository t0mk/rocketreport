package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/t0mk/rocketreport/exchanges"
)

const (
	RocketStorageAddressEnv = "ROCKETSTORAGE_ADDRESS"
	NodeAddressEnv          = "NODE_ADDRESS"
	Eth1UrlEnv              = "ETH1_URL"
	Eth2UrlEnv              = "ETH2_URL"
	DebugEnv                = "DEBUG"
	fiatEnv                 = "FIAT"
	telegramTokenEnv        = "TELEGRAM_TOKEN"
	telegramChatIdEnv       = "TELEGRAM_CHAT_ID"
	NetworkEnv              = "NETWORK"
	ConsensusClientEnv      = "CONSENSUS_CLIENT"
)

var Debug bool
var CachedRplPrice *float64

type Telegram struct {
	Token  string
	ChatId int64
}

type EthClientType string

const (
	Eth1 EthClientType = "eth1"
	Eth2 EthClientType = "eth2"
)

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("missing env var %s", key))
	}
	return value
}

func findEthClientUrl(t EthClientType) string {
	urlVar := Eth1UrlEnv
	if t == Eth2 {
		urlVar = Eth2UrlEnv
	}
	val := os.Getenv(urlVar)
	if val != "" {
		return val
	}
	ips, err := FindContainerIPs(string(t))
	if err != nil {
		panic(errMissingEthClientURL(t))
	}
	port := "8545"
	if t == Eth2 {
		port = "5052"
	}
	return fmt.Sprintf("http://%s:%s", ips[0], port)
}

var EC = sync.OnceValue(initEC)
var BC = sync.OnceValue(initBC)
var RpConfig = sync.OnceValue(initRpConfig)
var NodeAddress = sync.OnceValue(initNodeAddress)
var RocketStorageAddress = sync.OnceValue(initRocketStorageAddress)
var Network = sync.OnceValue(initNetworkValue)
var RP = sync.OnceValue(initRP)
var ChosenFiat = sync.OnceValue(initChosenFiat)
var TelegramChatID = sync.OnceValue(initTelegramChatID)
var TelegramBot = sync.OnceValue(initTelegramBot)



var XchMap = map[string]func(string) (float64, error){
	USD: exchanges.BitfinexPri,
	EUR:  exchanges.BitfinexPri,
	GBP:  exchanges.BitfinexPri,
	JPY:  exchanges.BitfinexPri,
	AUD:  exchanges.KrakenPri,
	CHF:  exchanges.KrakenPri,
	CZK:  exchanges.CoinmatePri,
}

func Setup() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	if (os.Getenv(DebugEnv) != "") && (os.Getenv(DebugEnv) != "0") {
		Debug = true
	}
}
