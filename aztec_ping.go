package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"

	golog "github.com/ipfs/go-log/v2"
	ma "github.com/multiformats/go-multiaddr"
)

const (
	AZTEC_PING_PROTOCOL = "/aztec/req/ping/0.1.0"
)

func main() {
	// Suppress libp2p logs for cleaner output
	golog.SetAllLoggers(golog.LevelError)

	// Parse command line arguments
	targetF := flag.String("d", "", "target peer to dial")
	flag.Parse()

	if *targetF == "" {
		log.Fatal("Please provide a target peer with -d")
	}

	// Create a libp2p host
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	host, err := libp2p.New(
		libp2p.Identity(priv),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.DisableRelay(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer host.Close()

	// Parse target peer multiaddr
	maddr, err := ma.NewMultiaddr(*targetF)
	if err != nil {
		log.Fatal(err)
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatal(err)
	}

	// Add peer to peerstore
	host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	// Open stream to target peer
	s, err := host.NewStream(context.Background(), info.ID, AZTEC_PING_PROTOCOL)
	if err != nil {
		log.Fatal("failed to open stream:", err)
	}
	defer s.Close()

	// Send ping request
	_, err = s.Write([]byte("ping"))
	if err != nil {
		log.Fatal("failed to send ping:", err)
	}

	// Close write side to signal end of request
	if err := s.CloseWrite(); err != nil {
		log.Fatal("failed to close write:", err)
	}

	// Read response
	response, err := io.ReadAll(s)
	if err != nil {
		log.Fatal("failed to read response:", err)
	}

	// Parse and display response
	if len(response) >= 6 && string(response[5:]) == "pong" {
		fmt.Println("✅ pong")
	} else {
		fmt.Printf("Response: %q\n", response)
	}
}
