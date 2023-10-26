package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"log"
	"os"
)

type MyNotifiee struct{}

func (n *MyNotifiee) Listen(_ network.Network, _ ma.Multiaddr)      {}
func (n *MyNotifiee) ListenClose(_ network.Network, _ ma.Multiaddr) {}
func (n *MyNotifiee) Connected(_ network.Network, c network.Conn) {
	fmt.Println("Connected RemotePeer:", c.RemotePeer())
}
func (n *MyNotifiee) Disconnected(_ network.Network, c network.Conn) {
	fmt.Println("Disconnected RemotePeer:", c.RemotePeer())
}

func main() {
	host, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}

	js, err := os.ReadFile("bootstrap.json")
	if err != nil {
		panic(err)
	}

	nodeAPeerInfo := &peer.AddrInfo{}
	err = json.Unmarshal(js, nodeAPeerInfo)
	if err != nil {
		panic(err)
	}

	if err := host.Connect(context.Background(), *nodeAPeerInfo); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Client Node ID: %+v\n", host.ID())

	//connManager, err := connmgr.NewConnManager(1000, 1000, connmgr.WithGracePeriod(0))
	//if err != nil {
	//	panic(err)
	//}

	host.Network().Notify(&MyNotifiee{})

	peerChan := initDHT(host, "client")

	defer func() {
		fmt.Println("Shutting down...")
		host.Close()
	}()
	//peerChan := initMDNS(host, "client")
	for { // allows multiple peers to join
		peer_ := <-peerChan // will block until we discover a peer
		if peer_.ID.Validate() == peer.ErrEmptyPeerID {
			continue
		}
		fmt.Println("Found peer:", peer_, ", connecting")

		if err := host.Connect(context.Background(), peer_); err != nil {
			fmt.Println("Connection failed:", err)
			continue
		}

		fmt.Println("Connected to:", peer_)
	}
}
