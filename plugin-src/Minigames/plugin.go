package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Origin-Net/FernMC/server"
	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/plugin-src/Minigames/framework"
	"github.com/Origin-Net/FernMC/plugin-src/Minigames/games/skywars"
	"github.com/sandertv/gophertunnel/minecraft/resource"
)

var fwGlobal *framework.Framework


var Plugin server.Plugin = minigamesPlugin{}

var _ server.Plugin = minigamesPlugin{}

type minigamesPlugin struct{}

func (minigamesPlugin) Meta() server.PluginMeta {
	return server.PluginMeta{
		Name:        "Minigames",
		Version:     "1.0.0",
		Description: "Minigames framework for FernMC",
		Authors:     []string{"FernMC Server"},
	}
}

func (minigamesPlugin) OnLoad(ctx server.PluginContext) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logWrapper := framework.NewLogWrapper(log)

	fw, err := framework.New(ctx.DataFolder, logWrapper, ctx.Server.World().BlockRegistry())
	if err != nil {
		log.Error("Failed to initialise framework", "error", err)
		return
	}
	fw.SetDefaultWorld(ctx.Server.World())

	skywarsGame := skywars.New(fw)
	if err := fw.Games().Registry().Register(skywarsGame); err != nil {
		log.Error("Failed to register SkyWars", "error", err)
		return
	}

	fwGlobal = fw
	fw.SetServer(ctx.Server)

	
	packDir := "resources"
	entries, err := os.ReadDir(packDir)
	if err == nil {
		for _, entry := range entries {
			packPath := filepath.Join(packDir, entry.Name())
			p, err := resource.ReadPath(packPath)
			if err != nil {
				log.Warn("Failed to load resource pack", "path", packPath, "error", err)
				continue
			}
			ctx.Server.AddResourcePack(p)
			log.Info("Loaded resource pack", "name", p.Name(), "path", packPath)
		}
	}
}

func (minigamesPlugin) OnEnable() {
	if fwGlobal == nil {
		return
	}

	cmd.Register(cmd.New("join", "Join a minigame", []string{"j"},
		joinCmd{fw: fwGlobal},
	))
	cmd.Register(cmd.New("leave", "Leave the current game", []string{"l"},
		leaveCmd{fw: fwGlobal},
	))
	cmd.Register(cmd.New("games", "List available minigames", nil,
		gamesCmd{fw: fwGlobal},
	))
	cmd.Register(cmd.New("setspawn", "Set the lobby spawn point", nil,
		setSpawnCmd{fw: fwGlobal},
	))

	
	for _, hc := range framework.HologramCommands(fwGlobal.Holograms()) {
		cmd.Register(hc)
	}

	fwGlobal.Logger().Info("Minigames plugin enabled")

	if err := fwGlobal.SetupLobby(); err != nil {
		fwGlobal.Logger().Error("Failed to setup lobby", "error", err)
	}

	if err := fwGlobal.DiscoverMaps(); err != nil {
		fwGlobal.Logger().Error("Map discovery failed", "error", err)
	}

	fwGlobal.AttachLobbyHandler()
}

func (minigamesPlugin) OnDisable() {
	if fwGlobal != nil {
		fwGlobal.Shutdown()
	}
}

func (minigamesPlugin) OnUnload() {}
