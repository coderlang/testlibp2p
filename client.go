package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"log"
	"os"
)

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

	//if err := host.Connect(context.Background(), *nodeAPeerInfo); err != nil {
	//	log.Fatal(err)
	//}

	fmt.Printf("Client Node ID: %+v\n", host.ID())

	peerChan := initDHT(host, "client")
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
