package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func respond(conn net.Conn, msg string) (n int, err error) {
	return conn.Write([]byte(msg))
}

func respondSuccess(conn net.Conn, msg string) (n int, err error) {
	return respond(conn, msg)
}

func respondFail(conn net.Conn, msg string) (n int, err error) {
	return respond(conn, msg)
}

func (context *dotServer) queryPieceList(conn net.Conn, pieceHash string) (n int, err error) {
	for _, dc := range context.connections {
		if _, ok := dc.pieceList[pieceHash]; ok {
			return respondSuccess(conn, "FOUND|"+dc.conn.RemoteAddr().String()+"|"+dc.pieceList[pieceHash].pieceSize+"|"+dc.pieceList[pieceHash].pieceHash)
		}
	}
	return respondFail(conn, "Piece not found")
}

func (context *dotServer) addPiece(conn net.Conn, pieceHash string) (n int, err error) {
	if _, ok := context.connections[conn.RemoteAddr().String()].pieceList[pieceHash]; ok {
		return respondFail(conn, "Piece already exists")
	}
	context.connections[conn.RemoteAddr().String()].pieceList[pieceHash] = &dotPiece{pieceHash: pieceHash, pieceSize: "0"}
	return respondSuccess(conn, "Piece Added")
}

func (context *dotServer) handleConnectionRequest(conn net.Conn, args []string) (n int, err error) {
	if len(args) < 3 {
		return respondFail(conn, "Not enough arguments")
	}
	connectTo := args[1]
	if connectTo == conn.RemoteAddr().String() {
		return respondFail(conn, "Cannot connect to self")
	}
	if _, ok := context.connections[connectTo]; !ok {
		return respondFail(conn, "Connection does not exist")
	}
	if _, ok := context.connections[connectTo]; ok {
		pieceHash := args[2]
		if _, ok := context.connections[connectTo].pieceList[pieceHash]; !ok {
			return respondFail(conn, "Piece does not exist")
		}
		respond(context.connections[connectTo].conn, "CONNECTREQ|"+conn.RemoteAddr().String()+"|"+pieceHash)
	}
	return respondSuccess(conn, "Sent connection request")
}

func (context *dotServer) handleConnectionResponse(conn net.Conn, args []string) (n int, err error) {
	// args[1] = remote address
	// args[2] = piece hash
	// args[3] = p2p node id
	// args[4] = p2p node address
	if len(args) < 5 {
		return respondFail(conn, "Not enough arguments")
	}
	remoteAddr := args[1]
	pieceHash := args[2]
	p2pNodeId := args[3]
	p2pNodeAddr := args[4]
	if _, ok := context.connections[remoteAddr]; !ok {
		return respondFail(conn, "Connection does not exist")
	}
	context.connections[remoteAddr].conn.Write([]byte("CONNECTRES|" + pieceHash + "|" + p2pNodeId + "|" + p2pNodeAddr))
	return respondSuccess(conn, "Sent connection response")
}

func (context *dotServer) handleMessage(msg string, conn net.Conn) (n int, err error) {
	msg = strings.TrimSuffix(msg, "\n")
	args := strings.Split(msg, " ")
	restOfMsg := strings.Join(args[1:], " ")

	cmd := strings.ToLower(args[0])
	log.Println("Message received: ", msg)
	switch cmd {
	case "ping":
		return respondSuccess(conn, "PONG")
	case "list":
		listOfConnections := ""
		for k := range context.connections {
			isMe := ""
			if k == conn.RemoteAddr().String() {
				isMe = " (me)"
			}
			listOfConnections += k + isMe + "\n"
		}
		return respondSuccess(conn, "Connections ("+fmt.Sprintf("%d", len(context.connections))+"): \n"+listOfConnections)
	case "piece":
		return context.addPiece(conn, restOfMsg)
	case "query":
		pieceHash := restOfMsg
		return context.queryPieceList(conn, pieceHash)
	case "escallate":
		if len(args) < 2 {
			return respondFail(conn, "Not enough arguments")
		}
		if args[1] != "4c78ubnyqrt08b234yqc3!#!@#$@#%ekldnfgioer34@#$%#$@5fdvkmlfg34CV$G$" {
			return respondFail(conn, "Incorrect password")
		}
		if context.connections[conn.RemoteAddr().String()].isAdmin {
			return respondFail(conn, "Already an admin")
		}
		context.connections[conn.RemoteAddr().String()].isAdmin = true
		return respondSuccess(conn, "Escallated to admin")
	case "broadcast":
		log.Println("Broadcasting message: " + restOfMsg)
		context.broadcast(restOfMsg)
		return respondSuccess(conn, "Broadcasted message")
	case "kickall":
		if context.connections[conn.RemoteAddr().String()].isAdmin {
			context.kickAll()
			return respondSuccess(conn, "Kicked all")
		}
		return respondFail(conn, "Not an admin")
	case "killserver":
		if context.connections[conn.RemoteAddr().String()].isAdmin {
			context.killServer()
			return respondSuccess(conn, "Killed server")
		}
		return respondFail(conn, "Not an admin")
	case "connect-request":
		return context.handleConnectionRequest(conn, args)
	case "connect-response":
		return context.handleConnectionResponse(conn, args)
	default:
		return respondFail(conn, "Unknown command")
	}
}
