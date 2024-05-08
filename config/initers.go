package config

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	rpgo "github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/services/config"
	"github.com/t0mk/rocketreport/zaplog"

	configtypes "github.com/rocket-pool/smartnode/shared/types/config"
)

func initNodeAddress() common.Address {
	if Config.NodeAddress == "" {
		panic("NodeAddress not set in config")
	}
	// check that it's valid eth address

	if !common.IsHexAddress(Config.NodeAddress) {
		s := fmt.Sprintf("wrong value for node_address in config file. \"%s\" is not a valid Ethereum address", Config.NodeAddress)
		panic(s)
	}

	return common.HexToAddress(Config.NodeAddress)
}

func Eth1Url() string {
	if Config.Eth1Url == "" {
		panic("Eth1Url not set in config")
	}
	return Config.Eth1Url
}

func Eth2Url() string {
	if Config.Eth2Url == "" {
		panic("Eth2Url not set in config")
	}
	return Config.Eth2Url
}

func ConsensusClient() configtypes.ConsensusClient {
	if Config.ConsensusClient == "" {
		panic("ConsensusClient not set in config")
	}
	return configtypes.ConsensusClient(Config.ConsensusClient)
}

func initNetwork() configtypes.Network {
	if Config.Network == "" {
		panic("Network not set in config")
	}
	switch Config.Network {
	case "mainnet":
		return configtypes.Network_Mainnet
	case "holesky":
		return configtypes.Network_Holesky
	}
	panic(fmt.Sprintf("Unknown network %s", Config.Network))
}

func initChosenFiat() string {
	if Config.Fiat == "" {
		return "USDT"
	}
	if len(Config.Fiat) != 3 {
		panic("Fiat currency must be 3 letters or \"USDT\"")
	}
	for _, c := range Config.Fiat {
		if c < 'A' || c > 'Z' {
			panic("Fiat currency must be all capital letters")
		}
	}
	return Config.Fiat
}

func initTelegramToken() string {
	if Config.TelegramToken == "" {
		panic("TelegramToken not set in config")
	}
	return Config.TelegramToken
}

func initTelegramChatId() int64 {
	if Config.TelegramChatId == 0 {
		panic("TelegramChatId not set in config")
	}
	return Config.TelegramChatId
}

func initRpConfig() *config.RocketPoolConfig {
	eth1Url := Eth1Url()
	eth2Url := Eth2Url()
	consensusClient := ConsensusClient()
	ret := &config.RocketPoolConfig{
		ConsensusClientMode: configtypes.Parameter{Value: configtypes.Mode_External},
		ExecutionClientMode: configtypes.Parameter{Value: configtypes.Mode_External},
	}
	if consensusClient != "" {
		ret.ExternalConsensusClient = configtypes.Parameter{Value: consensusClient}
		if eth2Url != "" {
			switch consensusClient {
			case configtypes.ConsensusClient_Lighthouse:
				ret.ExternalLighthouse = &config.ExternalLighthouseConfig{
					HttpUrl: configtypes.Parameter{Value: eth2Url},
				}
			case configtypes.ConsensusClient_Prysm:
				ret.ExternalPrysm = &config.ExternalPrysmConfig{
					HttpUrl: configtypes.Parameter{Value: eth2Url},
				}
			case configtypes.ConsensusClient_Teku:
				ret.ExternalTeku = &config.ExternalTekuConfig{
					HttpUrl: configtypes.Parameter{Value: eth2Url},
				}
			}
		}
	}
	if eth1Url != "" {
		ret.ExternalExecution = &config.ExternalExecutionConfig{
			HttpUrl: configtypes.Parameter{Value: eth1Url},
		}
	}
	nat := &config.NativeConfig{}
	if eth1Url != "" {
		nat.CcHttpUrl = configtypes.Parameter{Value: eth1Url}
	}
	if eth2Url != "" {
		nat.EcHttpUrl = configtypes.Parameter{Value: eth2Url}
	}
	if consensusClient != "" {
		nat.ConsensusClient = configtypes.Parameter{Value: consensusClient}
	}
	ret.Native = nat
	ret.Smartnode = config.NewSmartnodeConfig(ret)
	ret.Smartnode.Network = configtypes.Parameter{Value: initNetwork()}
	ret.Smartnode.DataPath = configtypes.Parameter{Value: "./"}

	//fmt.Printf("#%v\n", ret.Smartnode)
	return ret
}

func initEC() *services.ExecutionClientManager {
	ec, err := services.NewExecutionClientManager(RpConfig())
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = ec.BlockNumber(ctx)
	if err != nil {
		panic(fmt.Sprintf("Execution client probably not available at %s: %s", Eth1Url(), err))
	}
	return ec
}

func initBC() *services.BeaconClientManager {
	log := zaplog.New()
	bc, err := services.NewBeaconClientManager(RpConfig())
	if err != nil {
		panic(err)
	}

	log.Debug("Getting beacon head")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	errChan := make(chan error)
	go func() {
		_, err := bc.GetBeaconHead()
		errChan <- err
	}()
	eth2Url := Eth2Url()

	select {
	case err := <-errChan:
		if err != nil {
			panic(fmt.Sprintf("Beacon client probably not available at %s:\n%s", eth2Url, err))
		}
	case <-ctx.Done():
		panic(fmt.Sprintf("Timeout pinging Beacon client, make sure it's available at %s", eth2Url))
	}
	return bc
}

func initRP() *rpgo.RocketPool {
	rp, err := rpgo.NewRocketPool(EC(), RocketStorageAddress[string(Network())])
	if err != nil {
		panic(err)
	}
	return rp
}

func initTelegramBot() *tgbotapi.BotAPI {
	token := TelegramToken()
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	return bot
}
