package main

import (
  "math"

  "github.com/BattlesnakeOfficial/rules"
  "github.com/Battle-Bunker/cyphid-snake/agent"
)

// HeuristicDistanceFromEnemies calculates a score based on how far our snakes are from enemy snakes.
// A higher score indicates a better position (further from enemies).
func HeuristicDistanceFromEnemies(snapshot agent.GameSnapshot) float64 {
  var totalScore float64

  for _, allySnake := range snapshot.YourTeam() {
    if !allySnake.Alive() {
      continue
    }

    allyHead := allySnake.Head()
    minDistance := math.Inf(1)

    for _, enemySnake := range snapshot.Opponents() {
      if !enemySnake.Alive() {
        continue
      }

      enemyHead := enemySnake.Head()
      distance := manhattanDistance(allyHead, enemyHead)

      if distance < int(minDistance) {
          minDistance = float64(distance)
      }
    }

    // We want a higher score for greater distances
    if minDistance != math.Inf(1) {
      totalScore += minDistance
    }
  }

  return totalScore
}

// manhattanDistance calculates the Manhattan distance between two points
func manhattanDistance1(p1, p2 rules.Point) float64 {
  return math.Abs(float64(p1.X-p2.X)) + math.Abs(float64(p1.Y-p2.Y))
}