package utils

import (
	"github.com/rocket-pool/smartnode/shared/types/api"
	"fmt"
)



func EthClientStatusString(status api.ClientStatus) string {
	sentence := ""
	if !status.IsWorking {
		sentence += "not working,"
	}
	if status.IsSynced {
		sentence += "synced"
	} else {
		sentence += "not synced,"
	}
	if status.SyncProgress < 1 {
		sentence += fmt.Sprintf(" syncing, now at %d%%", int(100*(status.SyncProgress)))
	}
	if status.Error != "" {
		sentence += fmt.Sprintf(", Error: %s", status.Error)
	}
	return sentence
}

