package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/ethereum/go-ethereum/common"
)

const (
	NodeAddressEnv     = "NODE_ADDRESS"
	Eth1UrlEnv         = "ETH1_URL"
	Eth2UrlEnv         = "ETH2_URL"
	DebugEnv           = "DEBUG"
	fiatEnv            = "FIAT"
	telegramTokenEnv   = "TELEGRAM_TOKEN"
	telegramChatIdEnv  = "TELEGRAM_CHAT_ID"
	NetworkEnv         = "NETWORK"
	ConsensusClientEnv = "CONSENSUS_CLIENT"
)

func RewardTreesPath() (string, error) {
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("can't get os.UserHomeDir: %v", err)
	}
	return filepath.Join(hd, ".rocketreport/reward-trees", string(Network())), nil
}

var debugFlag bool

func SetDebugFlag() {
	debugFlag = true
}

func Debug() bool {
	return debugFlag || Config.Debug
}

var Config ConfigData
var EC = sync.OnceValue(initEC)
var BC = sync.OnceValue(initBC)
var RpConfig = sync.OnceValue(initRpConfig)
var NodeAddress = sync.OnceValue(initNodeAddress)
var Network = sync.OnceValue(initNetwork)
var RP = sync.OnceValue(initRP)
var ChosenFiat = sync.OnceValue(initChosenFiat)
var TelegramChatID = sync.OnceValue(initTelegramChatId)
var TelegramToken = sync.OnceValue(initTelegramToken)
var TelegramBot = sync.OnceValue(initTelegramBot)
var TelegramMessageSchedule = sync.OnceValue(initTelegramMessageSchedule)
var TelegramHeaderTemplate = sync.OnceValue(initTelegramHeaderTemplate)

var RocketStorageAddress = map[string]common.Address{
	"mainnet": common.HexToAddress("0x1d8f8f00cfa6758d7bE78336684788Fb0ee0Fa46"),
	"holesky": common.HexToAddress("0x594Fb75D3dc2DFa0150Ad03F99F97817747dd4E1"),
}

type EthClientType string

const (
	Eth1 EthClientType = "eth1"
	Eth2 EthClientType = "eth2"
)

type ConfigData struct {
	NodeAddress             string `env:"NODE_ADDRESS" yaml:"node_address" json:"node_address"`
	Eth1Url                 string `env:"ETH1_URL" yaml:"eth1_url" json:"eth1_url"`
	Eth2Url                 string `env:"ETH2_URL" yaml:"eth2_url" json:"eth2_url"`
	ConsensusClient         string `env:"CONSENSUS_CLIENT" yaml:"consensus_client" json:"consensus_client"`
	Network                 string `env:"NETWORK" yaml:"network" json:"network"`
	Fiat                    string `default:"USD" env:"FIAT" yaml:"fiat" json:"fiat"`
	TelegramToken           string `env:"TELEGRAM_TOKEN" yaml:"telegram_token" json:"telegram_token"`
	TelegramChatId          int64  `env:"TELEGRAM_CHAT_ID" yaml:"telegram_chat_id" json:"telegram_chat_id"`
	TelegramMessageSchedule string `env:"TELEGRAM_MESSAGE_SCHEDULE" yaml:"telegram_message_schedule" json:"telegram_message_schedule"`
	TelegramHeaderTemplate  string `env:"TELEGRAM_HEADER_TEMPLATE" yaml:"telegram_header_template" json:"telegram_header_template"`
	Debug                   bool   `env:"DEBUG" yaml:"debug" json:"debug"`
}

func (cd ConfigData) String() string {
	return fmt.Sprintf("NodeAddress: %s\nEth1Url: %s\nEth2Url: %s\nConsensusClient: %s\nNetwork: %s\nFiat: %s\nTelegramToken: %s\nTelegramChatId: %d\nDebug: %t\n", cd.NodeAddress, cd.Eth1Url, cd.Eth2Url, cd.ConsensusClient, cd.Network, cd.Fiat, cd.TelegramToken, cd.TelegramChatId, cd.Debug)

}

func LoadConfigFromFile(file string) {
	loader := aconfig.LoaderFor(&Config, aconfig.Config{
		SkipFlags: true,
		//SkipEnv:   true,
		Files: []string{file},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
			".yml":  aconfigyaml.New(),
		},
	})

	err := loader.Load()
	if err != nil {
		panic(err)
	}
}
