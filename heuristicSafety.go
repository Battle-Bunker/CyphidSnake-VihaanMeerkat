package main

import (
  "github.com/BattlesnakeOfficial/rules"
  "github.com/Battle-Bunker/cyphid-snake/agent"
  "math"
)

// HeuristicSafeSpace calculates a score based on the amount of safe space available to our team's snakes
func HeuristicSafeSpace(snapshot agent.GameSnapshot) float64 {
  var totalScore float64

  for _, allySnake := range snapshot.YourTeam() {
    if !allySnake.Alive() {
      continue
    }

    safeSpace := calculateSafeSpace(allySnake, snapshot)
    healthFactor := float64(allySnake.Health()) / 100.0
    snakeScore := safeSpace * healthFactor

    totalScore += snakeScore
  }

  return totalScore
}

func calculateSafeSpace(snake agent.SnakeSnapshot, snapshot agent.GameSnapshot) float64 {
  head := snake.Head()
  safeSpace := 0.0

  for x := 0; x < snapshot.Width(); x++ {
    for y := 0; y < snapshot.Height(); y++ {
      point := rules.Point{X: x, Y: y}
      if isSafe(point, snake, snapshot) {
        distance := manhattanDistance(head, point)
        safeSpace += 1.0 / (1.0 + distance) // Closer safe spaces are more valuable
      }
    }
  }

  return safeSpace
}

func isSafe(point rules.Point, snake agent.SnakeSnapshot, snapshot agent.GameSnapshot) bool {
  // Check if the point is not occupied by any snake body
  for _, otherSnake := range snapshot.AllSnakes() {
    if containsPoint(otherSnake.Body(), point) {
      return false
    }
  }

  // Check if the point is not a hazard
  for _, hazard := range snapshot.Hazards() {
    if hazard == point {
      return false
    }
  }

  // Check if the point is not adjacent to an enemy snake's head (unless we're longer)
  for _, enemy := range snapshot.Opponents() {
    if enemy.Alive() && manhattanDistance(enemy.Head(), point) == 1 && enemy.Length() >= snake.Length() {
      return false
    }
  }

  return true
}

func manhattanDistance(a, b rules.Point) float64 {
  return math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y))
}

func containsPoint(points []rules.Point, target rules.Point) bool {
  for _, p := range points {
    if p == target {
      return true
    }
  }
  return false
}