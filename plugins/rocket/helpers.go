package rocket

import (
	"github.com/jellydator/ttlcache/v3"
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
