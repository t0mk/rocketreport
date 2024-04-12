# Rocketreport

Rocketreport is a tool that fetches stats about Rocketpool node and other crypto data and sends it to you over Telegram. 
- It can be configured in cron-like fashing and send messages to existing chat
- It can serve a Telegram bot and reply on demand

## Basic usage

You can run it with Docker. To print current Ethereum gas price, do:

```
docker run --rm t0mk/rocketreport gasPrice
```

## Install binaries

In case you don't want to use Docker, you can get release artifacts from GitHub.
To install latest release for Linux:

```sh
wget -O /tmp/rocketreport https://github.com/t0mk/rocketreport/releases/latest/download/rocketreport-linux-amd64 && chmod +x /tmp/rocketreport && sudo cp /tmp/rocketreport /usr/local/bin/
```

.. for MacOS:

```sh
wget -O /tmp/rocketreport https://github.com/t0mk/rocketreport/releases/latest/download/rocketreport-darwin-amd64 && chmod +x /tmp/rocketreport && sudo cp /tmp/rocketreport /usr/local/bin/
```

## Plugins

Rocketreport messages are compiled from plugin output. That way you can configure what info you want to see in your messages.

| Name | Description | Args | Example Args |
|------|-------------|------|--------------|
| actualStake | Check actual RPL stake of Rocketpool node |  |  |
| bitfinex | Get the latest ticker price from Bitfinex | ticker (string), amount (float64) | ETHEUR, 1 |
| coinmate | Get the latest ticker price from Coinmate | ticker (string), amount (float64) | ETH_EUR, 1 |
| depositedEthFiat | Check the amount of deposited ETH in USD* |  |  |
| earnedConsensusFunds | Check the amount of consensus funds in USD* |  |  |
| earnedConsesusEth | Check the amount of consensus ETH in USD* |  |  |
| eth1sync | Check the sync status of Eth1 client (with Rocketpool Golang library) |  |  |
| eth2sync | Check the sync status of Eth2 client (with Rocketpool Golang library) |  |  |
| ethPrice | Check ETH/USD* price |  |  |
| gasPrice | Get the latest gas price |  |  |
| kraken | Get the latest ticker price from Kraken | ticker (string), amount (float64) | XETHZEUR, 1 |
| minStake | Check the minimum RPL stake |  |  |
| oracleRplPrice | Check the RPL price from Rocketpool oracle |  |  |
| ownEthDeposit | Check the amount of ETH deposited |  |  |
| rplFiat | Check the amount of RPL in USD* |  |  |
| rplPrice | Check RPL/USD* price |  |  |
| stakeReserve | Check the reserve of RPL stake |  |  |
| totalFunds | Check the total amount of funds in USD* |  |  |

&ast; you can use differnt fiat as quote currency in these plugins if you set `fiat` in config.yml

You can list existing plugins with `rocketreport plugins`. Plugin config file is passed with `-p` parameter to rocketreport. Plugin config file is yaml containing key `plugins` with list of selected plugins. See example in [plugins.yml](plugins.yml).

## Configuration

Some plugins need configuration, for example to get rocketpool minimum stake in RPL, you must set eth1 and eth2 clients and set your Rocketpool node address. Configuration file is in yaml format, passed with `-c` parameter to rocketreport. You can see example in [config_example.yml](config_example.yml).