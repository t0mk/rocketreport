package utils

import (
	"context"
	"fmt"

	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/t0mk/rocketreport/config"
)

func SmoothingPoolBalance() (*float64, error) {
	// Get the Smoothing Pool contract's balance
	smoothingPoolContract, err := config.RP().GetContract("rocketSmoothingPool", nil)
	if err != nil {
		return nil, fmt.Errorf("error getting smoothing pool contract: %w", err)
	}
	addr := *smoothingPoolContract.Address

	smoothingPoolBalance, err := config.RP().Client.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting smoothing pool balance: %w", err)
	}

	balance := eth.WeiToEth(smoothingPoolBalance)
	return &balance, nil
}

/*
func SmoothingPoolMiniPoolCount() (*uint64, error) {
	// Get the Smoothing Pool contract's minipool count
	rpnm, err := config.RP.GetContract("rocketNodeManager", nil)
	if err != nil {
		return nil, fmt.Errorf("error getting node manager contract: %w", err)
	}



	err := rpnm.Call(nil, "getNodeDetails", nil)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool count: %w", err)
	}

	return &miniPoolCount, nil
}


// Get a node address by index
func GetNodeAt(rp *rocketpool.RocketPool, index uint64, opts *bind.CallOpts) (common.Address, error) {
	rocketNodeManager, err := getRocketNodeManager(rp, opts)
	if err != nil {
		return common.Address{}, err
	}
	nodeAddress := new(common.Address)
	if err := rocketNodeManager.Call(opts, nodeAddress, "getNodeAt", big.NewInt(int64(index))); err != nil {
		return common.Address{}, fmt.Errorf("Could not get node %d address: %w", index, err)
	}
	return *nodeAddress, nil
}
*/
