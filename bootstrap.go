package main

import (
	"context"
	"encoding/json"
	"fmt"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	"log"
	"os"

	"github.com/libp2p/go-libp2p"
)

func main() {
	ctx := context.Background()

	// 创建一个 libp2p 主机
	host, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}

	// 创建 DHT 实例并配置为 ModeServer
	dhtInstance, err := dht.New(ctx, host, dht.Mode(dht.ModeServer))
	if err != nil {
		log.Fatal(err)
	}
	err = dhtInstance.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}

	peerInfo := peer.AddrInfo{
		ID:    host.ID(),
		Addrs: host.Addrs(),
	}
	fmt.Printf("Bootstrap Node ID: %+v\n", peerInfo.ID)
	js, err := json.Marshal(peerInfo)
	if err != nil {
		panic(err)
	}
	saveJSONToFile(js, "bootstrap.json")
	select {}
}

func saveJSONToFile(jsonData []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
