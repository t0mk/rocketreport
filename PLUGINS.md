
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
| rpBeaconRewards | Check the beacon rewards of Rocketpool node |
| rpCumulativeEthRewards | Check the cumulative ETH rewards of Rocketpool node |
| rpCumulativeRplRewards | Check the cumulative RPL rewards of Rocketpool node |
| rpEarnedConsesusEth | Check the amount of consensus ETH in USDT* |
| rpEffectiveRplStake | Check the effective RPL stake of Rocketpool node |
| rpEstimatedRewards | Check the estimated rewards of Rocketpool node |
| rpEth1sync | Check the sync status of Eth1 client (with Rocketpool Golang library) |
| rpEth2sync | Check the sync status of Eth2 client (with Rocketpool Golang library) |
| rpEthMatched | Check the matched ETH of Rocketpool node |
| rpFeeDistributorBalance | Check the balance of the Rocketpool fee distributor |
| rpIntervalEnd | Check the end of the current Rocketpool interval |
| rpMinStake | Check the minimum RPL stake for Rocketpool node |
| rpNodeBalance | Check the balance of the Rocketpool node |
| rpNodeStake | Check the RPL stake of Rocketpool node |
| rpOracleRplPrice | Check the RPL price from Rocketpool oracle |
| rpOracleRplPriceUpdate | Time of next RPL price update in Rocketpool oracle |
| rpOwnEthDeposit | Check the amount of ETH deposited in Rocketpool node |
| rpStakeRatio | Check how much % of the borrowed Eth value is staked |
| rpTotalRplStake | Check the total RPL stake of Rocketpool node |
| rpUnclaimedEthRewards | Check the unclaimed ETH rewards of Rocketpool node |
| rpUnclaimedRplRewards | Check the unclaimed RPL rewards of Rocketpool node |
| rpUntilIntervalEnd | Check the time until the end of the current Rocketpool interval |
| rpUntilOracleRplPriceUpdate | Time until next RPL price update in Rocketpool oracle |
| rpWithdrawalAddressBalance | Check the balance of the Rocketpool withdrawal address |


## Exchange Plugins
| Name | Description | Args | Defaults |
|------|-------------|------|--------------|
| binance | Get the latest ticker price from Binance | ticker (string), amount (float64) | ETHUSDT, 1 |
| bitfinex | Get the latest ticker price from Bitfinex | ticker (string), amount (float64) | ETHEUR, 1 |
| coinmate | Get the latest ticker price from Coinmate | ticker (string), amount (float64) | ETH_EUR, 1 |
| kraken | Get the latest ticker price from Kraken | ticker (string), amount (float64) | XETHZEUR, 1 |


## Meta Plugins
| Name | Description | Args | Defaults |
|------|-------------|------|--------------|
| add | Sum of given args, either numbers or plugin outputs, adds args and outputs a float | list of values - numbers or plugin outputs ([]interface {}) | [] |
| div | Divide first arg by second, either numbers or plugin outputs, divides args and outputs a float | list of values - numbers or plugin outputs ([]interface {}) | [] |
| mul | Product of given args, either numbers or plugin outputs, multiplies args and outputs a float | list of values - numbers or plugin outputs ([]interface {}) | [] |
| sub | Subtract second arg from first, either numbers or plugin outputs, subtracts args and outputs a float | list of values - numbers or plugin outputs ([]interface {}) | [] |


## Common Plugins
| Name | Description | Args | Defaults |
|------|-------------|------|--------------|
| addressBalance | Check the balance of an address | address (string) |  |
| ethPrice | Check ETH/USDT* price |  |  |
| gasPriceBeaconChain | Get the latest gas price from beaconcha.in |  |  |
| gasPriceExecutionClient | Get the latest gas price from the execution client |  |  |
| rplPriceRealtime | Check realtime RPL-ETH (based on RPL-USDT and ETH-USDT from Binance) |  |  |




&ast; you can use different fiat as quote currency in these plugins if you set "fiat" option in config.yml