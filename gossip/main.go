package main

import (
	"context"
	"os"
	"os/signal"
	"slices"
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
	_updateInterval = 300 * time.Millisecond
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	for range 1 {
		SpawnPeer(ctx)
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
			spawnPeer(ctx)
		}

		if rl.IsKeyPressed(rl.KeyMinus) {
			removePeer()
		}

		mu.Unlock()

		rl.EndDrawing()
	}
}
