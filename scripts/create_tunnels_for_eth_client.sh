#!/bin/bash
#set -x

# This script opens SSH tunnel to eth1 and eth2 JSON RPC from your RP node
# so that you can use it locally with rocketreport.
# It will tunnel eth1:8545 to 127.0.0.1:8545 and eth2:5052 to 127.0.0.1:5052.
# You can use those as eth1_url and eth2_url in rocketreport config.

# To use this, you need to set user, port and identity for your RP node
# in your .ssh/config, and change pass the hostname as argument to this script.


NODE=$1

[ -z "$NODE" ] && echo "Usage: $0 <sshconfig_hostname>" && exit 1


ETH1PORT=8545
ETH2PORT=5052
ETH1NAME=rocketpool_eth1
ETH2NAME=rocketpool_eth2

ETH1IP=`ssh $NODE "docker ps -q -f name=$ETH1NAME | xargs docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'"`
ETH2IP=`ssh $NODE "docker ps -q -f name=$ETH2NAME | xargs docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'"`

ssh -NTfL ${ETH1PORT}:${ETH1IP}:${ETH1PORT} $NODE
ssh -NTfL ${ETH2PORT}:${ETH2IP}:${ETH2PORT} $NODE

echo "Eth1 RPC available at:"
echo "http://127.0.0.1:${ETH1PORT}"
echo
echo "Eth2 RPC available at:"
echo "http://127.0.0.1:${ETH2PORT}"

