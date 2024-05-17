# Rocketreport

Rocketreport is a tool that fetches stats about Rocketpool node and other crypto data. It can then send a report over Telegram.

The motivation was to get regular updates about Rocketpool node to phone. Telegram seemed least bad. I'm open to implement other channels too.

## Install binaries

To install latest release for Linux:

```sh
wget -O /tmp/rocketreport https://github.com/t0mk/rocketreport/releases/latest/download/rocketreport-amd64 && chmod +x /tmp/rocketreport && sudo cp /tmp/rocketreport /usr/local/bin/
```

You can use rocketreport from Docker, the image is `t0mk/rocketreport`.

## Usage

### Print ETH gas price
```
./rocketreport plugin gasPriceBeaconcha.in
```

.. or with Docker container
```
docker run --rm t0mk/rocketreport plugin gasPriceBeaconcha.in
```

### Print output based on plugin configuration
Plugin configuration looks like [_examples/basic/plugins.yml](_examples/basic/plugins.yml)
```
./rocketreport -p _examples/basic/plugins.yml print
```

.. or with Docker
```
docker run --rm -v $(pwd)/_examples/basic:/conf t0mk/rocketreport -p /conf/plugins.yml print
```

That will print
```
Gas price is              	11.79
ETH-USDT                  	3,035 $T
Binance RPLUSDT             19.24
```

### Print Rocketpool node stats output

If you want to output Rocketpool node statistics, you need to [configure](#configuration) rocketreport.

For this to work, you need to have eth1 RPC available at http://127.0.0.1:8545 and eth2 client at http://127.0.0.1:5052. See how to do it in section [Eth1 and Eth2 client port tunnelling](#eth1-and-eth2-client-port-tunnelling).

For this example we have [_examples/rocketpool/config.yml](_examples/rocketpool/config.yml) and [_examples/rocketpool/plugins.yml](_examples/rocketpool/plugins.yml)
```
./rocketreport -c _examples/rocketpool/config.yml -p _examples/rocketpool/plugins.yml print
```

.. or with Docker
```
docker run --network host --rm -v $(pwd)/_examples/rocketpool:/conf t0mk/rocketreport -c /conf/config.yml -p /conf/plugins.yml print
```

Output might look like [this](https://i.ibb.co/17ZmXFF/ouput.png).

You need to use `--network host` for the Docker container to reach the SSH tunnels.

If this example doesn't work, try to change the `consensus_client` in [_examples/rocketpool/config.yml](_examples/rocketpool/config.yml). 

### Print portfolio details

[_examples/portfolio/plugins.yml](_examples/portfolio/plugins.yml) implements following scenario:

You have 0.5 BTC somewhere, and some ETH in address 0xC450c0F2d99c0eAFC3b53336Ac65b7f94f846478. You want to know (be regularly reminded) how much is it alltogether in USDT.

```
docker run --rm -v $(pwd)/_examples/portfolio:/conf t0mk/rocketreport -p /conf/plugins.yml print
```

Output might be:
```
My 0.5 BTC worth in USDT 	31,124
Eth in my address         	1.1176 Ξ
My eth in USDT            	3,357
My total portfolio in USDT  34,482
```


### Report over Telegram
```
./rocketreport -p plugins.yml -c config.yml send -s
```

You can run it with Docker too

```
docker run --rm t0mk/rocketreport plugin gasPrice
docker run --rm -v $(pwd):/confs/ t0mk/rocketreport -p /confs/plugins.yml -c /confs/config.yml print
```

## Plugins ([list](PLUGINS.md))

Rocketreport messages are compiled from plugin outputs, one plugin per row. That way you can configure what info you want to see in your reports.

Plugin configuration is a yaml file with list of plugins to evaluate, see `plugins.yml` in subdirs or [_examples](_examples), for example [_examples/portfolio/plugins.yml](_examples/portfolio/plugins.yml).

Plugin configuration has parameters

- `name` of the plugin. See implemented plugins in  [PLUGINS.md](PLUGINS.md), or `docker run --rm t0mk/rocketreport list-plugins`.
- `desc` is description text for your report. There's always some default but you might want to set this and change the description.
- `args` some plugins take arguments, for example `addressBalance` plugin needs an ethereum address. `args` is a list
- `labl` is a label for this plugin call. You can use it later with "meta" plugins, for adding, multuplying etc.
- `hide` is a flag that will hide this plugin call from your report.

You can see all plugins listed in [PLUGINS.md](PLUGINS.md).

Once you configure plugins, you can evaluate and print them to console with the `print` command:

```
docker run --rm -v $(pwd)/_examples/portfolio:/conf t0mk/rocketreport -p /conf/plugins.yml print
```

You can also run a single plugin:

```
docker run --rm t0mk/rocketreport plugin gasPriceBeaconcha.in
```

## Configuration

Some plugins need configuration, for example to use plugin that gets Rocketpool minimum stake in RPL, you must set URLs to eth1 and eth2 clients and set a Rocketpool node address. Configuration file is in yaml format, passed with `-c` parameter to rocketreport. You can see example in [_examples/rocketpool/config.yml](_examples/rocketpool/config.yml).

You can also configure from environment variables, envvar names are the same as in config yml but capitalized. In other words, you can use `TELEGRAM_TOKEN` envvar instead of field `telegram_token` in `config.yml`.

Configuration is "lazy". You only need to set config options which your selected plugins need. You can find out experimentally. The panic messages are confusing, but focus on the top and you should see what the problem is.

## Eth1 and Eth2 client port tunnelling

For most of the Rocketpool plugins, you need to have eth1 and eth2 client RPC API available. You pass those in `config.yml` as `eth1_url` anbd `eth2_url`. If you run Rocketpool node in managed mode, you have both running in the server, in containers `rocketpool_eth1` and `rocketpool_eth2`. If you run rocketreport from the Rocketpool node server, you can either find out IP addresses of the containers, or run rocketreport from docker (`t0mk/rocketreport`) and use container names in the config.

However, if you run rocketreport from elsewhere (like your Linux desktop), you can tunnel ports from the Rocketpool nodes to `https://127.0.0.1:8545` (eth1) and `http://127.0.0.1:5052` (eth2), and then set the localhost urls in `config.yml`. I prepared script that connects to Rocketpool server , finds IPs of the eth1 and eth2 containers, and tunnels the ports to localhost: [scripts/create_tunnels_for_eth_client.sh](scripts/create_tunnels_for_eth_client.sh).

This is further complicated if you run rocketreport from Docker. Then, if you have eth1 and eth2 clients forwarded to localhost (like suggested in previous paragraph), you need to have the localhost URLs in `eth{1,2}_url` in `config.yml`, and you need to use `docker run --network host`, so that the container can reach the tunnels. Or, if you run `t0mk/rocketreport` container on the Rocketpool node, you can use container names in config (in other words `eth1_url: http://rocketpool_eth1:8545` and `eth2_url: http://rocketpool_eth2:5052`). Then, you need to run the container in the `rocketpool_net` network - `docker run --rm --network rocketpool_net`.

## Build

If you want to change code and test, just do `make build`. Interoperable static build that works on various Linux distros is a bit more complicated because there's some C crypto in dependencies (I use code from RP smartnode). It's in makefile, so just do `make static-build-amd64`. There's still some issues with static build for arm64.

## Telegram

You can send reports over Telegrem. You can have either
- rocketreport listening on Telegram bot message, aka "serve" mode
- rocketreport just sending message from bot to a chat - `rocketreport send`

### Create bot and get token

To use Telegram bot with Rocketreport, you need to create your bot first. Follow https://core.telegram.org/bots/tutorial until "Obtain Your Bot Token", and then use the token in config.yml as `telegram_token`.

### Find out Chat ID

Rocketreport bot can only send message to a single chat (for the sake of security). Chat is a Telegram chat in your phone app. Once you get your bot token, you need to send message to the bot and find out ID of your chat. Just run

```
TELEGRAM_TOKEN=... rocketreport report-chat-id
```
.. or with Docker

```
docker run -e TELEGRAM_TOKEN="..." --rm t0mk/rocketreport report-chat-id
```

.. and then send a message to the bot. You can open Telegram chat with the bot by visiting link `https://t.me/<bot_username>`. The bot will reply with chat ID and will also print it to stdout.

### "serve" mode

Once you have Telegram token and a chat ID, you can send messages via Telegram.

You can run rocketreport in "serve" mode, where the rocketreport process stays on and reacts on Telegram messages. You can then get reports in your Telegram chat on-demand, by pressing a button.

#### Message schedule

You can also have reports sent regularly by setting TELEGRAM_MESSAGE_SCHEDULE to a 
[Cron expression](https://en.wikipedia.org/wiki/Cron#Overview), for example, with  `plugins.yml` as

```yaml
plugins:
  - name: gasPriceBeaconcha.in
    desc: Gas price is
```

and `config.yml` as

```yaml
telegram_token: ...
telegram_chat_id: ...
```

You will get report with Gas price to Telegram every minute by running:
```
TELEGRAM_MESSAGE_SCHEDULE="* * * * *" rocketreport -c config.yml -p plugins.yml serve
```

### Single send

If you'd like to get reports regularly, but don't want to have rocketreport running all the time (in the "serve" mode), you can also send a single message to a Telegram chat. Considering you put `telegram_token` and `telegram_chat_id` to config.yml, and some plugin configuration to `plugins.yml`, you can run

```
rocketreport -c config.yml -p plugins.yml send -s
```

Rocketreport will evaluate the plugins and send report to Telegram.

You can put this command to Cron if you want to have the report sent regularly.

### Telegram message "header"

The report in Telegram message is in format of "inline keyboard", so that it looks like a table. The message body is meant to be very short, but meaningful because you will see it in notification:

![notification](https://iamges.com/rr_notification.png)

You can template the message body to a small extent by the TELEGRAM_HEADER_TEMPLATE configuration value. It's a string with space-delimited identifiers. If identifier starts with "_", it's substituted by looked-up value. Otherwise, identifier is put to the message.

The substitution identifiers can refer to configured plugins by label ("labl" find in plugin conf), or they refer to plugins with no arguments.

Examples of header templates:
- `TELEGRAM_HEADER_TEMPLATE="Hi!"` will send report with "Hi!" in the body
- `TELEGRAM_HEADER_TEMPLATE="_timeMin ETH: _ethPrice` will send report with "2024-05-16_14:34 ETH: 2,998 $T" in the header.
- If you use [plugins.yml from _examples/portfolio](_examples/portfolio/plugins.yml), you can do `TELEGRAM_HEADER_TEMPLATE="Total: _total"`, and the header will have value of plugin marked with "labl" "total".
