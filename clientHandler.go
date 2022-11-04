package main

import (
	"crypto/tls"
	"log"
	"strings"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func (p *p2pClient) responseHandler(conn *tls.Conn, response string) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	resParts := strings.Split(response, "|")
	switch resParts[0] {
	case "FOUND":
		if len(resParts) == 4 {
			log.Println("Found piece at", resParts[1])
		} else {
			log.Println("Invalid response from server")
		}
		log.Println(response)
		conn.Write([]byte("connect-request " + resParts[1] + " " + resParts[3]))
	case "CONNECTREQ":
		if len(resParts) == 3 {
			log.Println("Received connection request from", resParts[1])
		} else {
			log.Println("Invalid response from server")
		}
		p.init(false)
		peerId := p.getNodeID()
		addrs := p.getNodeAddrs()
		// log.Panicln(res)
		conn.Write([]byte("connect-response " + resParts[1] + " " + resParts[2] + " " + peerId + " " + addrs))
	case "CONNECTRES":
		// fileHash := resParts[1]
		// peerId := resParts[2]
		// peerAddrs := resParts[3]
		if len(resParts) == 4 {
			log.Println("Received connection response from", resParts[2])
		} else {
			log.Println("Invalid response from server")
		}

		p.init(true)
		// peerId := peer.ID("/ip4/192.168.2.222/tcp/41577" + "/p2p/" + resParts[2])
		pId, err := peer.Decode(resParts[2])
		if err != nil {
			log.Panicln(err)
		}
		// log.Println(len(resParts[2]))
		// log.Panicln(pId)
		log.Println("Connecting to peer", pId)

		peerConnAddrs := strings.Split(resParts[3], ",")
		// []string to multiaddr.Multiaddr
		var peerAddrs []multiaddr.Multiaddr
		for _, addr := range peerConnAddrs {
			maddr, err := multiaddr.NewMultiaddr(addr)
			if err != nil {
				log.Println(err)
			}
			peerAddrs = append(peerAddrs, maddr)
		}

		peerAddrInfo := peer.AddrInfo{
			ID:    pId,
			Addrs: peerAddrs,
		}

		pC, err := p.connectToPeer(pId, peerAddrInfo)
		if err != nil {
			log.Panicln(err)
		}

		pC.Write([]byte("hello"))
		// stm.Write([]byte("Hello from client"))
		log.Println(p.node.Peerstore().Peers())
		log.Println(p.node.ID().Pretty())
		log.Println("Received connection response from ", resParts)
	default:
		log.Println(response)
	}
}
