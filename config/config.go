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
var Plugins PluginConfs

var RocketStorageAddress = map[string]common.Address{
	"mainnet": common.HexToAddress("0x1d8f8f00cfa6758d7bE78336684788Fb0ee0Fa46"),
	"holesky": common.HexToAddress("0x594Fb75D3dc2DFa0150Ad03F99F97817747dd4E1"),
}

type EthClientType string

const (
	Eth1 EthClientType = "eth1"
	Eth2 EthClientType = "eth2"
)

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
	Labl string        `yaml:"labl" json:"labl"`
	Args []interface{} `yaml:"args" json:"args"`
	Hide bool          `yaml:"hide" json:"hide"`
}

type PluginConfs []PluginConf

func (pcs PluginConfs) String() string {
	s := ""
	for i, p := range pcs {
		if p.Name != "" {
			s += fmt.Sprintf("name: %s\n", p.Name)
		}
		if p.Desc != "" {
			s += fmt.Sprintf("desc: %s\n", p.Desc)
		}
		if p.Labl != "" {
			s += fmt.Sprintf("labl: %s\n", p.Labl)
		}
		if p.Args != nil {
			s += "args:\n"
			for _, a := range p.Args {
				s += fmt.Sprintf("  - %v\n", a)
			}
		}
		if p.Hide {
			s += "hide: true\n"
		}
		if i < len(pcs)-1 {
			s += "----\n"
		}
	}
	return s
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

func (cd ConfigData) String() string {
	return fmt.Sprintf("NodeAddress: %s\nEth1Url: %s\nEth2Url: %s\nConsensusClient: %s\nNetwork: %s\nFiat: %s\nTelegramToken: %s\nTelegramChatId: %d\nDebug: %t\n", cd.NodeAddress, cd.Eth1Url, cd.Eth2Url, cd.ConsensusClient, cd.Network, cd.Fiat, cd.TelegramToken, cd.TelegramChatId, cd.Debug)

}

func FileToPlugins(file string) []PluginConf {
	pluginsWrap := struct {
		Plugins []PluginConf `yaml:"plugins"`
	}{}
	loader := aconfig.LoaderFor(&pluginsWrap, aconfig.Config{
		SkipFlags: true,
		SkipEnv:   true,
		Files:     []string{file},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
			".yml":  aconfigyaml.New(),
		},
	})
	err := loader.Load()
	if err != nil {
		panic(err)
	}
	//fmt.Println(pluginsWrap)
	return pluginsWrap.Plugins
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
