package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
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
	PluginListEnv           = "PLUGIN_LIST"
)

var CachedRplPrice *float64

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

type PluginConf map[string]struct {
	Key  string
	Desc string
	Args []string
}

type ConfigData struct {
	NodeAddress          string `env:"NODE_ADDRESS" yaml:"node_address" json:"node_address"`
	Eth1Url              string `env:"ETH1_URL" yaml:"eth1_url" json:"eth1_url"`
	Eth2Url              string `env:"ETH2_URL" yaml:"eth2_url" json:"eth2_url"`
	ConsensusClient      string `env:"CONSENSUS_CLIENT" yaml:"consensus_client" json:"consensus_client"`
	Network              string `env:"NETWORK" yaml:"network" json:"network"`
	RocketStorageAddress string `env:"ROCKETSTORAGE_ADDRESS" yaml:"rocketstorage_address" json:"rocketstorage_address"`
	Fiat                 string `default:"USD" env:"FIAT" yaml:"fiat" json:"fiat"`
	TelegramToken        string `env:"TELEGRAM_TOKEN" yaml:"telegram_token" json:"telegram_token"`
	TelegramChatId       int64  `env:"TELEGRAM_CHAT_ID" yaml:"telegram_chat_id" json:"telegram_chat_id"`
	Debug                bool   `env:"DEBUG" yaml:"debug" json:"debug"`
	PluginConf           []PluginConf `env:"PLUGIN_LIST" yaml:"plugin_list" json:"plugin_list"`
}

var c ConfigData

var EC = sync.OnceValue(initEC)
var BC = sync.OnceValue(initBC)
var RpConfig = sync.OnceValue(initRpConfig)
var NodeAddress = sync.OnceValue(initNodeAddress)
var RocketStorageAddress = sync.OnceValue(initRocketStorageAddress)
var Network = sync.OnceValue(initNetwork)
var RP = sync.OnceValue(initRP)
var PluginList = sync.OnceValue(initPluginList)
var ChosenFiat = sync.OnceValue(initChosenFiat)
var TelegramChatID = sync.OnceValue(initTelegramChatId)
var TelegramToken = sync.OnceValue(initTelegramToken)
var TelegramBot = sync.OnceValue(initTelegramBot)
var Debug bool

func Setup() {
	loader := aconfig.LoaderFor(&c, aconfig.Config{
		Files: []string{"config.yaml", "config.yml", "config.json"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
			".yml":  aconfigyaml.New(),
		},
	})
	err := loader.Load()
	if err != nil {
		panic(err)
	}
	/*
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	*/
}
