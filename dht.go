package main

import (
	"context"
	"fmt"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

func initDHT(host host.Host, rendezvous string) <-chan peer.AddrInfo {
	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	ctx := context.Background()
	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		panic(err)
	}

	fmt.Println("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}
	// We use a rendezvous point "meet me here" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	fmt.Println("Announcing ourselves...")
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, rendezvous)
	fmt.Println("Successfully announced!")

	// Now, look for others who have announced
	// This is like your friend telling you the location to meet you.
	fmt.Println("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(ctx, rendezvous)
	if err != nil {
		panic(err)
	}

	return peerChan
}
