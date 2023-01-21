package config

import (
	"github.com/spf13/viper"
	"time"
)

const (
	SelectedServerSlug = "selected_server_slug"
	SelectedClientSlug = "selected_client_slugs"
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

func SetSelectedClient(clientSlug string) error {
	viper.Set(SelectedClientSlug, clientSlug)
	return viper.WriteConfig()
}

func GetSelectedClient() string {
	return GetSelectedClientForServer(GetSelectedServer())
}

func GetSelectedClientForServer(serverSlug string) string {
	slugs := viper.GetStringMapString(SelectedClientSlug)

	return slugs[serverSlug]
}

func SetDefaults() {
	SetDefaultSessionDuration()
}

const sessionDuration = "session_duration"

const defaultSessionDuration = time.Minute * 10

func GetSessionDuration() time.Duration {
	return viper.GetDuration(sessionDuration)
}

func SetDefaultSessionDuration() {
	viper.SetDefault(sessionDuration, defaultSessionDuration)
}
