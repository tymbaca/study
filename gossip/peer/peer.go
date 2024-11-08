package peer

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

// TODO: we need reverce gossip mechanism: if p1 gossips to p2 and it appears
// that p2 has newer gossip it must return it to p1.

var ErrDown = fmt.Errorf("node is down")

type Peer struct {
	me     string
	ctx    context.Context
	cancel context.CancelFunc

	mu        sync.RWMutex
	sheeps    Sheeps
	peers     PeersList
	transport transport
}

func New(ctx context.Context, addr string, transport transport) *Peer {
	ctx, cancel := context.WithCancel(ctx)
	return &Peer{
		me:        addr,
		ctx:       ctx,
		cancel:    cancel,
		transport: transport,
		peers:     Gossip[map[string]struct{}]{Val: map[string]struct{}{addr: {}}},
	}
}

type transport interface {
	SetSheeps(sender string, peer string, sheeps Sheeps) error
	SetPeers(sender string, peer string, peers PeersList) error
}

func (p *Peer) Launch(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for range t.C {
		for peer := range p.GetPeersMap() {
			if peer != p.me {
				fmt.Println(p.me, "got inself in random get")
				continue
			}

			if p.ctx.Err() != nil {
				return
			}

			if err := p.transport.SetSheeps(p.me, peer, p.sheeps); errors.Is(err, ErrDown) {
				p.RemovePeer(peer)
			}

			if err := p.transport.SetPeers(p.me, peer, p.peers); errors.Is(err, ErrDown) {
				p.RemovePeer(peer)
			}
		}
	}
}

func (p *Peer) Stop() {
	p.cancel()
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

func (p *Peer) SetPeers(sender string, newPeers PeersList) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// bad
	if sender != "" {
		p.peers.Val[sender] = struct{}{}
		p.peers.Time = time.Now()
	}

	if p.peers.Time.After(newPeers.Time) {
		return
	}

	p.peers = newPeers
}

func (p *Peer) AddPeer(sender, peer string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.peers.Val[peer] = struct{}{}
	if sender != peer && sender != "" {
		p.peers.Val[sender] = struct{}{}
	}
	p.peers.Time = time.Now()
}

func (p *Peer) RemovePeer(addr string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for peerAddr := range p.peers.Val { // FIX: out of bounce if quickly remove twice
		if peerAddr == addr {
			delete(p.peers.Val, peerAddr)
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

func (p *Peer) GetPeersMap() map[string]struct{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	m := make(map[string]struct{}, len(p.peers.Val))
	for addr := range p.peers.Val {
		m[addr] = struct{}{}
	}

	return m
}

func (p *Peer) GetPeersList() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.getPeersList()
}

func (p *Peer) getPeersList() []string {
	list := lo.Keys(p.peers.Val)
	slices.Sort(list)

	return list
}

func (p *Peer) GetSheepsTime() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.sheeps.Time
}

type (
	Sheeps    = Gossip[int]
	PeersList = Gossip[map[string]struct{}]
)

type Gossip[T any] struct {
	Val  T
	Time time.Time
}

func (p *Peer) getRandomPeer() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, addr := range lo.Shuffle(p.getPeersList()) {
		if addr != p.me {
			return addr
		}
	}

	return p.me
}

func toBytes[T any](val T) []byte {
	b := bytes.NewBuffer(nil)
	err := gob.NewEncoder(b).Encode(val)
	if err != nil {
		panic(err)
	}

	return b.Bytes()
}
