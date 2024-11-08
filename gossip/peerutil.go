package main

import (
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
	newAddr := gofakeit.Name()

	randomPeer := choosePeer()
	if randomPeer != nil {
		randomPeer.AddPeer(newAddr)
	}

	newPeer := peer.New(newAddr, &mapTransport{})
	peers[newAddr] = newPeer
	go newPeer.Launch(_updateInterval)
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

func (t *mapTransport) SetSheeps(addr string, sheeps peer.Sheeps) error {
	mu.RLock()
	defer mu.RUnlock()

	toPeer, ok := peers[addr]
	if !ok {
		return peer.ErrDown
	}

	toPeer.SetSheeps(sheeps)
	return nil
}

func (t *mapTransport) SetPeers(addr string, addrs peer.PeersList) error {
	mu.RLock()
	defer mu.RUnlock()

	toPeer, ok := peers[addr]
	if !ok {
		return peer.ErrDown
	}

	toPeer.SetPeers(addrs)
	return nil
}
