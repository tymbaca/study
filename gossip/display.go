package main

import (
	"image/color"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tymbaca/study/gossip/peer"
)

func drawNodes(allPeers map[string]*peer.Peer, addrs []string, positions []Vector2) {
	for i, addr := range addrs {
		peer := allPeers[addr]
		pos := positions[i]

		rl.DrawCircleV(rl.Vector2(pos), _nodeRadius, getColor(peer))
		rl.DrawText(strconv.Itoa(peer.GetSheeps()), int32(pos.X)-10, int32(pos.Y-10), _textSize, rl.Black)
		hisPeers := peer.GetPeers()
		rl.DrawText(addr, int32(pos.X)+10, int32(pos.Y+20), _addrSize, rl.DarkGreen)
		rl.DrawText(strings.Join(hisPeers, "\n"), int32(pos.X)+10, int32(pos.Y+25+_addrSize), _infoSize, rl.DarkBrown)
	}
}

func drawLinks(allPeers map[string]*peer.Peer, addrs []string, positions []Vector2) {
	for i, addr := range addrs {
		from := rl.Vector2(positions[i])
		this := allPeers[addr]
		peers := this.GetPeers()

		for _, peer := range peers {
			if addr == peer {
				continue
			}
			toIdx := slices.Index(addrs, peer)
			if toIdx < 0 {
				continue
			}

			to := rl.Vector2(positions[toIdx])
			rl.DrawLineV(from, to, rl.DarkGray)

			pointFac := (rl.Vector2Distance(from, to) - _nodeRadius - 5) / rl.Vector2Distance(from, to)
			pointPos := rl.Vector2Lerp(from, to, pointFac)
			rl.DrawCircleV(pointPos, 3, rl.DarkGray)
		}
	}
}

func getColor(peer *peer.Peer) rl.Color {
	t := peer.GetSheepsTime()
	oldness := time.Since(t)
	oldest := 4 * time.Second

	factor := rl.Clamp(float32(oldness)/float32(oldest), 0, 1)
	val := rl.Lerp(20, 255, factor) // from green to red
	return color.RGBA{R: uint8(val), G: 122, B: 122, A: 255}
}

// Vector2 struct represents a 2D vector with x and y coordinates
type Vector2 struct {
	X, Y float32
}

// CircleLayout arranges `n` points evenly spaced in a circle
// and returns their positions as a slice of Vector2
func CircleLayout(n int, radius float64, offsetX, offsetY float32) []Vector2 {
	if n <= 0 {
		return nil // No points to place
	}

	// Slice to hold the calculated positions
	positions := make([]Vector2, n)

	for i := 0; i < n; i++ {
		// Calculate angle for each point
		angle := 2 * math.Pi * float64(i) / float64(n)
		// Calculate x and y coordinates
		x := radius * math.Cos(angle)
		y := radius * math.Sin(angle)
		positions[i] = Vector2{X: float32(x) + offsetX, Y: float32(y) + offsetY}
	}

	return positions
}
