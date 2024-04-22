package rocket

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/t0mk/rocketreport/cache"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/utils"
	"github.com/t0mk/rocketreport/zaplog"
)

func GetActualStake(...interface{}) (interface{}, error) {
	actualStake, err := node.GetNodeRPLStake(config.RP(), config.NodeAddress(), nil)
	if err != nil {
		return nil, err
	}
	actualStakeFloat, _ := utils.WeiToEther(actualStake).Float64()
	return actualStakeFloat, nil
}

func GetMinStake(...interface{}) (interface{}, error) {
	minStake, err := node.GetNodeMinimumRPLStake(config.RP(), config.NodeAddress(), nil)
	if err != nil {
		return nil, err
	}
	minStakeFloat, _ := utils.WeiToEther(minStake).Float64()
	return minStakeFloat, nil
}

func CachedGetMinipoolDetails(cacheKey string) (*utils.MinipoolInterestingDetails, error) {
	log := zaplog.New()
	item := cache.Cache.Get(cacheKey)
	if (item != nil) && (!item.IsExpired()) {
		mpd := item.Value().(utils.MinipoolInterestingDetails)
		return &mpd, nil
	}
	log.Debug("Cache miss for ", cacheKey)
	details, err := utils.GetMinipoolInterestingDetails(
		config.RP(),
		config.BC(),
		config.NodeAddress(),
		nil,
	)
	log.Debug("Got details ", details)
	if err != nil {
		return nil, err
	}
	cache.Cache.Set("minipoolDetails", details, ttlcache.DefaultTTL)
	return &details, nil
}

func RplEthOraclePriceCached() (float64, error) {
	priceRaw, ok := cache.Get("rplEthOraclePrice")

	if ok {
		return priceRaw.(float64), nil
	}
	price, err := RplEthOraclePrice()
	if err != nil {
		return 0, err
	}
	cache.Set("rplEthOraclePrice", price)
	return price, nil
}

func RplEthOraclePrice() (float64, error) {
	rplPrice, err := network.GetRPLPrice(config.RP(), nil)
	if err != nil {
		return 0, err
	}
	floatRplPrice, _ := utils.WeiToEther(rplPrice).Float64()
	return floatRplPrice, nil
}

func RplPriceBlock() (uint64, error) {
	rplPriceBlock, err := network.GetPricesBlock(config.RP(), nil)
	if err != nil {
		return 0, err
	}

	// Return the price
	return rplPriceBlock, nil
}

func NextRplPriceUpdate() (time.Time, error) {
	c := context.Background()
	latestBlock, err := config.EC().BlockNumber(c)
	if err != nil {
		return time.Time{}, err
	}
	now := time.Now()
	lastRplPriceUpdateBlock, err := RplPriceBlock()
	if err != nil {
		return time.Time{}, err
	}
	elapsedBlocks := latestBlock - lastRplPriceUpdateBlock
	nextPriceUpdateBlock := 5760 - elapsedBlocks
	nextPriceUpdateTime := now.Add(time.Duration(nextPriceUpdateBlock) * 12 * time.Second)
	return nextPriceUpdateTime, nil
}

func GetFeeDistributorBalance() (float64, error) {
	fdAddr, err := node.GetDistributorAddress(config.RP(), config.NodeAddress(), nil)
	if err != nil {
		return 0, err
	}
	balanceRaw, err := config.RP().Client.BalanceAt(context.Background(), fdAddr, nil)
	if err != nil {
		return 0, err
	}
	balance, _ := utils.WeiToEther(balanceRaw).Float64()
	return balance, nil
}
