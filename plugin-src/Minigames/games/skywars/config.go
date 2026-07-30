package skywars


type Config struct {
	ID                   string `yaml:"id"`
	MinPlayers           int    `yaml:"min_players"`
	MaxPlayers           int    `yaml:"max_players"`
	CountdownSeconds     int    `yaml:"countdown_seconds"`
	GameDurationSeconds  int    `yaml:"game_duration_seconds"`
	GracePeriodSeconds   int    `yaml:"grace_period_seconds"`
	InvincibilitySeconds int    `yaml:"invincibility_seconds"`
	AutoRefillChests     bool   `yaml:"auto_refill_chests"`
	RefillInterval       int    `yaml:"refill_interval_seconds"`
	DropLootOnDeath      bool   `yaml:"drop_loot_on_death"`
}


func DefaultConfig() *Config {
	return &Config{
		ID:                   "skywars",
		MinPlayers:           1,
		MaxPlayers:           8,
		CountdownSeconds:     5,
		GameDurationSeconds:  600,
		GracePeriodSeconds:   30,
		InvincibilitySeconds: 5,
		AutoRefillChests:     true,
		RefillInterval:       120,
		DropLootOnDeath:      true,
	}
}
