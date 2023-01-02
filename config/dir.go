package config

import "github.com/kirsle/configdir"

func ConfigDir() string {
	localConfig := configdir.LocalConfig("oauth-commander")

	err := configdir.MakePath(localConfig)
	if err != nil {
		panic(err)
	}

	return localConfig
}

func ServerDir() string {
	severConfig := configdir.LocalConfig("oauth-commander", "servers")

	err := configdir.MakePath(severConfig)
	if err != nil {
		panic(err)
	}

	return severConfig
}
