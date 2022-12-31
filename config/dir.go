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
