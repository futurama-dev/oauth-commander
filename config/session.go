package config

import (
	"errors"
	"github.com/spf13/viper"
	"net/url"
	"time"
)

const (
	authorizationSessions  = "auth_sessions"
	sessionDuration        = "session_duration"
	defaultSessionDuration = time.Minute * 10
)

func GetSessionDuration() time.Duration {
	return viper.GetDuration(sessionDuration)
}

func SetDefaultSessionDuration() {
	viper.SetDefault(sessionDuration, defaultSessionDuration)
}

type Session struct {
	State      string
	AuthReqUrl string    `yaml:"auth_req_url" mapstructure:"auth_req_url"`
	ServerSlug string    `yaml:"server_slug" mapstructure:"server_slug"`
	ClientSlug string    `yaml:"client_slug" mapstructure:"client_slug"`
	CreatedAt  time.Time `yaml:"created_at" mapstructure:"created_at"`
	ExpiresAt  time.Time `yaml:"expires_at" mapstructure:"expires_at"`
}

func NewSession(state string, authReqUrl url.URL, serverSlug, clientSlug string) Session {
	now := time.Now().Truncate(time.Second)

	return Session{
		State:      state,
		AuthReqUrl: authReqUrl.String(),
		ServerSlug: serverSlug,
		ClientSlug: clientSlug,
		CreatedAt:  now,
		ExpiresAt:  now.Add(GetSessionDuration()),
	}
}

func (s Session) Save() error {
	sessions, err := GetAuthorizationSessions()
	if err != nil {
		return err
	}

	if _, found := sessions.FindByState(s.State); found {
		return errors.New("a session with this state already exists: " + s.State)
	}

	sessions = append(sessions, s)

	return SetAuthorizationSessions(sessions)
}

func (s Session) IsExpired() bool {
	now := time.Now().Truncate(time.Second)

	//return s.ExpiresAt.Before(now) || s.ExpiresAt.Equal(now)
	return !s.ExpiresAt.After(now)
}

type Sessions []Session

func (ss Sessions) FindByState(state string) (Session, bool) {
	for _, session := range ss {
		if session.State == state {
			return session, true
		}
	}

	return Session{}, false
}

func GetAuthorizationSessions() (Sessions, error) {
	var sessions Sessions
	err := viper.UnmarshalKey(authorizationSessions, &sessions)
	if err != nil {
		return Sessions{}, err
	}

	return sessions, nil
}

func SetAuthorizationSessions(sessions Sessions) error {
	viper.Set(authorizationSessions, sessions)
	return viper.WriteConfig()
}
