package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Info  string = ">>> MongoDB Ping <<<\n"
	Usage string = "\nUsage of %s:\n"

	Uri      string
	UriName  string = "uri"
	UriValue string = ""
	UriUsage string = "MongoDB connection string URI format"
	UriError string = "URI missing or malformed"

	Tls      bool
	TlsName  string = "tls"
	TlsValue bool   = false
	TlsUsage string = "Enable TLS"

	Ca      string
	CaName  string = "ca"
	CaValue string = ""
	CaUsage string = "Specify TLS CA file"
	CaError string = "Invalid CA file"

	Timeout      int
	TimeoutName  string = "timeout"
	TimeoutValue int    = 10
	TimeoutUsage string = "Connection timeout"

	ConnErr string = "!!! Connection failed !!!"
	ConnOk  string = "*** Successful ***"
)

func init() {
	fmt.Println(Info)

	flag.StringVar(&Uri, UriName, UriValue, UriUsage)
	flag.BoolVar(&Tls, TlsName, TlsValue, TlsUsage)
	flag.StringVar(&Ca, CaName, CaValue, CaUsage)
	flag.IntVar(&Timeout, TimeoutName, TimeoutValue, TimeoutUsage)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), Usage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NFlag() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	if len(Uri) < 1 {
		fmt.Println(UriError)
		flag.Usage()
		os.Exit(1)
	}

	if Tls && len(Ca) < 1 {
		fmt.Println(CaError)
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(Timeout),
	)
	defer cancel()

	tls := new(tls.Config)
	if Tls {
		certs, err := ioutil.ReadFile(Ca)
		if err != nil {
			fmt.Println(CaError)
			fmt.Println(err)
			os.Exit(1)
		}

		tls.RootCAs = x509.NewCertPool()
		if ok := tls.RootCAs.AppendCertsFromPEM(certs); !ok {
			fmt.Println(CaError)
			os.Exit(1)
		}
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(Uri).SetTLSConfig(tls))
	if err != nil {
		fmt.Println(ConnErr)
		fmt.Println(err)
		os.Exit(1)
	}

	if err := client.Connect(ctx); err != nil {
		fmt.Println(ConnErr)
		fmt.Println(err)
		os.Exit(1)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		fmt.Println(ConnErr)
		fmt.Println(err)
	} else {
		fmt.Println(ConnOk)
		fmt.Println(Uri)
	}
}
