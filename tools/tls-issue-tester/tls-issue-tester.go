package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gomodule/redigo/redis"
)

func redis_connect() {
	log.Println("Connecting to redis..")
	tlsConfig, err := CreateClientTLSConfig("/tls-data/curl.crt", "/tls-data/curl.key", "/tls-data/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	connectionTimeouts, err := time.ParseDuration("1s")
	if err != nil {
		log.Fatal(err)
	}
	options := []redis.DialOption{
		redis.DialConnectTimeout(connectionTimeouts),
		redis.DialReadTimeout(connectionTimeouts),
		redis.DialWriteTimeout(connectionTimeouts),
		redis.DialTLSConfig(tlsConfig),
	}

	c, err := redis.DialURL("rediss://localhost:6379", options...)
	if err != nil {
		log.Print(err)
		return
	}
	defer c.Close()

	_, err = c.Do("PING")
	if err != nil {
		log.Print(err)
		return
	}
	log.Println("Disconnecting")
}

// CreateClientTLSConfig verifies configured files and return a prepared tls.Config
func CreateClientTLSConfig(ClientCertFile, ClientKeyFile, CaCertFile string) (*tls.Config, error) {
	tlsConfig := tls.Config{
		InsecureSkipVerify: false,
	}

	cert, err := LoadKeyPair(ClientCertFile, ClientKeyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig.Certificates = []tls.Certificate{*cert}
	certificates, err := LoadCAFile(CaCertFile)
	if err != nil {
		return nil, err
	}
	tlsConfig.RootCAs = certificates

	return &tlsConfig, nil
}

// The files must contain PEM encoded data.
func LoadKeyPair(certFile, keyFile string) (*tls.Certificate, error) {
	log.Printf("Load key pair: %s %s", certFile, keyFile)
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// LoadCAFile reads and parses CA certificates from a file into a pool.
// The file must contain PEM encoded data.
func LoadCAFile(caFile string) (*x509.CertPool, error) {
	log.Printf("Load CA cert file: %s", caFile)
	pemCerts, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCerts)
	return pool, nil
}

func MetricsScrape(w http.ResponseWriter, req *http.Request) {
	log.Println("Metrics scraping by connecting to redis")
	redis_connect()
	io.WriteString(w, fmt.Sprintf("current_time_seconds %d", time.Now().Unix()))
}

func main() {
	log.Printf("Using Go version: %s\n", runtime.Version())

	http.HandleFunc("/metrics", MetricsScrape)
	err := http.ListenAndServeTLS(":9121", "/tls-data/exporter-s.crt", "/tls-data/exporter-s.key", nil)
	//err := http.ListenAndServe(":9121", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
