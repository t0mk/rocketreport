package plugins

import "github.com/t0mk/rocketreport/utils"

type Plugin interface {
	Desc() string
	Value () string
	Help() string
	Exec() error
}

type Eth1Sync struct {
	Plugin
	value *string
}

func (e *Eth1Sync) Desc() string {
	return "Eth1 client"
}

func (e *Eth1Sync) Value() string {
	if e.value == nil {
		*e.value = utils.EthClientStatusString(ecs.PrimaryClientStatus)
	}
	return *e.value
}

func (e *Eth1Sync) Help() string {
	return "Check the sync status of the Eth1 client"
}

