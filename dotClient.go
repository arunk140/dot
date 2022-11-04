package main

import (
	"log"
	"strings"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
)

type dotClient struct {
	node host.Host
}

func (context *dotClient) makeNode() host.Host {
	node, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}
	context.node = node
	return node
}

func (context *dotClient) getNodeAddress() string {
	addressStr := make([]string, 0)
	for _, address := range context.node.Addrs() {
		addressStr = append(addressStr, address.String())
	}
	return strings.Join(addressStr, ",")
}
