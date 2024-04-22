# Rocketreport

Rocketreport is a tool that fetches stats about Rocketpool node and other crypto data and sends it to you over Telegram. 
- It can be configured in cron-like fashion and send messages to existing chat
- It can serve a Telegram bot and reply on demand

## Basic usage

Use `gasPrice` plugin to display Eth gas price:
```
./rocketreport plugin gasPriceBeaconChain
```

Print stats based on plugin configiuration, pass config file
```
./rocketreport -p plugins.yml -c config.yml print
```

Send same stats as Telegram message
```
./rocketreport -p plugins.yml -c config.yml send -s
```

You can run it with Docker too

```
docker run --rm t0mk/rocketreport plugin gasPrice
docker run --rm -v $(pwd):/confs/ t0mk/rocketreport -p /confs/plugins.yml -c /confs/config.yml print
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

Plugin configuration is a yaml file with list of plugins to evaluate, for example:

```yaml
plugins:
  - name: rpActualStake
    id: stake
  - name: rpEth1sync
  - name: rpEth2sync
  - name: rpIntervalEnds
  - name: rpMinStake
    id: min
  - name: rpOracleRplPrice
  - name: rplPrice
    id: rplusd
  - name: sub
    id: reserve
    desc: RPL reserve
    args:
      - stake
      - min
  - name: mul
    desc: RPL reserve in USD
    args:
        - reserve
        - rplusd  
```

- `name` is plugins name (look at `rocketreport list-plugins`)
- `id` is arbitrary ID you pick to refer to output of a plugin
- `desc` is description you want to see in report message (if your prefer to change it over the default)
- `args` is a list of arguments to a plugin. You can see plugin description is `list-plugins` command.

You can see all plugins listed in [PLUGINS.md](PLUGINS.md).

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

Some plugins need configuration, for example to get rocketpool minimum stake in RPL, you must set URLs to eth1 and eth2 clients and set a Rocketpool node address. Configuration file is in yaml format, passed with `-c` parameter to rocketreport. You can see example in [config_example.yml](config_example.yml).

You can also configure from environment variables, envvar names are the same as in config yml but capitalized. In other words, you can use `TELEGRAM_TOKEN` envvar instead of field `telegram_token` in `config.yml`.

Configuration is "lazy". You only need to set config options which your selected plugins need. You can find out experimentally.

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

