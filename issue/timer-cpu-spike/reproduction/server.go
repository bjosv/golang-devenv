package main

import (
	"crypto/tls"
	"log"
	"net"
)

func main() {
	cert, err := tls.LoadX509KeyPair("test.crt", "test.key")
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	log.Println("server: start listen")
	l, err := tls.Listen("tcp", ":10000", config)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			defer c.Close()
			buf := make([]byte, 512)
			for {
				n, err := conn.Read(buf)
				if err != nil {
					log.Printf("server: read error: %s", err)
					break
				}

				log.Printf("server: echo %q\n", string(buf[:n]))
				n, err = conn.Write(buf[:n])
				if err != nil {
					log.Printf("server: write error: %s", err)
					break
				}
			}
		}(conn)
	}
}
