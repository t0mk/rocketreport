package plugins

import "github.com/t0mk/rocketreport/utils"

func RegisterExtraPlugins() {
	Plugins = append(Plugins, []Plugin{
		{
			Key:  "smoothingPoolBalance",
			Desc: "Smoothing Pool Balance",
			Help: "ETH in the smoothing pool",
			Formatter: FloatSuffixFormatter(2, "ETH"),
			Refresh: func() (interface{}, error) {
				b, err := utils.SmoothingPoolBalance()
				if err != nil {
					return nil, err
				}
				return *b, nil
			},
		},
	}...)

}
