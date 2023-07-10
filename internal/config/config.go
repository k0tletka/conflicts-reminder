package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	configLocation = "config.toml"
)

type Config struct {
	RemindCycleMinutes int          `mapstructure:"remind_cycle_minutes"`
	Slack              SlackConfig  `mapstructure:"slack"`
	Binds              []UserBind   `mapstructure:"user_bind"`
	Gitlab             GitlabConfig `mapstructure:"gitlab"`

	bindMapGitlabToSlack map[int]string
}

func (c *Config) fillMaps() {
	for _, bindCfg := range c.Binds {
		c.bindMapGitlabToSlack[bindCfg.GitlabID] = bindCfg.SlackID
	}
}

func (c *Config) GetSlackIDByGitlabID(gitlabId int) string {
	slackId, _ := c.bindMapGitlabToSlack[gitlabId]
	return slackId
}

func (c *Config) CheckGitlabId(gitlabId int) bool {
	_, ok := c.bindMapGitlabToSlack[gitlabId]
	return ok
}

type SlackConfig struct {
	Token                 string `mapstructure:"token"`
	NotificationChannelID string `mapstructure:"notify_channel_id"`
}

type GitlabConfig struct {
	Token         string `mapstructure:"token"`
	GitlabAddress string `mapstructure:"gitlab_address"`
	ProjectName   string `mapstructure:"project_name"`
}

type UserBind struct {
	SlackID  string `json:"slack_id"`
	GitlabID int    `json:"gitlab_id"`
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
