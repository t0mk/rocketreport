package common

import (
	"time"

	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
)

func now(...interface{}) (interface{}, error) {
	return time.Now(), nil
}

func TimePlugins() map[string]types.RRPlugin {
	return map[string]types.RRPlugin{
		"timeSec": {
			Cat:       types.PluginCatCommon,
			Desc:      "Time",
			Help:      "current time up to seconds",
			Formatter: formatting.Time("2006-01-02_15:04:05"),
			Refresh:   now,
		},
		"timeMin": {
			Cat:       types.PluginCatCommon,
			Desc:      "Time",
			Help:      "current time up to minutes",
			Formatter: formatting.Time("2006-01-02_15:04"),
			Refresh:   now,
		},
		"date": {
			Cat:       types.PluginCatCommon,
			Desc:      "Date",
			Help:      "current date",
			Formatter: formatting.Time("2006-01-02"),
			Refresh:   now,
		},
	}
}
