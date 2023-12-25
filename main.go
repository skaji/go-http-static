package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	apachelog "github.com/lestrrat-go/apache-logformat/v2"
)

var version = "dev"

func main() {
	port := flag.Int("port", 5001, "port number")
	sslCrt := flag.String("ssl-crt", "", "ssl crt file")
	sslKey := flag.String("ssl-key", "", "ssl key file")
	showVersion := flag.Bool("v", false, "show version and exit")
	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "too many arguments")
		os.Exit(1)
	}

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if (*sslCrt == "") != (*sslKey == "") {
		fmt.Fprintln(os.Stderr, "must specify both -ssl-crt and -ssl-key")
		os.Exit(1)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: apachelog.CombinedLog.Wrap(http.FileServer(http.Dir(".")), os.Stdout),
	}
	var err error
	if *sslCrt != "" && *sslKey != "" {
		fmt.Printf("Listen https://localhost:%d\n", *port)
		err = server.ListenAndServeTLS(*sslCrt, *sslKey)
	} else {
		fmt.Printf("Listen http://localhost:%d\n", *port)
		err = server.ListenAndServe()
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
