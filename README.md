# Rocketreport

Rocketreport is a tool that fetches stats about Rocketpool node and other crypto __ and sends it to Telegram. It can 
send it to existing Chat or it can serve a Telegram bot and reply on demand.

## Basic usage

You can run it with Docker. To print current Ethereum gas price, do:

```
docker run --rm t0mk/rocketreport gasPrice
```

## Install binaries

In can you don't want to use Docker, you can get release artifacts from GitHub.

To install latest release for Linux:

```sh
wget -O /tmp/rocketreport https://github.com/t0mk/rocketreport/releases/latest/download/rocketreport-linux-amd64 && chmod +x /tmp/rocketreport && sudo cp /tmp/rocketreport /usr/local/bin/
```

.. for MacOS:

```sh
wget -O /tmp/rocketreport https://github.com/t0mk/rocketreport/releases/latest/download/rocketreport-darwin-amd64 && chmod +x /tmp/rocketreport && sudo cp /tmp/rocketreport /usr/local/bin/
```

## Plugins

Rocketreport 

You can list existing plugins with `rocketreport plugins`. Plugin config file is passed with `-p` parameter to rocketreport. Plugin config file is yaml containing key `plugins` with list of selected plugins. See example in [plugins.yml](plugins.yml).

## Configuration

Some plugins need configuration, for example to get rocketpool minimum stake in RPL, you must set eth1 and eth2 clients and set your node address. Configuration file is also in yaml format, passed with `-c` parameter to rocketreport. You can see example in [config_example.yml](config_example.yml).