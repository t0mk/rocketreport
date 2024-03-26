package main

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/rewards"
	"github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils/state"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/plugins"
	"github.com/t0mk/rocketreport/prices"
	"github.com/t0mk/rocketreport/zaplog"

	"github.com/schollz/progressbar/v3"
)

type Context struct {
	Debug bool
}

type ListPluginsCmd struct {
	Eval bool `short:"e" help:"Evaluate all plugins"`
}

func (l *ListPluginsCmd) Run(ctx *Context) error {
	for _, p := range plugins.Plugins {
		line := fmt.Sprintf("%-20s %-20s", p.Key, p.Desc)
		if l.Eval {
			p.Eval()
			if p.Err != "" {
				line += fmt.Sprintf(" (%s)", p.Err)
			} else {
				line += fmt.Sprintf(" (%s)", p.Output)
			}
		}
		fmt.Println(line)
	}
	return nil
}

type SendCmd struct {
	DoSend bool `short:"s" help:"Send to Telegram"`
}

func (s *SendCmd) Run(ctx *Context) error {
	log := zaplog.New()
	ethFiat, err := prices.PriEth(config.ChosenFiat)
	if err != nil {
		log.Error("Error getting eth price", err)
	}
	suff := fmt.Sprintf("%s/Ξ", config.ChosenFiat.String())
	ethFiatStr := plugins.FloatSuffixFormatter(0, suff)(ethFiat)

	fmt.Println("sending")
	if s.DoSend && config.Bot != nil {
		ts := time.Now().Format("Mon 02-Jan 15:04")
		subj := fmt.Sprintf("%s - %s", ts, ethFiatStr)
		nm := plugins.Plugins.TelegramFormat(subj)
		_, err := config.Bot.Send(nm)
		return err
	} else {
		fmt.Println("Not sending to Telegram, use -s to send.")
		txt := plugins.ToPlaintext(plugins.Plugins)
		fmt.Println(txt)
	}
	return nil
}

type PrintCmd struct{}

type TCmd struct{}

func (t *TCmd) Run(ctx *Context) error {

	nodeTimeInInterval := map[common.Address]time.Duration{}

	claimIntervalStart, err := rewards.GetClaimIntervalTimeStart(config.RP, nil)
	if err != nil {
		return err
	}
	log := zaplog.New()
	nas, err := node.GetNodeAddresses(config.RP, nil)
	if err != nil {
		return err
	}
	smnas := make([]common.Address, 0)

	pbn := progressbar.Default(int64(len(nas)))

	for _, na := range nas {
		isIn, err := node.GetSmoothingPoolRegistrationState(config.RP, na, nil)
		if err != nil {
			log.Error("Error getting smoothing pool registration state", err)
			continue
		}
		spChanged, err := node.GetSmoothingPoolRegistrationChanged(config.RP, na, nil)
		if err != nil {
			log.Error("Error getting smoothing pool changed", err)
			continue
		}
		if !spChanged.IsZero() {
			if isIn {
				smnas = append(smnas, na)
				nodeTimeInInterval[na] = time.Since(claimIntervalStart)
			} else {
				if spChanged.After(claimIntervalStart) {
					smnas = append(smnas, na)
					nodeTimeInInterval[na] = spChanged.Sub(claimIntervalStart)
				}
			}
		}

		if spChanged.IsZero() {
			if isIn {
				fmt.Println(na, "never changed and isIn")
			}
		}
		pbn.Add(1)
		nodeTimeInInterval[na] = time.Since(claimIntervalStart)
		//fmt.Println("spChanged", spChanged)
	}

	fmt.Println("totale nodes", len(nas))
	fmt.Println("smoothing nodes", len(smnas))

	oks := 0
	fails := 0
	okmps := 0
	failmps := 0

	statusCounts := map[types.MinipoolStatus]int{}
	depositCounts := map[types.MinipoolDeposit]int{}

	feeCount := 0
	feeSum := 0.0

	multicallerAddress := common.HexToAddress(config.SnConfig.GetMulticallAddress())
	balanceBatcherAddress := common.HexToAddress(config.SnConfig.GetBalanceBatcherAddress())
	contracts, err := state.NewNetworkContracts(config.RP, multicallerAddress, balanceBatcherAddress, nil)
	if err != nil {
		return err
	}

	pb := progressbar.Default(int64(len(smnas)))
	for _, na := range smnas {
		mps, err := state.GetNodeNativeMinipoolDetails(config.RP, contracts, na)
		if err != nil {
			log.Error("Error getting minipool details", err)
			failmps++
			continue
		}
		for _, mp := range mps {
			if na == config.NodeAddress {
				time := time.Unix(mp.StatusTime.Int64(), 0)
				fmt.Println(na, time)
			}
			statusCounts[mp.Status]++
			depositCounts[mp.DepositType]++
			fmt.Println(mp.NodeFee.Div()
			/*
				mp, err := minipool.NewMinipool(config.RP, a, nil)
				if err != nil {
					log.Error("Error getting minipool deposit type", err)
					failmps++
					continue
				}
				fee, err := mp.GetNodeFee(nil)
				if err != nil {
					log.Error("Error getting minipool fee", err)
					failmps++
					continue
				}
				feeSum += fee
				feeCount++

				sds, err := mp.GetStatusDetails(nil)
				if err != nil {
					log.Error("Error getting minipool status details", err)
					continue
				}
				dt, err := mp.GetDepositType(nil)
				if err != nil {
					log.Error("Error getting minipool deposit type", err)
					continue
				}
			*/
		}
		time.Sleep(7 * time.Millisecond)
		pb.Add(1)
	}
	fmt.Println("oks", oks)
	fmt.Println("fails", fails)
	fmt.Println("okmps", okmps)
	fmt.Println("failmps", failmps)

	fmt.Println("==========================")
	fmt.Println("Status counts")
	fmt.Println("==========================")
	for k, v := range statusCounts {
		fmt.Printf("%s: %d\n", k, v)
	}
	fmt.Println("==========================")
	fmt.Println("Deposit counts")
	fmt.Println("==========================")
	for k, v := range depositCounts {
		fmt.Printf("%s: %d\n", k, v)
	}
	fmt.Println("==========================")
	fmt.Println("Fee sum")
	fmt.Println("==========================")
	fmt.Printf("%f\n", feeSum)
	fmt.Println("==========================")
	fmt.Println("Fee count")
	fmt.Println("==========================")
	fmt.Printf("%d\n", feeCount)
	fmt.Println("==========================")
	fmt.Println("Fee avg")
	fmt.Println("==========================")
	fmt.Printf("%f\n", feeSum/float64(feeCount))
	return nil
}

type MyMinipoolDetails struct {
	Status      types.MinipoolStatus
	depositType types.MinipoolDeposit
}

var cli struct {
	Send  SendCmd        `cmd:"" help:"Send to configured telegram chat"`
	Print PrintCmd       `cmd:"" help:"Print to stdout"`
	List  ListPluginsCmd `cmd:"" help:"List all plugins"`
	T     TCmd           `cmd:"" help:"T"`
}

func (p *PrintCmd) Run(ctx *Context) error {
	fmt.Println(plugins.Plugins.TermText())
	return nil
}

func main() {
	config.Setup()
	plugins.RegisterAll()
	ctx := kong.Parse(&cli)
	err := ctx.Run(&Context{Debug: true})
	ctx.FatalIfErrorf(err)
	kong.Parse(&cli)
}
