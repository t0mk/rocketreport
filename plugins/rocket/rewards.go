package rocket

import (
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/t0mk/rocketreport/config"

	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/minipool"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/rewards"
	"github.com/rocket-pool/rocketpool-go/tokens"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	rpstate "github.com/rocket-pool/rocketpool-go/utils/state"
	"github.com/rocket-pool/smartnode/shared/services/beacon"
	rprewards "github.com/rocket-pool/smartnode/shared/services/rewards"
	"github.com/rocket-pool/smartnode/shared/utils/eth2"
	"golang.org/x/sync/errgroup"
)

type NodeRewardsResponse struct {
	Status               string        `json:"status"`
	Error                string        `json:"error"`
	NodeRegistrationTime time.Time     `json:"nodeRegistrationTime"`
	RewardsInterval      time.Duration `json:"rewardsInterval"`
	LastCheckpoint       time.Time     `json:"lastCheckpoint"`
	Registered           bool          `json:"registered"`
	EffectiveRplStake    float64       `json:"effectiveRplStake"`
	TotalRplStake        float64       `json:"totalRplStake"`
	EstimatedRewards     float64       `json:"estimatedRewards"`
	CumulativeRplRewards float64       `json:"cumulativeRplRewards"`
	CumulativeEthRewards float64       `json:"cumulativeEthRewards"`
	UnclaimedRplRewards  float64       `json:"unclaimedRplRewards"`
	UnclaimedEthRewards  float64       `json:"unclaimedEthRewards"`
	BeaconRewards        float64       `json:"beaconRewards"`
}

func GetIntervalInfo(interval uint64) (info rprewards.IntervalInfo, err error) {
	info.Index = interval
	var event rewards.RewardsEvent
	cfg := config.RpConfig()

	// Get the event details for this interval
	event, err = rprewards.GetRewardSnapshotEvent(config.RP(), cfg, interval, nil)
	if err != nil {
		return
	}

	info.CID = event.MerkleTreeCID
	info.StartTime = event.IntervalStartTime
	info.EndTime = event.IntervalEndTime
	merkleRootCanon := event.MerkleRoot
	info.MerkleRoot = merkleRootCanon

	// Check if the tree file exists
	//info.TreeFilePath = cfg.Smartnode.GetRewardsTreePath(interval, true)
	n := string(config.Network())
	info.TreeFilePath = fmt.Sprintf("./rewards-trees/%s/rp-rewards-%s-%d.json", n, n, interval)
	_, err = os.Stat(info.TreeFilePath)
	if os.IsNotExist(err) {
		info.TreeFileExists = false
		err = nil
		return
	}
	info.TreeFileExists = true

	// Unmarshal it
	localRewardsFile, err := rprewards.ReadLocalRewardsFile(info.TreeFilePath)
	if err != nil {
		err = fmt.Errorf("error reading %s: %w", info.TreeFilePath, err)
		return
	}

	proofWrapper := localRewardsFile.Impl()

	// Make sure the Merkle root has the expected value
	merkleRootFromFile := common.HexToHash(proofWrapper.GetHeader().MerkleRoot)
	if merkleRootCanon != merkleRootFromFile {
		info.MerkleRootValid = false
		return
	}
	info.MerkleRootValid = true

	// Get the rewards from it
	rewards, exists := proofWrapper.GetNodeRewardsInfo(config.NodeAddress())
	info.NodeExists = exists
	if exists {
		info.CollateralRplAmount = rewards.GetCollateralRpl()
		info.ODaoRplAmount = rewards.GetOracleDaoRpl()
		info.SmoothingPoolEthAmount = rewards.GetSmoothingPoolEth()

		var proof []common.Hash
		proof, err = rewards.GetMerkleProof()
		if err != nil {
			err = fmt.Errorf("error deserializing merkle proof for %s, node %s: %w", info.TreeFilePath, config.NodeAddress().Hex(), err)
			return
		}
		info.MerkleProof = proof
	}

	return
}

func GetRewards() (*NodeRewardsResponse, error) {

	rp := config.RP()
	bc := config.BC()
	cfg := config.RpConfig()
	nodeAddress := config.NodeAddress()

	// Response
	response := NodeRewardsResponse{}

	var totalEffectiveStake *big.Int
	var totalRplSupply *big.Int
	var inflationInterval *big.Int
	var nodeOperatorRewardsPercent float64
	var totalDepositBalance float64
	var totalNodeShare float64
	var addresses []common.Address
	var beaconHead beacon.BeaconHead

	// Sync
	var wg errgroup.Group

	// Check if the node is registered or not
	wg.Go(func() error {
		exists, err := node.GetNodeExists(rp, nodeAddress, nil)
		if err == nil {
			response.Registered = exists
		}
		return err
	})

	// Get the node registration time
	wg.Go(func() error {
		var time time.Time
		var err error
		time, err = node.GetNodeRegistrationTime(rp, nodeAddress, nil)

		if err == nil {
			response.NodeRegistrationTime = time
		}
		return err
	})

	// Get claimed and pending rewards
	wg.Go(func() error {
		// Legacy rewards
		unclaimedRplRewardsWei := big.NewInt(0)
		rplRewards := big.NewInt(0)
		// TEMP removal of the legacy rewards crawler for now, TODO performance improvements here
		/*
			rplRewards, err := legacyrewards.CalculateLifetimeNodeRewards(rp, nodeAddress, big.NewInt(int64(eventLogInterval)), nil, &legacyRocketRewardsAddress, &legacyClaimNodeAddress)*/
		unclaimedEthRewardsWei := big.NewInt(0)
		ethRewards := big.NewInt(0)

		// Get the claimed and unclaimed intervals
		unclaimed, claimed, err := rprewards.GetClaimStatus(rp, nodeAddress)
		if err != nil {
			return err
		}

		// Get the info for each claimed interval
		for _, claimedInterval := range claimed {
			//intervalInfo, err := rprewards.GetIntervalInfo(rp, cfg, nodeAddress, claimedInterval, nil)
			intervalInfo, err := GetIntervalInfo(claimedInterval)
			if err != nil {
				return err
			}
			if !intervalInfo.TreeFileExists {
				return fmt.Errorf("Error calculating lifetime node rewards: rewards file %s doesn't exist but interval %d was claimed", intervalInfo.TreeFilePath, claimedInterval)
			}
			rplRewards.Add(rplRewards, &intervalInfo.CollateralRplAmount.Int)
			ethRewards.Add(ethRewards, &intervalInfo.SmoothingPoolEthAmount.Int)
		}

		// Get the unclaimed rewards
		for _, unclaimedInterval := range unclaimed {
			//intervalInfo, err := rprewards.GetIntervalInfo(rp, cfg, nodeAddress, unclaimedInterval, nil)
			intervalInfo, err := GetIntervalInfo(unclaimedInterval)
			if err != nil {
				return err
			}
			if !intervalInfo.TreeFileExists {
				return fmt.Errorf("Error calculating lifetime node rewards: rewards file %s doesn't exist and interval %d is unclaimed", intervalInfo.TreeFilePath, unclaimedInterval)
			}
			if intervalInfo.NodeExists {
				unclaimedRplRewardsWei.Add(unclaimedRplRewardsWei, &intervalInfo.CollateralRplAmount.Int)
				unclaimedEthRewardsWei.Add(unclaimedEthRewardsWei, &intervalInfo.SmoothingPoolEthAmount.Int)
			}
		}

		if err == nil {
			response.CumulativeRplRewards = eth.WeiToEth(rplRewards)
			response.UnclaimedRplRewards = eth.WeiToEth(unclaimedRplRewardsWei)
			response.CumulativeEthRewards = eth.WeiToEth(ethRewards)
			response.UnclaimedEthRewards = eth.WeiToEth(unclaimedEthRewardsWei)
		}
		return err
	})

	// Get the start of the rewards checkpoint
	wg.Go(func() error {
		lastCheckpoint, err := rewards.GetClaimIntervalTimeStart(rp, nil)
		if err == nil {
			response.LastCheckpoint = lastCheckpoint
		}
		return err
	})

	// Get the rewards checkpoint interval
	wg.Go(func() error {
		rewardsInterval, err := rewards.GetClaimIntervalTime(rp, nil)
		if err == nil {
			response.RewardsInterval = rewardsInterval
		}
		return err
	})

	// Get the node's effective stake
	wg.Go(func() error {
		effectiveStake, err := node.GetNodeEffectiveRPLStake(rp, nodeAddress, nil)
		if err == nil {
			response.EffectiveRplStake = eth.WeiToEth(effectiveStake)
		}
		return err
	})

	// Get the node's total stake
	wg.Go(func() error {
		stake, err := node.GetNodeRPLStake(rp, nodeAddress, nil)
		if err == nil {
			response.TotalRplStake = eth.WeiToEth(stake)
		}
		return err
	})

	// Get the total network effective stake
	wg.Go(func() error {
		multicallerAddress := common.HexToAddress(cfg.Smartnode.GetMulticallAddress())
		balanceBatcherAddress := common.HexToAddress(cfg.Smartnode.GetBalanceBatcherAddress())
		contracts, err := rpstate.NewNetworkContracts(rp, multicallerAddress, balanceBatcherAddress, nil)
		if err != nil {
			return fmt.Errorf("error creating network contract binding: %w", err)
		}
		totalEffectiveStake, err = rpstate.GetTotalEffectiveRplStake(rp, contracts)
		if err != nil {
			return fmt.Errorf("error getting total effective RPL stake: %w", err)
		}
		return nil
	})

	// Get the total RPL supply
	wg.Go(func() error {
		var err error
		totalRplSupply, err = tokens.GetRPLTotalSupply(rp, nil)
		if err != nil {
			return err
		}
		return nil
	})

	// Get the RPL inflation interval
	wg.Go(func() error {
		var err error
		inflationInterval, err = tokens.GetRPLInflationIntervalRate(rp, nil)
		if err != nil {
			return err
		}
		return nil
	})

	// Get the node operator rewards percent
	wg.Go(func() error {
		nodeOperatorRewardsPercentRaw, err := rewards.GetNodeOperatorRewardsPercent(rp, nil)
		nodeOperatorRewardsPercent = eth.WeiToEth(nodeOperatorRewardsPercentRaw)
		if err != nil {
			return err
		}
		return nil
	})

	// Get the list of minipool addresses for this node
	wg.Go(func() error {
		_addresses, err := minipool.GetNodeMinipoolAddresses(rp, nodeAddress, nil)
		if err != nil {
			return fmt.Errorf("Error getting node minipool addresses: %w", err)
		}
		addresses = _addresses
		return nil
	})

	// Get the beacon head
	wg.Go(func() error {
		_beaconHead, err := bc.GetBeaconHead()
		if err != nil {
			return fmt.Errorf("Error getting beacon chain head: %w", err)
		}
		beaconHead = _beaconHead
		return nil
	})

	// Wait for data
	if err := wg.Wait(); err != nil {
		return nil, err
	}

	// Calculate the total deposits and corresponding beacon chain balance share
	minipoolDetails, err := eth2.GetBeaconBalances(rp, bc, addresses, beaconHead, nil)
	if err != nil {
		return nil, err
	}
	for _, minipool := range minipoolDetails {
		totalDepositBalance += eth.WeiToEth(minipool.NodeDeposit)
		totalNodeShare += eth.WeiToEth(minipool.NodeBalance)
	}
	response.BeaconRewards = totalNodeShare - totalDepositBalance

	// Calculate the estimated rewards
	rewardsIntervalDays := response.RewardsInterval.Seconds() / (60 * 60 * 24)
	inflationPerDay := eth.WeiToEth(inflationInterval)
	totalRplAtNextCheckpoint := (math.Pow(inflationPerDay, float64(rewardsIntervalDays)) - 1) * eth.WeiToEth(totalRplSupply)
	if totalRplAtNextCheckpoint < 0 {
		totalRplAtNextCheckpoint = 0
	}

	if totalEffectiveStake.Cmp(big.NewInt(0)) == 1 {
		response.EstimatedRewards = response.EffectiveRplStake / eth.WeiToEth(totalEffectiveStake) * totalRplAtNextCheckpoint * nodeOperatorRewardsPercent
	}

	// Return response
	return &response, nil

}
