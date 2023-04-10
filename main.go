package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/shadiestgoat/log"
)

var TWITCH_ID = ""

func main() {
	InitConfig()

	logCB := []log.LogCB{log.NewLoggerFile("log"), log.NewLoggerPrint()}
	if DEBUG_WEBHOOK != "" {
		logCB = append(logCB, log.NewLoggerDiscordWebhook(DEBUG_PREFIX, DEBUG_WEBHOOK))
	}
	log.Init(logCB...)
	log.Warn("This is a test run - any future warnings or errors will appear like this!")

	TWITCH_ID = GetOwnID()
	var tokenResp *struct {
		Token string `json:"token"`
	}
	DonFetch(`GET`, `/discordToken`, nil, &tokenResp)
	if tokenResp == nil {
		panic("Couldn't fetch discord token :(")
	}
	DISCORD_TOKEN = tokenResp.Token
	
	go InitDonations()

	fmt.Printf("Twitch ID %s\n", TWITCH_ID)
	
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<- stop
}
