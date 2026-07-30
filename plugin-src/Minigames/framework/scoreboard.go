package framework

import (
	"fmt"
	"time"

	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/player/scoreboard"
)


func (pm *PlayerManager) SendLobbyScoreboard(p *player.Player) {
	var stats *PlayerStats
	if pm.stats != nil {
		stats = pm.stats.LoadPlayerStats(p.UUID())
	} else {
		stats = &PlayerStats{}
	}
	kd := 0.0
	if stats.Deaths > 0 {
		kd = float64(stats.Kills) / float64(stats.Deaths)
	} else if stats.Kills > 0 {
		kd = float64(stats.Kills)
	}

	pm.mu.RLock()
	srv := pm.srv
	pm.mu.RUnlock()
	online := 0
	if srv != nil {
		online = srv.PlayerCount()
	}

	sb := scoreboard.New("§6§lSKYWARS")
	sb.Set(0, "")
	sb.Set(1, "§fYour Stats")
	sb.Set(2, fmt.Sprintf("§fWins: §a%d", stats.Wins))
	sb.Set(3, fmt.Sprintf("§fKills: §c%d", stats.Kills))
	sb.Set(4, fmt.Sprintf("§fK/D: §e%.2f", kd))
	sb.Set(5, "")
	sb.Set(6, fmt.Sprintf("§fPlayers Online: §a%d", online))
	sb.Set(7, "")
	sb.Set(8, "§7play.originpvp.net")
	sb.RemovePadding()
	p.SendScoreboard(sb)
}


func SendMatchScoreboard(p *player.Player, matchKills int, aliveCount int, chestTimeLeft time.Duration) {
	sb := scoreboard.New("§6§lSKYWARS")
	sb.Set(0, "")
	sb.Set(1, fmt.Sprintf("§fPlayers: §a%d", aliveCount))
	sb.Set(2, fmt.Sprintf("§fKills: §c%d", matchKills))
	sb.Set(3, "")

	chestStr := "§eN/A"
	if chestTimeLeft > 0 {
		minutes := int(chestTimeLeft.Minutes())
		seconds := int(chestTimeLeft.Seconds()) % 60
		chestStr = fmt.Sprintf("§e%02d:%02d", minutes, seconds)
	} else if chestTimeLeft == 0 {
		chestStr = "§aRefilled!"
	}
	sb.Set(4, "§fChest Refill:")
	sb.Set(5, chestStr)
	sb.Set(6, "")
	sb.Set(7, "§7play.originmc.net")
	sb.RemovePadding()
	p.SendScoreboard(sb)
}


func RemoveScoreboard(p *player.Player) {
	p.RemoveScoreboard()
}
