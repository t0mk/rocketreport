# Rocketreport

Rocketreport is a tool that fetches stats about Rocketpool node and other crypto data. It can then send it over Telegram.

The motivation was to get regular updates about Rocketpool node to phone. Telegram seemed least bad. I'm open to implement other channels too.

## Install binaries

To install latest release for Linux:

```sh
wget -O /tmp/rocketreport https://github.com/t0mk/rocketreport/releases/latest/download/rocketreport-amd64 && chmod +x /tmp/rocketreport && sudo cp /tmp/rocketreport /usr/local/bin/
```

You can use rocketreport from Docker, the image is `t0mk/rocketreport`.

## Usage

### Display ETH gas price
```
./rocketreport plugin gasPriceBeaconcha.in
```

.. or with Docker container
```
docker run --rm t0mk/rocketereport plugin gasPriceBeaconcha.in
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

### Print Rocketpool stats output

If you want to output Rocketpool status, you need to [configure](#configuration) rocketreport.

For this to work, you need to have eth1 RPC available at http://127.0.0.1:8545 and eth2 client at http://127.0.0.1:5052. See how to do it in section [Eth1 and Eth2 client port tunnelling](#eth1-and-eth2-client-port-tunnelling).

For this example we have [_examples/rocketpool/config.yml](_examples/rocketpool/config.yml) and [_examples/rocketpool/plugins.yml](_examples/rocketpool/plugins.yml)
```
./rocketreport -c _examples/rocketpool/config.yml -p _examples/rocketpool/plugins.yml print
```

.. or with Docker
```
docker run --network host --rm -v $(pwd)/_examples/rocketpool:/conf t0mk/rocketreport -c /conf/config.yml -p /conf/plugins.yml print
```

You need to use `--network host` for the Docker container to reach the SSH tunnels.

If this example doesn't work, try to change the `consensus_client` in [_examples/rocketpool/config.yml](_examples/rocketpool/config.yml). 

### Report portfolio value in USDT

[_examples/portfolio/plugins.yml](_examples/portfolio/plugins.yml) implements following scenario:

You have 0.5 BTC, and some ETH in address 0xC450c0F2d99c0eAFC3b53336Ac65b7f94f846478. You want to know (be regularly reminded) how much is it USDT.

```
docker run --rm -v $(pwd)/_examples/portfolio:/conf t0mk/rocketreport -p /conf/plugins.yml print
```

Output might be:
```
Binance ticker BTCUSDT 0.5 	31,124
Eth in my address         	1.1176 Ξ
My eth in USDT            	3,357
My portfolio in USDT      	34,482
```


### Send same stats as Telegram message
```
./rocketreport -p plugins.yml -c config.yml send -s
```

You can run it with Docker too

```
docker run --rm t0mk/rocketreport plugin gasPrice
docker run --rm -v $(pwd):/confs/ t0mk/rocketreport -p /confs/plugins.yml -c /confs/config.yml print
```

## Plugins

Rocketreport messages are compiled from plugin output, one plugin per row. That way you can configure what info you want to see in your messages.

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
docker run --rm t0mk/rocketereport plugin gasPriceBeaconcha.in
```

## Configuration

Some plugins need configuration, for example to use plugin that gets Rocketpool minimum stake in RPL, you must set URLs to eth1 and eth2 clients and set a Rocketpool node address. Configuration file is in yaml format, passed with `-c` parameter to rocketreport. You can see example in [_examples/rocketpool/config.yml](_examples/rocketpool/config.yml).

You can also configure from environment variables, envvar names are the same as in config yml but capitalized. In other words, you can use `TELEGRAM_TOKEN` envvar instead of field `telegram_token` in `config.yml`.

Configuration is "lazy". You only need to set config options which your selected plugins need. You can find out experimentally. The panic messages are confusing, but focus on the top and you should see what the problem is.

## Eth1 and Eth2 client port tunnelling

For most of the Rocketpool plugins, you need to have eth1 and eth2 client RPC API available. You pass those in `config.yml` as `eth1_url` anbd `eth2_url`. If you run Rocketpool node in managed mode, you have both running in the server, in containers `rocketpool_eth1` and `rocketpool_eth2`. If you run rocketreport from the Rocketpool node server, you can either find out IP addresses of the containers, or run rocketreport from docker (`t0mk/rocketreport`) and use container names in the config.

However, if you run rocketreport from elsewhere (like your Linux desktop), you can tunnel ports from the Rocketpool nodes to `https://127.0.0.1:8545` (eth1) and `http://127.0.0.1:5052` (eth2), and then set the localhost urls in `config.yml`. I prepared script that connects to Rocketpool server , finds IPs of the eth1 and eth2 containers, and tunnels the ports to localhost: [scripts/create_tunnels_for_eth_client.sh](scripts/create_tunnels_for_eth_client.sh).

This is further complicated if you run rocketreport from Docker. Then, if you have eth1 and eth2 clients forwarded to localhost, you need to have the localhost URLs in `eth{1,2}_url` and you need to use `docker run --network host`, so that the container can reach the tunnels. Or, if you run `t0mk/rocketreport` container on the Rocketpool node, you can use container names in config (in other words `http://rocketpool_eth{1,2}`), but you need to run the container in the `rocketpool_net` network - run it like `docker run --rm --network rocketpool_net`.

## Build

If you want to change code and test, just do `make build`. Interoperable static build that works on various Linux distros is a bit more complicated because there's some C crypto in dependencies (or RP smartnode that I'm using). It's in makefile, so just do `make static-build.

## Telegram bot

To serve Telegram bot with Rocketreport, you need to create your bot first. Follow https://core.telegram.org/bots/tutorial until "Obtain Your Bot Token", and then use the token in config.yml as `telegram_token`.

## Telegram chat send

If you'd like to get reports regularly, but don't want to haver Telegram bot running all the time (in the "serve" mode), you can configure it to send message to a Telegram Chat, maybe to your mobile device. To do that, you need to specify Telegram Chat ID. Frist create the bot, put `telegram_token` to `config.yml` and run

```
rocketreport -c config.yml report-chat-id
```

rocketreport will wait for a messagei from your device, and then print Chat ID which you can put to `config.yml` as `telegram_chat_id`.

Once you have both bot token and Chat ID set, you can create a cronjob to trigger

```
rocketreport -c config.yml -p plugins.yml send -s
```
