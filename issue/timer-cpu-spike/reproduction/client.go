package main

import (
	"crypto/tls"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// main creates a Ticker which triggers a TLS dial from a goroutine each tick
func main() {
	tick := time.NewTicker(time.Second * 10)
	done := make(chan bool)
	go scheduler(tick, done)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	done <- true
}

// scheduler triggers a TLS dial each tick
func scheduler(tick *time.Ticker, done chan bool) {
	for {
		select {
		case <-tick.C:
			connect()
		case <-done:
			return
		}
	}
}

func connect() {
	log.Println("connecting...")

	// Triggers initial calls to SetWriteDeadline
	d := net.Dialer{
		Timeout: time.Second * 30,
	}

	// Requires an address which gives fallbacks, i.e sysDialer.dialParallel() is used
	c, err := d.Dial("tcp", "[::]:10000")
	if err != nil {
		log.Fatal("Failed to dial: %v", err)
		return
	}

	// Setup a TLS client using existing transport
	cert, err := tls.LoadX509KeyPair("test.crt", "test.key")
	if err != nil {
		log.Printf("Failed to load keypair: %v", err)
		return
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true,
		Certificates: []tls.Certificate{cert}}
	c = tls.Client(c, tlsConfig)

	// Close in TLS calls closeNotify() which calls SetWriteDeadline
	defer c.Close()

	// Trigger TLS handshake and send data
	n, err := c.Write([]byte("PING\n"))
	if err != nil {
		log.Println(n, err)
		return
	}

	// Read response
	buf := make([]byte, 100)
	n, err = c.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}
}
