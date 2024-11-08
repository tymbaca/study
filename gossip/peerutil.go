package main

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/tymbaca/study/gossip/peer"
)

func ChoosePeer() *peer.Peer {
	mu.RLock()
	defer mu.RUnlock()

	return choosePeer()
}

func choosePeer() *peer.Peer {
	for _, peer := range peers {
		return peer
	}

	return nil
}

func SpawnPeer() {
	mu.Lock()
	defer mu.Unlock()

	spawnPeer()
}

func spawnPeer() {
	// newPeer := gofakeit.UUID()
	newPeer := gofakeit.Name()

	randomPeer := choosePeer()
	if randomPeer != nil {
		randomPeer.AddPeer(newPeer)
	}

	peers[newPeer] = peer.New(newPeer, &mapTransport{})
}

func RemovePeer() {
	mu.Lock()
	defer mu.Unlock()

	removePeer()
}

func removePeer() {
	for addr := range peers {
		delete(peers, addr)
		break
	}
}

type mapTransport struct{}

func (t *mapTransport) SetSheeps(addr string, sheeps peer.Sheeps) {
	mu.RLock()
	defer mu.RUnlock()

	peer, ok := peers[addr]
	if !ok {
		fmt.Println("got unknown peer: " + addr)
		return
	}

	peer.SetSheeps(sheeps)
}

func (t *mapTransport) SetPeers(addr string, addrs peer.PeersList) {
	mu.RLock()
	defer mu.RUnlock()

	peer, ok := peers[addr]
	if !ok {
		fmt.Println("got unknown peer: " + addr)
		return
	}

	peer.SetPeers(addrs)
}
