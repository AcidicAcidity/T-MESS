package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type Node struct {
	Host   host.Host
	DHT    *dht.IpfsDHT
	PubSub *pubsub.PubSub
	ctx    context.Context
	cancel context.CancelFunc
}

func NewNode(ctx context.Context, privKey interface{}) (*Node, error) {
	ctx, cancel := context.WithCancel(ctx)

	// Создаём хост с случайным портом
	host, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.Identity(privKey.(libp2p.PrivKey)),
	)
	if err != nil {
		cancel()
		return nil, err
	}

	// Создаём DHT
	kadDHT, err := dht.New(ctx, host)
	if err != nil {
		cancel()
		return nil, err
	}

	// Запускаем DHT
	if err := kadDHT.Bootstrap(ctx); err != nil {
		cancel()
		return nil, err
	}

	// Создаём PubSub
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		cancel()
		return nil, err
	}

	return &Node{
		Host:   host,
		DHT:    kadDHT,
		PubSub: ps,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (n *Node) ID() string {
	return n.Host.ID().String()
}

func (n *Node) Addrs() []string {
	addrs := make([]string, 0)
	for _, addr := range n.Host.Addrs() {
		addrs = append(addrs, addr.String())
	}
	return addrs
}

func (n *Node) Close() error {
	n.cancel()
	return n.Host.Close()
}

// ConnectToPeer подключается к другому узлу
func (n *Node) ConnectToPeer(ctx context.Context, addr string) error {
	ma, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return err
	}
	peerInfo, err := peer.AddrInfoFromP2pAddr(ma)
	if err != nil {
		return err
	}
	return n.Host.Connect(ctx, *peerInfo)
}
