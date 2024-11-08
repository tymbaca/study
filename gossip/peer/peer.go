package peer

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Peer struct {
	me string

	mu        sync.RWMutex
	sheeps    Gossip[int]
	peers     Gossip[[]string]
	transport transport
}

func New(addr string, transport transport) *Peer {
	return &Peer{
		me:        addr,
		transport: transport,
		peers:     Gossip[[]string]{Val: []string{addr}},
	}
}

type transport interface {
	SetSheeps(peer string, sheeps Sheeps)
	SetPeers(peer string, peers PeersList)
}

func (p *Peer) Launch(interval time.Duration) {
	for range time.Tick(interval) {
		peer := p.peers.Val[rand.Intn(len(p.peers.Val))]
		if peer == p.me {
			continue
		}

		fmt.Println(p.me, "sheep count is", p.sheeps.Val)
		// fmt.Println(p.me, "it talking to", peer)

		p.transport.SetSheeps(peer, p.sheeps)
		p.transport.SetPeers(peer, p.peers)
	}
}

func (p *Peer) SetSheeps(newSheeps Sheeps) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.sheeps.Time.After(newSheeps.Time) {
		return
	}

	p.sheeps = newSheeps
}

func (p *Peer) SetPeers(newPeers PeersList) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.peers.Time.After(newPeers.Time) {
		return
	}

	p.peers = newPeers
}

func (p *Peer) AddPeer(peer string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.peers.Val = append(p.peers.Val, peer)
	p.peers.Time = time.Now()
}

func (p *Peer) GetSheeps() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.sheeps.Val
}

func (p *Peer) GetPeers() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.peers.Val
}

func (p *Peer) GetPeersListTime() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.peers.Time
}

type (
	Sheeps    = Gossip[int]
	PeersList = Gossip[[]string]
)

type Gossip[T any] struct {
	Val  T
	Time time.Time
}

func toBytes[T any](val T) []byte {
	b := bytes.NewBuffer(nil)
	err := gob.NewEncoder(b).Encode(val)
	if err != nil {
		panic(err)
	}

	return b.Bytes()
}
