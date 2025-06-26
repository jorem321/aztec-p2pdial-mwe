package main

import (
	"context"
	"crypto/rand"
	"flag"
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
	AZTEC_STATUS_PROTOCOL = "/aztec/req/status/0.1.0"
)

func main() {
	// Parse target peer from command line
	targetF := flag.String("d", "", "target peer to dial")
	flag.Parse()

	if *targetF == "" {
		log.Fatal("Please provide a target peer to dial with -d")
	}

	// Set log level
	golog.SetAllLoggers(golog.LevelWarn) // Reduce noise

	// Create minimal host
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	host, err := libp2p.New(
		libp2p.Identity(priv),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.NoListenAddrs, // Don't listen, just dial out
		libp2p.DisableRelay(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer host.Close()

	// Parse target multiaddr
	maddr, err := ma.NewMultiaddr(*targetF)
	if err != nil {
		log.Fatal("Invalid multiaddr:", err)
	}

	// Extract peer info
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatal("Invalid peer address:", err)
	}

	// Add to peerstore
	host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	// Open stream and send status request
	log.Println("Requesting status from Aztec node...")
	s, err := host.NewStream(context.Background(), info.ID, AZTEC_STATUS_PROTOCOL)
	if err != nil {
		log.Fatal("Failed to open stream:", err)
	}
	defer s.Close()

	// Send status request
	_, err = s.Write([]byte("status"))
	if err != nil {
		log.Fatal("Failed to send request:", err)
	}

	// Signal end of request
	s.CloseWrite()

	// Read response
	response, err := io.ReadAll(s)
	if err != nil {
		log.Fatal("Failed to read response:", err)
	}

	// Display results
	log.Printf("Status response (%d bytes): %q", len(response), response)
	log.Printf("Status response (hex): %x", response)
}
