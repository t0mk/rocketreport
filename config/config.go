package config

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	rpgo "github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/services/config"
	configtypes "github.com/rocket-pool/smartnode/shared/types/config"
	"github.com/t0mk/rocketreport/prices"
)

var rocketStorageAddress common.Address
var nodeAddress common.Address
var eth1Url string
var eth2Url string
var debug bool
var cachedRplPrice *float64
var chosenFiat prices.Fiat
var rpConfig *config.RocketPoolConfig
var bc *services.BeaconClientManager
var ec *services.ExecutionClientManager
var rp *rpgo.RocketPool

var bot *tgbotapi.BotAPI
var telegramToken string

const (
	rocketStorageAddressEnv = "ROCKETSTORAGE_ADDRESS"
	nodeAddressEnv          = "NODE_ADDRESS"
	eth1UrlEnv              = "ETH1_URL"
	eth2UrlEnv              = "ETH2_URL"
	debugEnv                = "DEBUG"
	fiatEnv                 = "FIAT"
	telegramTokenEnv        = "TELEGRAM_TOKEN"
)

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("missing env var %s", key))
	}
	return value
}

func doConfig() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	if (os.Getenv(debugEnv) != "") && (os.Getenv(debugEnv) != "0") {
		debug = true
	}

	eth1Url = getEnvOrPanic("ETH1_URL")
	nodeAddress = common.HexToAddress(getEnvOrPanic("NODE_ADDRESS"))
	rocketStorageAddress = common.HexToAddress(
		getEnvOrPanic("ROCKETSTORAGE_ADDRESS"),
	)
	fiatValue := os.Getenv(fiatEnv)
	if fiatValue == "" {
		chosenFiat = USD
	}
	chosenFiat = Fiat(fiatValue)

	eth2Url = getEnvOrPanic("ETH2_URL")

	rpConfig = &config.RocketPoolConfig{
		//IsNativeMode: true,
		ConsensusClientMode:     configtypes.Parameter{Value: configtypes.Mode_External},
		ExternalConsensusClient: configtypes.Parameter{Value: configtypes.ConsensusClient_Lighthouse},
		ExternalLighthouse: &config.ExternalLighthouseConfig{
			HttpUrl: configtypes.Parameter{Value: eth2Url},
		},
		ExecutionClientMode: configtypes.Parameter{Value: configtypes.Mode_External},
		ExternalExecution: &config.ExternalExecutionConfig{
			HttpUrl: configtypes.Parameter{Value: eth1Url},
		},

		Native: &config.NativeConfig{
			CcHttpUrl: configtypes.Parameter{Value: eth1Url},
			EcHttpUrl: configtypes.Parameter{Value: eth2Url},
			ConsensusClient: configtypes.Parameter{
				Value: configtypes.ConsensusClient_Lighthouse,
			},
		},
	}

	bc, err = services.NewBeaconClientManager(rpConfig)
	if err != nil {
		panic(err)
	}
	ec, err = services.NewExecutionClientManager(rpConfig)
	if err != nil {
		panic(err)
	}
	rp, err = rpgo.NewRocketPool(ec, rocketStorageAddress)
	if err != nil {
		panic(err)
	}
}
