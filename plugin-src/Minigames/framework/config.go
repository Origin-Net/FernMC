package framework

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)


type Config struct {
	GameSettings map[string]GameConfig `yaml:"games"`
	Global       GlobalConfig          `yaml:"global"`
	Database     DatabaseConfig        `yaml:"database"`
}


type GameConfig struct {
	Enabled              bool `yaml:"enabled"`
	MinPlayers           int  `yaml:"min_players"`
	MaxPlayers           int  `yaml:"max_players"`
	CountdownSeconds     int  `yaml:"countdown_seconds"`
	GameDurationSeconds  int  `yaml:"game_duration_seconds"`
	GracePeriodSeconds   int  `yaml:"grace_period_seconds"`
	InvincibilitySeconds int  `yaml:"invincibility_seconds"`
	BackendTaskTicks     int  `yaml:"backend_task_ticks"`
}


type GlobalConfig struct {
	PluginPrefix    string `yaml:"plugin_prefix"`
	JoinCommand     string `yaml:"join_command"`
	LeaveCommand    string `yaml:"leave_command"`
	Language        string `yaml:"language"`
	EnablePvP       bool   `yaml:"enable_pvp"`
	DataSaveInterval int   `yaml:"data_save_interval_seconds"`
}


type DatabaseConfig struct {
	Type string `yaml:"type"` 
	Path string `yaml:"path"`
}


type ConfigManager struct {
	dataDir string
	config  *Config
	logger  *LogWrapper
}


func NewConfigManager(dataDir string, logger *LogWrapper) (*ConfigManager, error) {
	cm := &ConfigManager{
		dataDir: dataDir,
		logger:  logger,
	}

	if err := cm.load(); err != nil {
		return nil, err
	}

	return cm, nil
}

func (cm *ConfigManager) load() error {
	configPath := filepath.Join(cm.dataDir, "config", "config.yml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cm.logger.Warn("config.yml not found, creating default")
			if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
				return fmt.Errorf("create config dir: %w", err)
			}
			cm.config = DefaultConfig()
			return cm.save()
		}
		return fmt.Errorf("read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	cm.config = &config
	return nil
}

func (cm *ConfigManager) save() error {
	configPath := filepath.Join(cm.dataDir, "config", "config.yml")
	data, err := yaml.Marshal(cm.config)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}


func (cm *ConfigManager) Config() *Config {
	return cm.config
}


func (cm *ConfigManager) GameConfig(gameID string) GameConfig {
	if gc, ok := cm.config.GameSettings[gameID]; ok {
		return gc
	}
	return GameConfig{
		Enabled:              true,
		MinPlayers:           2,
		MaxPlayers:           8,
		CountdownSeconds:     30,
		GameDurationSeconds:  600,
		GracePeriodSeconds:   30,
		InvincibilitySeconds: 5,
		BackendTaskTicks:     10,
	}
}


func DefaultConfig() *Config {
	return &Config{
		GameSettings: map[string]GameConfig{
			"skywars": {
				Enabled:              true,
				MinPlayers:           2,
				MaxPlayers:           8,
				CountdownSeconds:     30,
				GameDurationSeconds:  600,
				GracePeriodSeconds:   30,
				InvincibilitySeconds: 5,
				BackendTaskTicks:     10,
			},
		},
		Global: GlobalConfig{
			PluginPrefix:     "§6[Minigames]§r",
			JoinCommand:      "join",
			LeaveCommand:     "leave",
			Language:         "en",
			EnablePvP:        true,
			DataSaveInterval: 300,
		},
		Database: DatabaseConfig{
			Type: "json",
			Path: "stats",
		},
	}
}
