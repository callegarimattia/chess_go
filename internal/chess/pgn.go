package chess

import (
	"fmt"
	"strings"
	"time"
)

// buildPGN constructs a minimal PGN string from the game.
func buildPGN(g Game) string {
	var sb strings.Builder

	result := g.Result()
	resultStr := "*"
	switch result {
	case WhiteWins:
		resultStr = "1-0"
	case BlackWins:
		resultStr = "0-1"
	case Stalemate, DrawFiftyMove, DrawThreefoldRepetition, DrawInsufficientMaterial:
		resultStr = "1/2-1/2"
	}

	// Tag pairs.
	fmt.Fprintf(&sb, "[Event \"?\"]\n")
	fmt.Fprintf(&sb, "[Site \"?\"]\n")
	fmt.Fprintf(&sb, "[Date \"%s\"]\n", time.Now().Format("2006.01.02"))
	fmt.Fprintf(&sb, "[Round \"?\"]\n")
	fmt.Fprintf(&sb, "[White \"?\"]\n")
	fmt.Fprintf(&sb, "[Black \"?\"]\n")
	fmt.Fprintf(&sb, "[Result \"%s\"]\n", resultStr)
	fmt.Fprintf(&sb, "\n")

	// Move text (reconstructed from history).
	// For the skeleton, we emit UCI strings as a placeholder.
	// Full SAN requires SANString() which is implemented separately.
	moveNum := 1
	for i, state := range g.history {
		_ = state
		if i%2 == 0 {
			fmt.Fprintf(&sb, "%d. ", moveNum)
		}
		// We don't store moves explicitly in history; this is a skeleton.
		// Full implementation will track moves separately.
		if i%2 == 1 {
			moveNum++
		}
	}
	fmt.Fprintf(&sb, "%s\n", resultStr)

	return sb.String()
}
