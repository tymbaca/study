package main

import (
	"image/color"
	"slices"
	"strconv"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/samber/lo"
	"github.com/tymbaca/study/gossip/peer"
	"golang.org/x/exp/rand"
)

var (
	peers = map[string]*peer.Peer{}
	mu    = new(sync.RWMutex)
)

const (
	_nodeRadius     = 20
	_textSize       = 20
	_captionSize    = 8
	_updateInterval = 100 * time.Millisecond
)

func main() {
	for range 1 {
		SpawnPeer()
	}

	addrs := lo.Keys(peers)
	for key := range peers {
		peers[key].SetPeers(peer.Gossip[[]string]{Val: addrs, Time: time.Now()})
		go peers[key].Launch(_updateInterval)
	}

	entry := ChoosePeer()
	entry.SetSheeps(peer.Gossip[int]{Val: 10, Time: time.Now()})

	// go func() {
	// 	for range time.Tick(1500 * time.Millisecond) {
	// 		entry.SetSheeps(peer.Gossip[int]{Val: rand.Intn(100), Time: time.Now()})
	// 	}
	// }()

	// go func() {
	// 	for range time.Tick(5000 * time.Millisecond) {
	// 		SpawnPeer()
	// 	}
	// }()

	// go func() {
	// 	for range time.Tick(1700 * time.Millisecond) {
	// 		RemovePeer()
	// 	}
	// }()

	//--------------------------------------------------------------------------------------------------

	rl.InitWindow(800, 600, "gossip")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawFPS(10, 10)

		mu.Lock()
		positions := CircleLayout(len(peers), float64(rl.Lerp(100, 500, float32(len(peers))/100)), 400, 300)
		addrs := lo.Keys(peers)
		slices.Sort(addrs)

		// Draw links
		drawLinks(peers, addrs, positions)

		drawNodes(peers, addrs, positions)

		if rl.IsKeyPressed(rl.KeySpace) {
			choosePeer().SetSheeps(peer.Gossip[int]{Val: rand.Intn(100), Time: time.Now()})
		}

		if rl.IsKeyPressed(rl.KeyEqual) {
			spawnPeer()
		}

		if rl.IsKeyPressed(rl.KeyMinus) {
			removePeer()
		}

		mu.Unlock()

		rl.EndDrawing()
	}
}

func drawNodes(allPeers map[string]*peer.Peer, addrs []string, positions []Vector2) {
	for i, addr := range addrs {
		peer := peers[addr]
		pos := positions[i]

		rl.DrawCircleV(rl.Vector2(pos), _nodeRadius, getColor(peer))
		rl.DrawText(strconv.Itoa(peer.GetSheeps()), int32(pos.X)-10, int32(pos.Y-10), _textSize, rl.Black)
		rl.DrawText(addr, int32(pos.X)+10, int32(pos.Y+20), _captionSize, rl.DarkGreen)
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
	oldest := 10 * time.Second

	factor := float32(oldness) / float32(oldest)
	val := rl.Lerp(122, 255, factor)
	return color.RGBA{R: uint8(val), G: 122, B: 122, A: 255}
}
