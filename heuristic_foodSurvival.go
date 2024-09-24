package main

import (
  "github.com/BattlesnakeOfficial/rules"
  "github.com/Battle-Bunker/cyphid-snake/agent"
  "math"
)

// HeuristicFoodAndSurvival calculates a heuristic score based on health, food proximity, and survival potential
func HeuristicFoodAndSurvival(snapshot agent.GameSnapshot) float64 {
  var score float64
  you := snapshot.You()

  // Health component
  healthScore := float64(you.Health())
  if you.Health() < 70 {
    healthScore *= 1.5 // Increase importance of health when below 70
  }
  score += healthScore

  // Food proximity and safety component
  foodScore := calculateFoodScore(snapshot, you)
  score += foodScore

  // Survival potential component
  survivalScore := calculateSurvivalScore(snapshot, you)
  score += survivalScore

  return score
}

func calculateFoodScore(snapshot agent.GameSnapshot, you agent.SnakeSnapshot) float64 {
  var foodScore float64
  for _, food := range snapshot.Food() {
    distanceToFood := manhattanDistance(you.Head(), food)
    if distanceToFood == 0 {
      continue // Skip food at the same position as the head
    }

    // Penalize food close to borders
    borderPenalty := calculateBorderPenalty(food, snapshot.Width(), snapshot.Height())

    // Penalize food close to enemies
    enemyPenalty := calculateEnemyPenalty(food, snapshot.Opponents())

    // Calculate food value (inverse of distance, adjusted by penalties)
    foodValue := 100.0 / (float64(distanceToFood) * (1 + borderPenalty + enemyPenalty))
    foodScore += foodValue
  }
  return foodScore
}

func calculateSurvivalScore(snapshot agent.GameSnapshot, you agent.SnakeSnapshot) float64 {
  var survivalScore float64
  numForwardMoves := len(you.ForwardMoves())
  survivalScore += float64(numForwardMoves) * 10 // Reward having more move options

  // Compare food count with other snakes
  yourFoodCount := you.Length() - 3 // Assuming initial length is 3
  for _, opponent := range snapshot.Opponents() {
    opponentFoodCount := opponent.Length() - 3
    if yourFoodCount > opponentFoodCount {
      survivalScore += 50 // Bonus for having more food than an opponent
    }
  }

  return survivalScore
}

func manhattanDistance2(p1, p2 rules.Point) int {
  return int(math.Abs(float64(p1.X-p2.X)) + math.Abs(float64(p1.Y-p2.Y)))
}

func calculateBorderPenalty(food rules.Point, width, height int) float64 {
  distToBorder := math.Min(
    math.Min(float64(food.X), float64(width-1-food.X)),
    math.Min(float64(food.Y), float64(height-1-food.Y)),
  )
  return math.Max(0, 2-distToBorder) // Penalty decreases as distance to border increases
}

func calculateEnemyPenalty(food rules.Point, opponents []agent.SnakeSnapshot) float64 {
  var minDistance int = math.MaxInt32
  for _, opponent := range opponents {
    dist := manhattanDistance(food, opponent.Head())
    if dist < minDistance {
      minDistance = dist
    }
  }
  return math.Max(0, 5-float64(minDistance)) // Penalty decreases as distance to nearest enemy increases
}