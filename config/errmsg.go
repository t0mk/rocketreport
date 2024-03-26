package config

const (
	errEc = `Couldn't find execution client url. You can:
	       - set ETH1_URL envvar, for example to http://localhost:8545
		   - run docker container called "eth1". This tool will try to find it and use it
		   You can for example tunnel eth1 JSON RPC from rocketpool node as 
		   $ ssh -NTfL 8545:172.19.0.10:8545 rocketpoolnode, where 172.19.0.10 is local private IP of eth1 container`

)
