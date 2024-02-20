package main

import (
	"fmt"

	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/alecthomas/kong"
)

func tabPrint(s string, v ...interface{}) {
	fmt.Printf("%-30s", s)
	for _, i := range v {
		fmt.Printf("%v ", i)
	}
	fmt.Println()
}

var cli struct {
	Send struct {
		DoSend bool `short:"s" help:"Send to configured telegram chat" default:"true"`
	} `cmd:"" help:Send to configured telegram chat"`
	Print struct {
	} `cmd:"" help:Print to stdout"`
}

func main() {
	doConfig()

	ecs := ec.CheckStatus(rpConfig)
	bcs := bc.CheckStatus()

	fmt.Println("Eth1 client", EthClientStatusString(ecs.PrimaryClientStatus))
	fmt.Println("Eth2 client", EthClientStatusString(bcs.PrimaryClientStatus))
	details, err := GetMinipoolInterestingDetails(
		rp,
		bc,
		nodeAddress,
		nil,
	)
	if err != nil {
		panic(err)
	}

	minStake, err := node.GetNodeMinimumRPLStake(rp, nodeAddress, nil)
	if err != nil {
		panic(err)
	}
	actualStake, err := node.GetNodeRPLStake(rp, nodeAddress, nil)
	if err != nil {
		panic(err)
	}

	actualFloatStake, _ := weiToEther(actualStake).Float64()
	normalFloatMinStake, _ := weiToEther(minStake).Float64()
	stakeReserve := actualFloatStake - normalFloatMinStake

	rplPrice, err := priRplEth()
	if err != nil {
		panic(err)
	}
	fiat := chosenFiat
	rplFiat, err := priRplFiat(fiat)
	if err != nil {
		panic(err)
	}
	ethFiat, err := priEthFiat(fiat)
	if err != nil {
		panic(err)
	}
	rplFiatAmount := rplFiat * actualFloatStake
	ethFiatAmount := ethFiat * (details.TotalDeposit + details.Earned)

	tabPrint("Min stake (10%): ", fmtRpl(normalFloatMinStake), "RPL")
	tabPrint("Max stake (150%): ", fmtRpl(actualFloatStake), "RPL")
	tabPrint("Stake over min:", fmtRpl(stakeReserve), "RPL")

	tabPrint("RPL/ETH", fmtEth(rplPrice), "ETH")
	tabPrint("ETH/"+fiat.String(), fmtFiat(ethFiat), fiat.String())
	tabPrint("RPL/"+fiat.String(), fmtRplEur(rplFiat), fiat.String())
	tabPrint("Funds RPL", fmtFiat(rplFiatAmount), fiat.String())
	tabPrint("Funds ETH", fmtFiat(ethFiatAmount), fiat.String())
	tabPrint("Funds Total", fmtFiat(rplFiatAmount+ethFiatAmount), fiat.String())

}
