package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	rpgo "github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/services/config"
	configtypes "github.com/rocket-pool/smartnode/shared/types/config"
	"github.com/t0mk/rocketreport/types"
	"github.com/t0mk/rocketreport/zaplog"
)

var RocketStorageAddress common.Address
var NodeAddress common.Address
var Eth1Url string
var Eth2Url string
var Debug bool
var CachedRplPrice *float64
var ChosenFiat types.Denom
var RpConfig *config.RocketPoolConfig
var BC *services.BeaconClientManager
var EC *services.ExecutionClientManager
var RP *rpgo.RocketPool
var SnConfig *config.SmartnodeConfig
var Network configtypes.Network

var Bot *tgbotapi.BotAPI
var TelegramToken string
var TelegramChatId int64

type EthClientType string

const (
	Eth1 EthClientType = "eth1"
	Eth2 EthClientType = "eth2"
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
)

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("missing env var %s", key))
	}
	return value
}

func getBeaconHeadTimed(ctx context.Context) error {
	errChan := make(chan error)
	go func() {
		_, err := BC.GetBeaconHead()
		errChan <- err
	}()

	select {
	case e := <-errChan:
		return e
	case <-ctx.Done():
		return fmt.Errorf("timeout, make sure consensus client is ready at %s", Eth2Url)
	}
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
		panic(err)
	}
	port := "8545"
	if t == Eth2 {
		port = "5052"
	}
	return fmt.Sprintf("http://%s:%s", ips[0], port)
}

func Setup() {

	log := zaplog.New()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	if (os.Getenv(DebugEnv) != "") && (os.Getenv(DebugEnv) != "0") {
		Debug = true
	}

	Eth1Url = findEthClientUrl(Eth1)
	Eth2Url = findEthClientUrl(Eth2)
	NodeAddress = common.HexToAddress(getEnvOrPanic("NODE_ADDRESS"))
	RocketStorageAddress = common.HexToAddress(
		getEnvOrPanic("ROCKETSTORAGE_ADDRESS"),
	)
	fiatValue := os.Getenv(fiatEnv)
	if fiatValue == "" {
		ChosenFiat = types.USD
	}
	ChosenFiat = types.Denom(fiatValue)

	network := getEnvOrPanic(NetworkEnv)
	switch network {
	case "mainnet":
		Network = configtypes.Network_Mainnet
	case "holesky":
		Network = configtypes.Network_Holesky
	default:
		panic(fmt.Sprintf("Unknown network: %s", network))
	}

	RpConfig = &config.RocketPoolConfig{
		//IsNativeMode: true,
		ConsensusClientMode:     configtypes.Parameter{Value: configtypes.Mode_External},
		ExternalConsensusClient: configtypes.Parameter{Value: configtypes.ConsensusClient_Lighthouse},
		ExternalLighthouse: &config.ExternalLighthouseConfig{
			HttpUrl: configtypes.Parameter{Value: Eth2Url},
		},
		ExecutionClientMode: configtypes.Parameter{Value: configtypes.Mode_External},
		ExternalExecution: &config.ExternalExecutionConfig{
			HttpUrl: configtypes.Parameter{Value: Eth1Url},
		},

		Native: &config.NativeConfig{
			CcHttpUrl: configtypes.Parameter{Value: Eth1Url},
			EcHttpUrl: configtypes.Parameter{Value: Eth2Url},
			ConsensusClient: configtypes.Parameter{
				Value: configtypes.ConsensusClient_Lighthouse,
			},
		},
	}

	SnConfig = config.NewSmartnodeConfig(RpConfig)
	SnConfig.Network = configtypes.Parameter{Value: Network}

	log.Debug("Setting up beacon client")
	BC, err = services.NewBeaconClientManager(RpConfig)
	if err != nil {
		panic(err)
	}

	log.Debug("Getting beacon head")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = getBeaconHeadTimed(ctx)
	/*
		_, err = BC.GetBeaconHead()
	*/
	if err != nil {
		panic(fmt.Sprintf("Beacon client maybe not working: %s", err))
	}
	log.Debug("Setting up execution client")
	EC, err = services.NewExecutionClientManager(RpConfig)
	if err != nil {
		panic(err)
	}
	log.Debug("Getting new RP object")
	RP, err = rpgo.NewRocketPool(EC, RocketStorageAddress)
	if err != nil {
		panic(err)
	}
	TelegramToken = getEnvOrPanic(telegramTokenEnv)
	TelegramChatIdStr := getEnvOrPanic(telegramChatIdEnv)
	TelegramChatId, err = strconv.ParseInt(TelegramChatIdStr, 10, 64)
	if err != nil {
		panic(err)
	}
	log.Debug("Setting up telegram bot")
	Bot, err = tgbotapi.NewBotAPI(TelegramToken)
	if err != nil {
		panic(err)
	}

}
