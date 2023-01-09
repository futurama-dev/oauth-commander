package config

import (
	"github.com/spf13/viper"
)

const (
	SelectedServerSlug  = "selected_server_slug"
	SelectedClientSlugs = "selected_client_slugs"
)

func GetSelectedServer() string {
	return viper.GetString(SelectedServerSlug)
}

func SetSelectedServer(serverSlug string) error {
	viper.Set(SelectedServerSlug, serverSlug)
	return viper.WriteConfig()
}

func IsSelectedServer() bool {
	return len(GetSelectedServer()) > 0
}

func GetSelectedClient() string {
	return GetSelectedClientForServer(GetSelectedServer())
}

func GetSelectedClientForServer(serverSlug string) string {
	slugs := viper.GetStringMapString(SelectedClientSlugs)

	return slugs[serverSlug]
}

func SetDefaults() {
	SetDefaultSessionDuration()
}
