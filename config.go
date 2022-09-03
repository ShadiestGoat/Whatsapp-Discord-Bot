package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type confItem struct {
	Res         *string
	Default     string
	PanicIfNone bool
}

var (
	TOKEN         = ""
	GUILD_CHANNEL = ""
	GUILD_ROLE    = ""
	CHAT_NAME	  = ""
)

func ConfigInit() {
	godotenv.Load(".env")

	var confMap = map[string]confItem{
		"TOKEN": {
			Res:         &TOKEN,
			PanicIfNone: true,
		},
		"GUILD_CHANNEL": {
			Res:         &GUILD_CHANNEL,
			PanicIfNone: true,
		},
		"GUILD_ROLE": {
			Res: &GUILD_ROLE,
			PanicIfNone: true,
		},
		"CHAT_NAME": {
			Res: &CHAT_NAME,
			PanicIfNone: true,
		},
	}

	for name, opt := range confMap {
		item := os.Getenv(name)

		if item == "" {
			if opt.PanicIfNone {
				panic(fmt.Sprintf("'%v' is a needed variable, but is not present! Please read the README.md file for more info.", name))
			}
			item = opt.Default
		}

		*opt.Res = item
	}

	fmt.Println("Config loaded!")
}
