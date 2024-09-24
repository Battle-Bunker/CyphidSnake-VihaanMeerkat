package main

import (
  "github.com/BattlesnakeOfficial/rules"
  "github.com/Battle-Bunker/cyphid-snake/agent"
  "math"
)

// HeuristicFoodSafety calculates a heuristic score based on the team's health and food safety
func HeuristicFoodSafety(snapshot agent.GameSnapshot) float64 {
  var score float64

  for _, allySnake := range snapshot.YourTeam() {
    // Base score is the snake's current health
    snakeScore := float64(allySnake.Health())

    // If health is below 70, add bonus for nearby safe food
    if allySnake.Health() < 70 {
      snakeScore += evaluateFoodSafety(snapshot, allySnake)
    }

    score += snakeScore
  }

  return score
}

// evaluateFoodSafety calculates a bonus score based on nearby safe food
func evaluateFoodSafety(snapshot agent.GameSnapshot, snake agent.SnakeSnapshot) float64 {
  var foodBonus float64
  head := snake.Head()

  for _, food := range snapshot.Food() {
    distance := manhattanDistance(head, food)
    if distance > 5 {
      continue // Only consider nearby food
    }

    safety := evaluateFoodPositionSafety(snapshot, food)

    // The closer and safer the food, the higher the bonus
    foodBonus += float64(6 - distance) * safety
  }

  return foodBonus
}

// evaluateFoodPositionSafety checks if the food is in a safe position
func evaluateFoodPositionSafety(snapshot agent.GameSnapshot, food rules.Point) float64 {
  safety := 1.0

  // Reduce safety for food near borders
  if food.X <= 1 || food.X >= snapshot.Width()-2 || food.Y <= 1 || food.Y >= snapshot.Height()-2 {
    safety *= 0.5
  }

  // Reduce safety for food near enemy snakes
  for _, enemy := range snapshot.Opponents() {
    if manhattanDistance(enemy.Head(), food) < 3 {
      safety *= 0.5
    }
  }

  return safety
}

// manhattanDistance calculates the Manhattan distance between two points
func manhattanDistance2(p1, p2 rules.Point) int {
  return int(math.Abs(float64(p1.X-p2.X)) + math.Abs(float64(p1.Y-p2.Y)))
}