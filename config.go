package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type confItem struct {
	Res      *string
	Default  string
	Required bool
}

var (
	STREAM_ELEMENTS_JWT = ""
	DONATION_TOKEN = ""
	DONATION_HOST = ""
	DISCORD_TOKEN = ""
)

func InitConfig() {
	godotenv.Load(".env")

	var confMap = map[string]confItem{
		"STREAM_ELEMENTS_JWT": {
			Res:      &STREAM_ELEMENTS_JWT,
			Required: true,
		},
		"DONATION_TOKEN": {
			Res:     &DONATION_TOKEN,
			Required: true,
		},
		"DONATION_HOST": {
			Res:      &DONATION_HOST,
			Required: true,
		},
	}

	for name, opt := range confMap {
		item := os.Getenv(name)

		if item == "" {
			if opt.Required {
				panic(fmt.Sprintf("'%v' is a needed variable, but is not present! Please read the README.md file for more info.", name))
			}
			item = opt.Default
		}

		*opt.Res = item
	}
}
