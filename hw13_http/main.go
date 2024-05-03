package main

import (
	"flag"
	"os"

	client "github.com/aasdhajkshd/home_work_basic/hw13_http/client"
	data "github.com/aasdhajkshd/home_work_basic/hw13_http/data"
	server "github.com/aasdhajkshd/home_work_basic/hw13_http/server"
)

func main() {
	var (
		mode               string
		printJSON, verbose bool
	)

	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.StringVar(&mode, "mode", "server", "Specify 'client' or 'server' mode")
	flag.BoolVar(&printJSON, "print", false, "print the content of test JSON file (\"data/data.json\")")
	flag.Parse()

	if verbose {
		for i, j := range flag.Args() {
			print("parsed arguments:", i, j)
		}
	}

	if printJSON {
		data.PrintData("data/data.json")
		os.Exit(0)
	}

	switch mode {
	case "client":
		client.RunClient()
	default:
		print("Starting 'server'... by default\n")
		server.RunServer()
	}
}
