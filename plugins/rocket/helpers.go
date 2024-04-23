package rocket

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/rewards"
	"github.com/t0mk/rocketreport/cache"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/utils"
	"github.com/t0mk/rocketreport/zaplog"
)

func GetNodeStake() (float64, error) {
	actualStake, err := node.GetNodeRPLStake(config.RP(), config.NodeAddress(), nil)
	if err != nil {
		return 0, err
	}
	actualStakeFloat, _ := utils.WeiToEther(actualStake).Float64()
	return actualStakeFloat, nil
}

func GetMinStake() (float64, error) {
	minStake, err := node.GetNodeMinimumRPLStake(config.RP(), config.NodeAddress(), nil)
	if err != nil {
		return 0., err
	}
	minStakeFloat, _ := utils.WeiToEther(minStake).Float64()
	return minStakeFloat, nil
}

func GetEthMatched() (float64, error) {
	ethMatched, err := node.GetNodeEthMatched(config.RP(), config.NodeAddress(), nil)
	if err != nil {
		return 0, err
	}
	ethMatchedFloat, _ := utils.WeiToEther(ethMatched).Float64()
	return ethMatchedFloat, nil
}

func GetStakeRatio() (float64, error) {
	price, err := cache.Float("rpOracleRplPrice", GetRplEthOraclePrice)
	if err != nil {
		return 0., err
	}
	ethMatched, err := cache.Float("rpEthMatched", GetEthMatched)
	if err != nil {
		return 0., err
	}
	nodeStake, err := cache.Float("rpNodeStake", GetNodeStake)
	if err != nil {
		return 0., err
	}
	ethRatio := nodeStake / (ethMatched / price)
	return ethRatio * 100, nil
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

func GetRplEthOraclePrice() (float64, error) {
	rplPrice, err := network.GetRPLPrice(config.RP(), nil)
	if err != nil {
		return 0, err
	}
	floatRplPrice, _ := utils.WeiToEther(rplPrice).Float64()
	return floatRplPrice, nil
}

func GetRplPriceBlock() (uint64, error) {
	rplPriceBlock, err := network.GetPricesBlock(config.RP(), nil)
	if err != nil {
		return 0, err
	}

	// Return the price
	return rplPriceBlock, nil
}

func GetNextRplPriceUpdate() (time.Time, error) {
	c := context.Background()
	latestBlock, err := config.EC().BlockNumber(c)
	if err != nil {
		return time.Time{}, err
	}
	now := time.Now()
	lastRplPriceUpdateBlock, err := GetRplPriceBlock()
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
	return utils.AddressBallance(fdAddr)
}

func GetNodeBalance() (float64, error) {
	return utils.AddressBallance(config.NodeAddress())
}

func GetIntervalEnd() (time.Time, error) {
	start, err := rewards.GetClaimIntervalTimeStart(config.RP(), nil)
	if err != nil {
		return time.Time{}, err
	}
	duration, err := rewards.GetClaimIntervalTime(config.RP(), nil)
	if err != nil {
		return time.Time{}, err
	}
	return start.Add(duration).UTC(), nil
}

func GetWithdrawalAddressBallance() (float64, error) {
	nd, err := node.GetNodeDetails(config.RP(), config.NodeAddress(), nil)
	if err != nil {
		return 0, err
	}
	return utils.AddressBallance(nd.WithdrawalAddress)
}
