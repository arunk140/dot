package main

import (
	"crypto/tls"
	"log"
	"net"
)

var serverCtx dotServer

func handleConnection(conn net.Conn, serverCtx *dotServer) {
	defer func() {
		serverCtx.kickConnection(conn)
		serverCtx.printActiveConnections()
	}()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() != "EOF" {
				log.Println("ERROR:", err)
			}
			return
		}
		_, err = serverCtx.handleMessage(string(buf[:n]), conn)
		if err != nil {
			log.Println("ERROR:", err)
			return
		}
	}
}

func runTCPServer() {
	certPvt := "server.key"
	certPem := "server.pem"
	cert, err := tls.LoadX509KeyPair(certPem, certPvt)
	if err != nil {
		log.Fatalf("server: loadkeys: %s\n", err)
		return
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS13, MaxVersion: tls.VersionTLS13}
	log.Println("Listening on port 3333")
	ln, err := tls.Listen("tcp", ":3333", &config)
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer ln.Close()

	serverCtx = dotServer{ln: ln, connections: make(map[string]*dotConnection)}

	for {
		conn, err := ln.Accept()
		if ln == nil || err != nil {
			log.Println(err)
			log.Fatalln("Server Killed")
			return
		}
		serverCtx.connections[conn.RemoteAddr().String()] = &dotConnection{conn: conn, isAdmin: false, pieceList: make(map[string]*dotPiece)}
		serverCtx.printActiveConnections()
		go handleConnection(conn, &serverCtx)
	}
}
