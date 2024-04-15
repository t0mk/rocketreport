# Rocketreport

Rocketreport is a tool that fetches stats about Rocketpool node and other crypto data and sends it to you over Telegram. 
- It can be configured in cron-like fashion and send messages to existing chat
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

| Name | Description | Args | Defaults |
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
| minStake | Check the minimum RPL stake for Rocketpool node |  |  |
| oracleRplPrice | Check the RPL price from Rocketpool oracle |  |  |
| ownEthDeposit | Check the amount of ETH deposited in Rocketpool node |  |  |
| prod | Product of given args, either numbers or plugin outputs, multiplies args and outputs a float | list of values - numbers or plugin outputs ([]interface {}) | [] |
| rplFiat | Check the amount of RPL in USD*, deposited to Rocketpool node |  |  |
| rplPrice | Check RPL/USD* price (RPL/ETH based on Rocketpool Oracle) |  |  |
| stakeReserve | Check the reserve of RPL stake in Rocketpool node (above 10%) |  |  |
| sum | Sum of given args, either numbers or plugin outputs, adds args and outputs a float | list of values - numbers or plugin outputs ([]interface {}) | [] |
| totalFunds | Check the total amount of funds in USD* |  |  |

&ast; you can use different fiat as quote currency in these plugins if you set `fiat` option in config.yml

You can list existing plugins with `rocketreport plugins`. Plugin config file is passed with `-p` parameter to rocketreport. Plugin config file is yaml containing key `plugins` with list of selected plugins. See example in [plugins.yml](plugins.yml).

Once you configure plugins, you can evaluate and print them to console with

```
rocketreport -c config.yml -p plugins.yml print
```

You can also run a single plugin:

```
rocketreport plugin 
```

## Configuration

Some plugins need configuration, for example to get rocketpool minimum stake in RPL, you must set URLs to eth1 and eth2 clients and set your Rocketpool node address. Configuration file is in yaml format, passed with `-c` parameter to rocketreport. You can see example in [config_example.yml](config_example.yml).

You can also configure from environment variables, envvar names are the same as in config yml but capitalized. In other words, you can use `TELEGRAM_TOKEN` envvar instead of field `telegram_token` in `config.yml`.

Configuration is "lazy". You only need to set config options which your selected plugins need. You can find out experimentally.

## Telegram bot

To serve Telegram bot with Rocketreport, you need to create your bot first. Follow https://core.telegram.org/bots/tutorial until "Obtain Your Bot Token", and then use the token in config.yml as `telegram_token`.

## Telegram chat send

If you'd like to get reports regularly, but don't want to have rocketreport running all the time (in the "serve" mode), you can configure it to send message to a Telegram Chat, maybe on your mobile device. To do that, create the bot, put `telegram_token` to `config.yml` and run

```
rocketreport -c config.yml report-chat-id
```

rocketreport will wait for a message and then print Chat ID which you can put to `config.yml` as `telegram_chat_id`.

Once you have both telegratm config options set, you can create a cronjob to trigger

```
rocketreport -c config.yml -p plugins.yml send -s
```
