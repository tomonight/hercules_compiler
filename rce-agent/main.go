package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"hercules/rce-agent/cmd"
	"hercules/rce-agent/rce"
	"log"
	"os"
)

var (
	showVersion = flag.Bool(
		"version", false,
		"Print version information.",
	)
	listenAddress = flag.String(
		"listen-address", ":5051",
		"Address to listen on.",
	)
	tlsCertFile = flag.String(
		"tls-cert-file", "",
		"Path to TLS certificate file.",
	)
	tlsKeyFile = flag.String(
		"tls-key-file", "",
		"Path to TLS key file.",
	)
	tlsRootCaFile = flag.String(
		"tls-rootca-file", "",
		"Path to TLS certificate root ca file.",
	)
)

const appVersion = "1.0.0"
const copyrightInfo = `znodelet starting, copyright enmotech.com 2018`

func main() {

	fmt.Println(copyrightInfo)

	flag.Parse()

	if *showVersion {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	var err error
	var tlsConfig *tls.Config

	if len(*tlsCertFile) > 0 && len(*tlsKeyFile) > 0 && len(*tlsRootCaFile) > 0 {

		tlsFiles := rce.TLSFiles{
			RootCert:   *tlsRootCaFile,
			ClientCert: *tlsCertFile,
			ClientKey:  *tlsKeyFile,
		}
		tlsConfig, err = tlsFiles.TLSConfig()
		if err != nil {
			log.Fatal(err)
		}
	}
	var commands cmd.Runnable
	var spec = cmd.Spec{Name: "exec_shell", Exec: []string{"/bin/bash", "-c"}}

	commands = append(commands, spec)

	s := rce.NewServer(*listenAddress, tlsConfig, commands)

	err = s.StartServer()
	if err != nil {
		log.Fatal(err)
	}
	defer s.StopServer()
}
