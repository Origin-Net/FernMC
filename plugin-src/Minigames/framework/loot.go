package framework

import (
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)


type LootEntry struct {
	Item        string            `yaml:"item"`
	AmountMin   int               `yaml:"amount_min"`
	AmountMax   int               `yaml:"amount_max"`
	Chance      float64           `yaml:"chance"`        
	Enchantments []LootEnchantment `yaml:"enchantments,omitempty"`
}


type LootEnchantment struct {
	ID    string `yaml:"id"`
	Level int    `yaml:"level"`
}


type LootTable struct {
	Entries   []LootEntry `yaml:"entries"`
	MinRolls  int         `yaml:"min_rolls"`
	MaxRolls  int         `yaml:"max_rolls"`
}


type LootManager struct {
	tables  map[string]*LootTable
	dataDir string
	logger  *LogWrapper
}


func NewLootManager(dataDir string, logger *LogWrapper) (*LootManager, error) {
	lm := &LootManager{
		tables:  make(map[string]*LootTable),
		dataDir: filepath.Join(dataDir, "config", "loot"),
		logger:  logger,
	}

	if err := lm.loadTables(); err != nil {
		return nil, err
	}

	return lm, nil
}

func (lm *LootManager) loadTables() error {
	entries, err := os.ReadDir(lm.dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			lm.logger.Warn("Loot directory not found, creating default", "path", lm.dataDir)
			if err := os.MkdirAll(lm.dataDir, 0755); err != nil {
				return fmt.Errorf("create loot dir: %w", err)
			}
			if err := lm.writeDefaultTables(); err != nil {
				return fmt.Errorf("write default loot: %w", err)
			}
			entries, _ = os.ReadDir(lm.dataDir)
		} else {
			return fmt.Errorf("read loot dir: %w", err)
		}
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yml" {
			continue
		}
		name := entry.Name()[:len(entry.Name())-4]
		data, err := os.ReadFile(filepath.Join(lm.dataDir, entry.Name()))
		if err != nil {
			lm.logger.Error("Failed to read loot table", "file", entry.Name(), "error", err)
			continue
		}
		var table LootTable
		if err := yaml.Unmarshal(data, &table); err != nil {
			lm.logger.Error("Failed to parse loot table", "file", entry.Name(), "error", err)
			continue
		}
		lm.tables[name] = &table
		lm.logger.Info("Loaded loot table", "name", name, "entries", len(table.Entries))
	}

	return nil
}


func (lm *LootManager) GenerateLoot(tableName string) []map[string]int {
	table, ok := lm.tables[tableName]
	if !ok {
		lm.logger.Warn("Loot table not found", "name", tableName)
		return nil
	}

	rolls := table.MinRolls
	if table.MaxRolls > table.MinRolls {
		rolls += rand.IntN(table.MaxRolls - table.MinRolls + 1)
	}

	var result []map[string]int
	for i := 0; i < rolls; i++ {
		for _, entry := range table.Entries {
			if rand.Float64() > entry.Chance {
				continue
			}
			amount := entry.AmountMin
			if entry.AmountMax > entry.AmountMin {
				amount += rand.IntN(entry.AmountMax - entry.AmountMin + 1)
			}
			result = append(result, map[string]int{
				entry.Item: amount,
			})
			break
		}
	}

	return result
}

func (lm *LootManager) writeDefaultTables() error {
	defaults := map[string]*LootTable{
		"normal": {
			MinRolls: 3,
			MaxRolls: 6,
			Entries: []LootEntry{
				{Item: "minecraft:stone_sword", AmountMin: 1, AmountMax: 1, Chance: 0.3},
				{Item: "minecraft:wooden_sword", AmountMin: 1, AmountMax: 1, Chance: 0.5},
				{Item: "minecraft:leather_chestplate", AmountMin: 1, AmountMax: 1, Chance: 0.2},
				{Item: "minecraft:chainmail_boots", AmountMin: 1, AmountMax: 1, Chance: 0.15},
				{Item: "minecraft:arrow", AmountMin: 4, AmountMax: 12, Chance: 0.4},
				{Item: "minecraft:bow", AmountMin: 1, AmountMax: 1, Chance: 0.2},
				{Item: "minecraft:cooked_beef", AmountMin: 2, AmountMax: 6, Chance: 0.5},
				{Item: "minecraft:bread", AmountMin: 2, AmountMax: 4, Chance: 0.6},
				{Item: "minecraft:golden_apple", AmountMin: 1, AmountMax: 1, Chance: 0.1},
				{Item: "minecraft:iron_ingot", AmountMin: 1, AmountMax: 3, Chance: 0.25},
				{Item: "minecraft:stick", AmountMin: 2, AmountMax: 6, Chance: 0.4},
			},
		},
		"center": {
			MinRolls: 4,
			MaxRolls: 8,
			Entries: []LootEntry{
				{Item: "minecraft:iron_sword", AmountMin: 1, AmountMax: 1, Chance: 0.4},
				{Item: "minecraft:diamond_sword", AmountMin: 1, AmountMax: 1, Chance: 0.15},
				{Item: "minecraft:iron_chestplate", AmountMin: 1, AmountMax: 1, Chance: 0.3},
				{Item: "minecraft:diamond_boots", AmountMin: 1, AmountMax: 1, Chance: 0.1},
				{Item: "minecraft:golden_apple", AmountMin: 1, AmountMax: 2, Chance: 0.3},
				{Item: "minecraft:ender_pearl", AmountMin: 1, AmountMax: 2, Chance: 0.2},
				{Item: "minecraft:diamond", AmountMin: 1, AmountMax: 3, Chance: 0.2},
				{Item: "minecraft:obsidian", AmountMin: 2, AmountMax: 6, Chance: 0.15},
				{Item: "minecraft:cooked_beef", AmountMin: 4, AmountMax: 8, Chance: 0.5},
				{Item: "minecraft:bow", AmountMin: 1, AmountMax: 1, Chance: 0.3},
				{Item: "minecraft:arrow", AmountMin: 8, AmountMax: 16, Chance: 0.4},
				{Item: "minecraft:enchanted_golden_apple", AmountMin: 1, AmountMax: 1, Chance: 0.05},
			},
		},
	}

	for name, table := range defaults {
		data, err := yaml.Marshal(table)
		if err != nil {
			return fmt.Errorf("marshal %s loot: %w", name, err)
		}
		if err := os.WriteFile(filepath.Join(lm.dataDir, name+".yml"), data, 0644); err != nil {
			return fmt.Errorf("write %s loot: %w", name, err)
		}
	}
	return nil
}
