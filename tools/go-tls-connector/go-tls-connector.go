package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetReportCaller(true)
	log.Infof("Using Go version: %s", runtime.Version())

	uri := flag.String("uri", getEnv("URI", "rediss://localhost:6380"), "")
	clientCertFile := "/redis-conf/curl.crt"
	clientKeyFile := "/redis-conf/curl.key"
	caCertFile := "/redis-conf/ca.crt"
	numberOfPings := flag.Int64("pings", getEnvInt64("PINGS", 10), "")

	tlsConfig, err := CreateClientTLSConfig(clientCertFile, clientKeyFile, caCertFile)
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

	for {
		log.Infof("Sending %d pings...", *numberOfPings)
		var wg sync.WaitGroup
		var i int64
		for i = 0; i < *numberOfPings; i++ {
			wg.Add(1)
			go ping(&wg, *uri, options)
		}
		wg.Wait()

		time.Sleep(1 * time.Second)
	}
}

func ping(wg *sync.WaitGroup, uri string, options []redis.DialOption) {
	defer wg.Done()

	c, err := redis.DialURL(uri, options...)
	if err != nil {
		log.Error(err)
		return
	}
	defer c.Close()

	_, err = c.Do("PING")
	if err != nil {
		log.Error(err)
		return
	}
}

func getEnv(key string, defaultVal string) string {
	if envVal, ok := os.LookupEnv(key); ok {
		return envVal
	}
	return defaultVal
}

func getEnvInt64(key string, defaultVal int64) int64 {
	if envVal, ok := os.LookupEnv(key); ok {
		envInt64, err := strconv.ParseInt(envVal, 10, 64)
		if err == nil {
			return envInt64
		}
	}
	return defaultVal
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

// GetConfigForClientFunc returns a function for tls.Config.GetConfigForClient
func GetConfigForClientFunc(certFile, keyFile, caCertFile string) func(*tls.ClientHelloInfo) (*tls.Config, error) {
	return func(*tls.ClientHelloInfo) (*tls.Config, error) {
		certificates, err := LoadCAFile(caCertFile)
		if err != nil {
			return nil, err
		}

		tlsConfig := tls.Config{
			ClientAuth:     tls.RequireAndVerifyClientCert,
			ClientCAs:      certificates,
			GetCertificate: GetServerCertificateFunc(certFile, keyFile),
		}
		return &tlsConfig, nil
	}
}

// GetServerCertificateFunc returns a function for tls.Config.GetCertificate
func GetServerCertificateFunc(certFile, keyFile string) func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		return LoadKeyPair(certFile, keyFile)
	}
}

// LoadKeyPair reads and parses a public/private key pair from a pair of files.
// The files must contain PEM encoded data.
func LoadKeyPair(certFile, keyFile string) (*tls.Certificate, error) {
	log.Debugf("Load key pair: %s %s", certFile, keyFile)
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// LoadCAFile reads and parses CA certificates from a file into a pool.
// The file must contain PEM encoded data.
func LoadCAFile(caFile string) (*x509.CertPool, error) {
	log.Debugf("Load CA cert file: %s", caFile)
	pemCerts, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCerts)
	return pool, nil
}
