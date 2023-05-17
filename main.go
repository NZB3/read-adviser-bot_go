package main

import (
	"flag"
	"fmt"
	"log"
	"read-adviser-bot/clients/telegram"
)

func main() {
	var host string
	var token string
	// host falg
	flag.StringVar(
		&host,
		"host",
		"",
		"host for access to the telegram bot",
	)
	//token falag
	flag.StringVar(
		&token,
		"token",
		"",
		"token for access to the telegram bot `token`",
	)

	flag.Parse()

	if host == "" {
		log.Fatal("Missed host")
	}
	if token == "" {
		log.Fatal("Missed token")
	}

	tgClient := telegram.New(host, token)
	fmt.Print(tgClient)
	// fetcher = fetcher.New(tgClien)

	// processor = processor.New(tgClient)

	// consumer.Start(fetcher,processor)
}
