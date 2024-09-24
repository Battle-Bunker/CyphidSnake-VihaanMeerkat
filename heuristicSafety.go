package main

import (
  "github.com/BattlesnakeOfficial/rules"
  "github.com/Battle-Bunker/cyphid-snake/agent"
  "math"
)

const (
  RecentMovesCount = 3 // Number of recent moves to consider
)

// HeuristicSafeSpaceWithVariedMovement calculates a score based on safe space, distance from dangers, and movement variety
func HeuristicSafeSpaceWithVariedMovement(snapshot agent.GameSnapshot) float64 {
  var totalScore float64

  for _, allySnake := range snapshot.YourTeam() {
    if !allySnake.Alive() {
      continue
    }

    head := allySnake.Head()
    safetyScore := calculateSafetyScore(snapshot, head)
    spaceScore := calculateSpaceScore(snapshot, head)
    movementVarietyScore := calculateMovementVarietyScore(allySnake)

    // Combine safety, space, and movement variety scores
    snakeScore := (safetyScore*2 + spaceScore + movementVarietyScore) / 4

    // Scale the score by the snake's health to prioritize survival
    totalScore += snakeScore * float64(allySnake.Health())
  }

  return totalScore
}

func calculateSafetyScore(snapshot agent.GameSnapshot, point rules.Point) float64 {
  width := snapshot.Width()
  height := snapshot.Height()

  // Calculate distance from borders
  borderDistance := math.Min(
    math.Min(float64(point.X), float64(width-1-point.X)),
    math.Min(float64(point.Y), float64(height-1-point.Y)),
  )

  // Calculate distance from other snakes
  minSnakeDistance := math.Inf(1)
  for _, snake := range snapshot.Opponents() {
    if !snake.Alive() {
      continue
    }
    for _, bodyPart := range snake.Body() {
      distance := manhattanDistance(point, bodyPart)
      minSnakeDistance = math.Min(minSnakeDistance, float64(distance))
    }
  }

  // Combine border and snake distances, normalizing to 0-1 range
  safetyScore := (borderDistance/float64(width/2) + minSnakeDistance/float64(width+height)) / 2
  return math.Min(safetyScore, 1.0) // Cap at 1.0
}

func calculateSpaceScore(snapshot agent.GameSnapshot, point rules.Point) float64 {
  width := snapshot.Width()
  height := snapshot.Height()
  totalSpace := width * height

  occupiedSpace := make(map[rules.Point]bool)

  // Mark occupied spaces
  for _, snake := range snapshot.AllSnakes() {
    if !snake.Alive() {
      continue
    }
    for _, bodyPart := range snake.Body() {
      occupiedSpace[bodyPart] = true
    }
  }

  // Flood fill to calculate available space
  availableSpace := floodFill(point, occupiedSpace, width, height)

  // Normalize available space to 0-1 range
  spaceScore := float64(availableSpace) / float64(totalSpace)
  return spaceScore
}

func calculateMovementVarietyScore(snake agent.SnakeSnapshot) float64 {
  body := snake.Body()
  if len(body) < RecentMovesCount+1 {
    return 1.0 // Not enough moves to calculate, return max score
  }

  recentMoves := make([]rules.Point, RecentMovesCount)
  for i := 0; i < RecentMovesCount; i++ {
    recentMoves[i] = rules.Point{
      X: body[i].X - body[i+1].X,
      Y: body[i].Y - body[i+1].Y,
    }
  }

  uniqueMoves := make(map[rules.Point]bool)
  for _, move := range recentMoves {
    uniqueMoves[move] = true
  }

  varietyScore := float64(len(uniqueMoves)) / float64(RecentMovesCount)
  return varietyScore
}

func manhattanDistance(p1, p2 rules.Point) int {
  return abs(p1.X-p2.X) + abs(p1.Y-p2.Y)
}

func abs(x int) int {
  if x < 0 {
    return -x
  }
  return x
}

func floodFill(start rules.Point, occupied map[rules.Point]bool, width, height int) int {
  if occupied[start] {
    return 0
  }

  queue := []rules.Point{start}
  visited := make(map[rules.Point]bool)
  count := 0

  for len(queue) > 0 {
    current := queue[0]
    queue = queue[1:]

    if visited[current] || occupied[current] ||
      current.X < 0 || current.X >= width ||
      current.Y < 0 || current.Y >= height {
      continue
    }

    visited[current] = true
    count++

    queue = append(queue, rules.Point{X: current.X + 1, Y: current.Y})
    queue = append(queue, rules.Point{X: current.X - 1, Y: current.Y})
    queue = append(queue, rules.Point{X: current.X, Y: current.Y + 1})
    queue = append(queue, rules.Point{X: current.X, Y: current.Y - 1})
  }

  return count
}