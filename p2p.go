package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"strings"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

const (
	fileProto protocol.ID = "/file/1.0.0"
)

type p2pClient struct {
	node host.Host
}

func (p *p2pClient) init(receiver bool) (host.Host, error) {
	r := rand.Reader
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	node, err := libp2p.New(
		libp2p.Identity(prvKey),
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// if receiver {
	node.SetStreamHandler(fileProto, handleStream("file"))
	// }
	p.node = node
	node.Addrs()
	return p.node, nil
}

func (p *p2pClient) connectToPeer(peerId peer.ID, addr peer.AddrInfo) (network.Stream, error) {
	ctx := context.Background()
	if err := p.node.Connect(ctx, addr); err != nil {
		log.Println(err)
		return nil, err
	}
	fileStream, err := p.node.NewStream(ctx, peerId, fileProto)
	return fileStream, err
}

func (p *p2pClient) getNodeID() string {
	return p.node.ID().Pretty()
}

func (p *p2pClient) getNodeAddrs() string {
	str := make([]string, 0)
	for _, address := range p.node.Addrs() {
		str = append(str, address.String())
	}
	return strings.Join(str, ",")
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')
		if str == "" {
			return
		}
		if str != "\n" {
			fmt.Printf(str)
		}
	}
}

func sendData(s network.Stream, data []byte) (int, error) {
	return s.Write(data)
}

func closeStream(s network.Stream) error {
	return s.Close()
}

func handleStream(typ string) func(s network.Stream) {
	return func(s network.Stream) {
		log.Println("Got a new stream!")
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		go readData(rw)
	}
}
