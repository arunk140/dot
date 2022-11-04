package main

import (
	"fmt"
	"log"
	"net"
)

type dotPiece struct {
	pieceHash string
	pieceSize string
}

type dotConnection struct {
	conn      net.Conn
	pieceList map[string]*dotPiece
	isAdmin   bool
}

type dotServer struct {
	ln          net.Listener
	connections map[string]*dotConnection
}

func (context *dotServer) kickConnection(conn net.Conn) {
	// log.Println("Connection closed")
	conn.Close()
	delete(context.connections, conn.RemoteAddr().String())
}

func (context *dotServer) printActiveConnections() {
	for _, dc := range context.connections {
		dc.conn.Write([]byte(fmt.Sprintf("Active connections: %d", len(context.connections))))
	}
	log.Println("Active Connections: ", len(context.connections))
}

func (context *dotServer) kickAll() {
	for _, dc := range context.connections {
		if dc.isAdmin {
			continue
		}
		context.kickConnection(dc.conn)
	}
}

func (context *dotServer) killServer() {
	context.ln.Close()
}

func (context *dotServer) broadcast(msg string) {
	for _, dc := range context.connections {
		_, err := dc.conn.Write([]byte(msg))
		if err != nil {
			log.Println("Error sending message to client")
			log.Println(err)
		}
	}
}
