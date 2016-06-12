package main

import (
	"flag"
	"fmt"
)

var (
	maxBet                                                                                   int
	teamA, teamB                                                                             string
	initialOddsA, initialOddsB, initialOddsDraw, currentOddsA, currentOddsDraw, currentOddsB float64
)

func init() {
	flag.IntVar(&maxBet, "maxBet", 10, "the maximum amount you want to bet")
	flag.StringVar(&teamA, "teamA", "Albania", "team A's name")
	flag.StringVar(&teamB, "teamB", "Belgium", "team B's name")
	flag.Float64Var(&initialOddsA, "initialOddsA", 1.5, "initial odds for team A to win")
	flag.Float64Var(&initialOddsB, "initialOddsB", 1.5, "inital odds for team B to win")
	flag.Float64Var(&initialOddsDraw, "initialOddsDraw", 1.5, "initial for team A and team B to draw")
	flag.Float64Var(&currentOddsA, "currentOddsA", 1.5, "odds for team A to win")
	flag.Float64Var(&currentOddsB, "currentOddsB", 1.5, "odds for team B to win")
	flag.Float64Var(&currentOddsDraw, "currentOddsDraw", 1.5, "odds for team A and team B to draw")
}

func main() {
	flag.Parse()
	fmt.Printf("Win rate %+v\n", getInitialProjectedWinRate(teamA, teamB, initialOddsA, initialOddsB, initialOddsDraw))
	strategy := getBestStrategy(teamA, teamB, maxBet, getInitialProjectedWinRate(teamA, teamB, initialOddsA, initialOddsB, initialOddsDraw), currentOddsA, currentOddsB, currentOddsDraw)
	fmt.Printf("The best strategy %+v\n", strategy)
}

func getInitialProjectedWinRate(teamA, teamB string, initialOddsA, initialOddsB, initialOddsDraw float64) map[string]float64 {
	denom := initialOddsA*initialOddsB + initialOddsA*initialOddsDraw + initialOddsB*initialOddsDraw

	result := make(map[string]float64)
	result[teamA] = (initialOddsB * initialOddsDraw) / denom
	result[teamB] = (initialOddsA * initialOddsDraw) / denom
	result["draw"] = (initialOddsA * initialOddsB) / denom
	return result
}

func getBestStrategy(teamA, teamB string, maxBet int, projectedRate map[string]float64, currentOddsA, currentOddsB, currentOddsDraw float64) map[string]float64 {
	var x, y int
	result := make(map[string]float64)
	maxBetF := float64(maxBet)
	maxWin := 0.0

	for x = 0; x < maxBet; x++ {
		for y = 0; y < maxBet-x; y++ {
			teamAWinGain := float64(x)*currentOddsA - maxBetF
			teamBWinGain := float64(y)*currentOddsB - maxBetF
			drawWinGain := float64(maxBet-x-y)*currentOddsDraw - maxBetF

			if teamAWinGain > 0 && teamBWinGain > 0 && drawWinGain > 0 {
				fmt.Printf("Bet %d on %s, %d on %s and %d on draw\n", x, teamA, y, teamB, maxBet-x-y)
				projectedGain := projectedRate[teamA]*teamAWinGain + projectedRate[teamB]*teamBWinGain + projectedRate["draw"]*drawWinGain
				if maxWin < projectedGain {
					maxWin = projectedGain
					result[teamA] = float64(x)
					result[teamB] = float64(y)
					result["draw"] = float64(maxBet - x - y)
				}
			}
		}
	}

	fmt.Printf("Expected gain: %f\n", maxWin)
	return result
}
