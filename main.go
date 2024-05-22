package main

import (
	"fmt"
	"log"
	"os"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func main() {
	f, err := os.Open("./testdemo.dem")
	if err != nil {
		log.Panic("failed to open demo file: ", err)
	}
	defer f.Close()

	p := dem.NewParser(f)
	defer p.Close()

	roundNumber := 0
	killsByPlayerInRound := make(map[int]map[string]int)
	firstKillTime := make(map[int]map[string]float64)
	lastKillTime := make(map[int]map[string]float64)

	// Handle round start
	p.RegisterEventHandler(func(e events.RoundStart) {
		roundNumber++
	})


	// Handle kills
	p.RegisterEventHandler(func(k events.Kill) {
		currentTime := float64(p.GameState().IngameTick()) / float64(p.TickRate())

		// Count kills per player per round
		killsByPlayerInRound = countKillsPerPlayerPerRound(killsByPlayerInRound, roundNumber, k.Killer.Name)

		// Calculate time between first and last kill
		firstKillTime, lastKillTime = calculateKillTimes(firstKillTime, lastKillTime, roundNumber, k.Killer.Name, currentTime)
	})


	// Handle round end
	p.RegisterEventHandler(func(e events.RoundEnd) {
		for player, kills := range killsByPlayerInRound[roundNumber] {
			if kills >= 3 {
				totalTime := lastKillTime[roundNumber][player] - firstKillTime[roundNumber][player]
				fmt.Printf("Round %d: Player %s had a multikill (%d) and it took %.0f seconds\n", roundNumber, player, kills, totalTime)
			}
		}
		delete(killsByPlayerInRound, roundNumber)
	})

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		log.Panic("failed to parse demo: ", err)
	}
}


func countKillsPerPlayerPerRound(killsByPlayerInRound map[int]map[string]int, roundNumber int, killerName string) map[int]map[string]int {
	if _, ok := killsByPlayerInRound[roundNumber]; !ok {
		killsByPlayerInRound[roundNumber] = make(map[string]int)
	}
	killsByPlayerInRound[roundNumber][killerName]++
	return killsByPlayerInRound
}

func calculateKillTimes(firstKillTime, lastKillTime map[int]map[string]float64, roundNumber int, killerName string, currentTime float64) (map[int]map[string]float64, map[int]map[string]float64) {
	if _, ok := firstKillTime[roundNumber][killerName]; !ok {
		if _, ok := firstKillTime[roundNumber]; !ok {
			firstKillTime[roundNumber] = make(map[string]float64)
		}
		firstKillTime[roundNumber][killerName] = currentTime
	}
	if _, ok := lastKillTime[roundNumber]; !ok {
		lastKillTime[roundNumber] = make(map[string]float64)
	}
	lastKillTime[roundNumber][killerName] = currentTime
	return firstKillTime, lastKillTime
}
