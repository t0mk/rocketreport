package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	rpgo "github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/services/config"
	"github.com/t0mk/rocketreport/zaplog"

	configtypes "github.com/rocket-pool/smartnode/shared/types/config"
)

func initRocketStorageAddress() common.Address {
	rocketStorageAddress := os.Getenv(RocketStorageAddressEnv)
	if rocketStorageAddress == "" {
		panic("Couldn't find rocket storage address. You can set ROCKETSTORAGE_ADDRESS envvar")
	}
	return common.HexToAddress(rocketStorageAddress)
}

func initNetworkValue() configtypes.Network {
	network := os.Getenv(NetworkEnv)
	switch network {
	case "mainnet":
		return configtypes.Network_Mainnet
	case "holesky":
		return configtypes.Network_Holesky
	}
	panic(fmt.Sprintf("Unknown network %s", network))
}

func initNodeAddress() common.Address {
	nodeAddress := os.Getenv(NodeAddressEnv)
	if nodeAddress == "" {
		panic("Couldn't find node address. You can set NODE_ADDRESS envvar")
	}
	return common.HexToAddress(nodeAddress)
}

func initRpConfig() *config.RocketPoolConfig {
	eth1Url := findEthClientUrl(Eth1)
	eth2Url := findEthClientUrl(Eth2)
	consensusClient := configtypes.ConsensusClient(os.Getenv(ConsensusClientEnv))
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
	return ret
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
	eth2Url := findEthClientUrl(Eth2)

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

func initEC() *services.ExecutionClientManager {
	fmt.Println("ecinit")
	ec, err := services.NewExecutionClientManager(RpConfig())
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = ec.BlockNumber(ctx)
	if err != nil {
		panic(fmt.Sprintf("Execution client probably not available at %s: %s", findEthClientUrl(Eth1), err))
	}
	return ec
}

func initRP() *rpgo.RocketPool {
	rp, err := rpgo.NewRocketPool(EC(), RocketStorageAddress())
	if err != nil {
		panic(err)
	}
	return rp
}

func initChosenFiat() string {
	fiatValue := os.Getenv(fiatEnv)
	if fiatValue == "" {
		return USD
	}
	// test

	if _, ok := XchMap[fiatValue]; !ok {
		panic(fmt.Sprintf("Unknown fiat %s", fiatValue))
	}
	return fiatValue
}

func initTelegramChatID() int64 {
	telegramChatIdStr := os.Getenv(telegramChatIdEnv)
	telegramChatId, err := strconv.ParseInt(telegramChatIdStr, 10, 64)
	if err != nil {
		panic(err)
	}
	return telegramChatId
}

func initTelegramBot() *tgbotapi.BotAPI {
	token := getEnvOrPanic(telegramTokenEnv)
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	return bot
}

