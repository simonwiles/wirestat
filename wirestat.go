package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	portOpt := flag.Uint("port", 8930, "Port to run the server on")
	rulesPathOpt := flag.String("rules", "/etc/wirestat/rules.txt", "Path to the file containing rules")

	flag.Parse()
	args := flag.Args()
	rulesPath, err := filepath.Abs(*rulesPathOpt)
	if err != nil {
		panic(err)
	}

	if len(args) > 0 && args[0] == "systemd" {
		config := GenerateSystemdConfig(*portOpt, rulesPath)
		fmt.Println(config)
		os.Exit(0)
	}

	_, err = os.Stat(rulesPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("Startup failed, error when checking rules file at %s", rulesPath))
		os.Exit(1)
	}

	rules, err := parseRuleFile(rulesPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("Startup failed, error when parsing rules file: %s", err.Error()))
		os.Exit(1)
	}

	responseBuilder := NewResponseBuilder(rules)

	fmt.Println(fmt.Sprintf("Starting server, listening on: http://0.0.0.0:%d", *portOpt))
	startServer(responseBuilder, *portOpt)
}

func dd(data ...interface{}) {
	for _, item := range data {
		fmt.Printf("%+v\n", item)
	}
	os.Exit(1)
}
