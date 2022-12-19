package main

import (
	"fmt"
	"os"
	"os/signal"
)

var TWITCH_ID = ""

func main() {
	InitConfig()
	TWITCH_ID = GetOwnID()
	var tokenResp *struct {
		Token string `json:"token"`
	}
	DonFetch(`GET`, `/discordToken`, nil, &tokenResp)
	if tokenResp == nil {
		panic("Couldn't fetch discord token :(")
	}
	DISCORD_TOKEN = tokenResp.Token
	go InitDonation()

	fmt.Printf("Twitch ID %s\n", TWITCH_ID)
	
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<- stop
}
