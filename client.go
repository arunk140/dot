package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
)

func waitForInput(conn *tls.Conn) {
	for {
		var input string
		bufioReader := bufio.NewReader(os.Stdin)
		input, err := bufioReader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		if input == "disconnect\n" {
			conn.Close()
			return
		}

		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Fatalln("error writing to connection")
			conn.Close()
			log.Fatalln(err)
		}
	}
}

func waitForResponse(conn *tls.Conn) {
	p := p2pClient{node: nil}
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			log.Fatalln("ERROR: Connection closed")
		}
		p.responseHandler(conn, string(buf[:n]))
	}
}

func runTCPClient(asAdmin bool) {
	certPem := "server.pem"
	roots := x509.NewCertPool()
	caCert, err := os.ReadFile(certPem)
	if err != nil {
		log.Fatalln(err)
		return
	}
	roots.AppendCertsFromPEM(caCert)

	if err != nil {
		log.Fatalf("client: loadkeys: %s\n", err)
		return
	}
	config := tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS13, MaxVersion: tls.VersionTLS13}
	conn, err := tls.Dial("tcp", "localhost:3333", &config)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println("Connection established")
	log.Println(conn.RemoteAddr().String())
	defer func() {
		defer conn.Close()
		log.Println("Connection closed")
	}()

	go waitForInput(conn)
	waitForResponse(conn)
}
