package main

import (
  "math"

  "github.com/BattlesnakeOfficial/rules"
  "github.com/Battle-Bunker/cyphid-snake/agent"
)

// HeuristicEnemyDistance calculates a score based on the average distance
// between your team's snakes and enemy snakes.
func HeuristicEnemyDistance(snapshot agent.GameSnapshot) float64 {
  yourTeam := snapshot.YourTeam()
  opponents := snapshot.Opponents()

  if len(yourTeam) == 0 || len(opponents) == 0 {
    return 0
  }

  totalDistance := 0.0
  count := 0

  for _, ally := range yourTeam {
    for _, enemy := range opponents {
      distance := manhattanDistance(ally.Head(), enemy.Head())
      totalDistance += distance
      count++
    }
  }

  averageDistance := totalDistance / float64(count)

  // Normalize the score to be between 0 and 100
  maxPossibleDistance := float64(snapshot.Width() + snapshot.Height())
  normalizedScore := (averageDistance / maxPossibleDistance) * 100

  return normalizedScore
}

// manhattanDistance calculates the Manhattan distance between two points
func manhattanDistance1(p1, p2 rules.Point) float64 {
  return math.Abs(float64(p1.X-p2.X)) + math.Abs(float64(p1.Y-p2.Y))
}