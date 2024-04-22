package config

import (
	"fmt"
	"os"
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

type EthClientType string

const (
	Eth1 EthClientType = "eth1"
	Eth2 EthClientType = "eth2"
)

var RocketStorageAddress = map[string]common.Address{
	"mainnet": common.HexToAddress("0x1d8f8f00cfa6758d7bE78336684788Fb0ee0Fa46"),
	"holesky": common.HexToAddress("0x594Fb75D3dc2DFa0150Ad03F99F97817747dd4E1"),
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

type PluginConf struct {
	Name string        `yaml:"name" json:"name"`
	Desc string        `yaml:"desc" json:"desc"`
	Id   string        `yaml:"id" json:"id"`
	Args []interface{} `yaml:"args" json:"args"`
	Opts []string      `yaml:"opts" json:"opts"`
	Mute bool          `yaml:"mute" json:"mute"`
}

type PluginConfs struct {
	Plugins []PluginConf `yaml:"plugins" json:"plugins"`
}

func PluginsString(pcs []PluginConf) string {
	s := "plugins:\n"
	for _, p := range pcs {
		s += fmt.Sprintf("  - name: %s\n", p.Name)
		if p.Args != nil {
			s += "    args:\n"
			for _, a := range p.Args {
				s += fmt.Sprintf("      - %v\n", a)
			}
		}
	}
	return s
}

type ConfigData struct {
	NodeAddress     string `env:"NODE_ADDRESS" yaml:"node_address" json:"node_address"`
	Eth1Url         string `env:"ETH1_URL" yaml:"eth1_url" json:"eth1_url"`
	Eth2Url         string `env:"ETH2_URL" yaml:"eth2_url" json:"eth2_url"`
	ConsensusClient string `env:"CONSENSUS_CLIENT" yaml:"consensus_client" json:"consensus_client"`
	Network         string `env:"NETWORK" yaml:"network" json:"network"`
	Fiat            string `default:"USD" env:"FIAT" yaml:"fiat" json:"fiat"`
	TelegramToken   string `env:"TELEGRAM_TOKEN" yaml:"telegram_token" json:"telegram_token"`
	TelegramChatId  int64  `env:"TELEGRAM_CHAT_ID" yaml:"telegram_chat_id" json:"telegram_chat_id"`
	Debug           bool   `env:"DEBUG" yaml:"debug" json:"debug"`
}

var c ConfigData

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
var Debug bool
var Plugins []PluginConf
var pluginsFileContent PluginConfs

func Setup(configFile, pluginsFile string) {
	cfgFiles := []string{configFile}
	if configFile == "" {
		cfgFiles = []string{}
	}

	loader := aconfig.LoaderFor(&c, aconfig.Config{
		SkipFlags: true,
		Files:     cfgFiles,
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
			".yml":  aconfigyaml.New(),
		},
	})
	err := loader.Load()
	if err != nil {
		panic(err)
	}
	if c.Debug {
		Debug = true
	}
	if Debug {
		fmt.Println("Loaded config")
		fmt.Println(c)
	}
	if pluginsFile == "" {
		pluginsFile = "defaultplugins.yaml"
	}
	loader = aconfig.LoaderFor(&pluginsFileContent, aconfig.Config{
		SkipFlags: true,
		SkipEnv:   true,
		Files:     []string{pluginsFile},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
			".yml":  aconfigyaml.New(),
		},
	})
	err = loader.Load()
	if err != nil {
		panic(err)
	}
	Plugins = pluginsFileContent.Plugins
	if Debug {
		fmt.Println("Loaded plugins")
		fmt.Println(Plugins)
	}
}
