
# Rocketreport Plugins

- [Rocketpool Plugins
](#rocketpool-plugins)
- [Exchange Plugins
](#exchange-plugins)
- [Meta Plugins
](#meta-plugins)
- [Common Plugins
](#common-plugins)
## Rocketpool Plugins
| Name | Description |
|------|-------------|
| rpBeaconRewards | Beacon rewards of Rocketpool node |
| rpCumulativeEthRewards | Cumulative ETH rewards of Rocketpool node |
| rpCumulativeRplRewards | Cumulative RPL rewards of Rocketpool node |
| rpEarnedConsesusEth | Amount of consensus ETH in USDT* |
| rpEffectiveRplStake | Effective RPL stake of Rocketpool node |
| rpEstimatedRewards | Estimated rewards of Rocketpool node |
| rpEth1sync | Sync status of Eth1 client (with Rocketpool Golang library) |
| rpEth2sync | Sync status of Eth2 client (with Rocketpool Golang library) |
| rpEthMatched | Matched ETH of Rocketpool node |
| rpFeeDistributorBalance | Balance of the Rocketpool fee distributor |
| rpIntervalEnd | End of the current Rocketpool interval |
| rpMinStake | Minimum RPL stake for Rocketpool node |
| rpNodeBalance | Balance of the Rocketpool node |
| rpNodeStake | RPL stake of Rocketpool node |
| rpOracleRplPrice | RPL price from Rocketpool oracle |
| rpOracleRplPriceUpdate | Time of next RPL price update in Rocketpool oracle |
| rpOwnEthDeposit | Amount of ETH deposited in Rocketpool node |
| rpStakeRatio | How much % of the borrowed Eth value is staked |
| rpTotalRplStake | Total RPL stake of Rocketpool node |
| rpUnclaimedEthRewards | Unclaimed ETH rewards of Rocketpool node |
| rpUnclaimedRplRewards | Unclaimed RPL rewards of Rocketpool node |
| rpUntilIntervalEnd | Time until the end of the current Rocketpool interval |
| rpUntilOracleRplPriceUpdate | Time until next RPL price update in Rocketpool oracle |
| rpWithdrawalAddressBalance | Balance of the Rocketpool withdrawal address |


## Exchange Plugins
| Name | Description | Args (type) | Defaults |
|------|-------------|------|--------------|
| binance | Latest ticker price from Binance | ticker (string), amount (float64) | ETHUSDT, 1 |
| bitfinex | Latest ticker price from Bitfinex | ticker (string), amount (float64) | ETHEUR, 1 |
| coinmate | Latest ticker price from Coinmate | ticker (string), amount (float64) | ETH_EUR, 1 |
| kraken | Latest ticker price from Kraken | ticker (string), amount (float64) | XETHZEUR, 1 |


## Meta Plugins
| Name | Description | Args (type) | Defaults |
|------|-------------|------|--------------|
| add | Sum of given args, either numbers or plugin outputs, adds args and outputs a float | list of values - numbers or plugin outputs ([]interface {}) | [] |
| div | Divide first arg by second, either numbers or plugin outputs, divides args and outputs a float | list of values - numbers or plugin outputs ([]interface {}) | [] |
| mul | Product of given args, either numbers or plugin outputs, multiplies args and outputs a float | list of values - numbers or plugin outputs ([]interface {}) | [] |
| sub | Subtract second arg from first, either numbers or plugin outputs, subtracts args and outputs a float | list of values - numbers or plugin outputs ([]interface {}) | [] |


## Common Plugins
| Name | Description | Args (type) | Defaults |
|------|-------------|------|--------------|
| addressBalance | Balance of an address via Execution client | address (string) |  |
| addressBalanceEtherscan | Balance of an address using Etherscan | address (string) |  |
| date | Current date |  |  |
| ethPrice | ETH/USDT* price |  |  |
| gasPriceBeaconcha.in | Latest gas price from beaconcha.in |  |  |
| gasPriceExecutionClient | Latest gas price from the execution client |  |  |
| rplPriceRealtime | Realtime RPL-ETH (based on RPL-USDT and ETH-USDT from Binance) |  |  |
| timeMin | Current time up to minutes |  |  |
| timeSec | Current time up to seconds |  |  |




&ast; you can use different fiat as quote currency in these plugins if you set "fiat" option in config.yml