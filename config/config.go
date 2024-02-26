package config

import (
	"fmt"
	"os"
	"strconv"

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

var Bot *tgbotapi.BotAPI
var TelegramToken string
var TelegramChatId int64

const (
	RocketStorageAddressEnv = "ROCKETSTORAGE_ADDRESS"
	NodeAddressEnv          = "NODE_ADDRESS"
	Eth1UrlEnv              = "ETH1_URL"
	Eth2UrlEnv              = "ETH2_URL"
	DebugEnv                = "DEBUG"
	fiatEnv                 = "FIAT"
	telegramTokenEnv        = "TELEGRAM_TOKEN"
	telegramChatIdEnv       = "TELEGRAM_CHAT_ID"
)

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("missing env var %s", key))
	}
	return value
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

	Eth1Url = getEnvOrPanic("ETH1_URL")
	NodeAddress = common.HexToAddress(getEnvOrPanic("NODE_ADDRESS"))
	RocketStorageAddress = common.HexToAddress(
		getEnvOrPanic("ROCKETSTORAGE_ADDRESS"),
	)
	fiatValue := os.Getenv(fiatEnv)
	if fiatValue == "" {
		ChosenFiat = types.USD
	}
	ChosenFiat = types.Denom(fiatValue)

	Eth2Url = getEnvOrPanic("ETH2_URL")

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
	log.Debug("Setting up beacon client")
	BC, err = services.NewBeaconClientManager(RpConfig)
	if err != nil {
		panic(err)
	}
	log.Debug("Getting beacon head")
	_, err = BC.GetBeaconHead()
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
