package peer

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"sync"
	"time"
)

var ErrDown = fmt.Errorf("node is down")

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
	SetSheeps(peer string, sheeps Sheeps) error
	SetPeers(peer string, peers PeersList) error
}

func (p *Peer) Launch(interval time.Duration) {
	for range time.Tick(interval) {
		peer := p.peers.Val[rand.Intn(len(p.peers.Val))]
		if peer == p.me {
			continue
		}

		if err := p.transport.SetSheeps(peer, p.sheeps); errors.Is(err, ErrDown) {
			p.RemovePeer(peer)
		}

		if err := p.transport.SetPeers(peer, p.peers); errors.Is(err, ErrDown) {
			p.RemovePeer(peer)
		}
	}
}

func (p *Peer) Addr() string {
	return p.me
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

func (p *Peer) RemovePeer(addr string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, peerAddr := range p.peers.Val { // FIX: out of bounce if quickly remove twice
		if peerAddr == addr {
			p.peers.Val = slices.Delete(p.peers.Val, i, i+1)
			p.peers.Time = time.Now()
			return
		}
	}
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

func (p *Peer) GetSheepsTime() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.sheeps.Time
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