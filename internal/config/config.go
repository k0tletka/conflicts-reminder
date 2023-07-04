package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	configLocation = "config.toml"
)

type Config struct {
	RemindCycleMinutes int         `mapstructure:"remind_cycle_minutes"`
	Slack              SlackConfig `mapstructure:"slack"`
	Binds              []UserBind  `mapstructure:"user_bind"`

	bindMapGitlabToSlack map[string]string
}

func (c *Config) fillMaps() {
	for _, bindCfg := range c.Binds {
		c.bindMapGitlabToSlack[bindCfg.GitlabID] = bindCfg.SlackID
	}
}

func (c *Config) GetSlackIDByGitlabID(gitlabId string) string {
	slackId, _ := c.bindMapGitlabToSlack[gitlabId]
	return slackId
}

type SlackConfig struct {
	Token                 string `mapstructure:"token"`
	NotificationChannelID string `json:"notify_channel_id"`
}

type UserBind struct {
	SlackID  string `json:"slack_id"`
	GitlabID string `json:"gitlab_id"`
}

func ReadConfig() *Config {
	viperReader := viper.New()
	viperReader.SetConfigType("toml")
	viperReader.SetConfigFile(configLocation)

	if err := viperReader.ReadInConfig(); err != nil {
		zap.S().With("err", err).Fatal("Can't initialize config")
	}

	var cfg *Config
	if err := viperReader.Unmarshal(cfg); err != nil {
		zap.S().With("err", err).Fatal("Can't unmarshal data to config")
	}

	cfg.fillMaps()
	return cfg
}
