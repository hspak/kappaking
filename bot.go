package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/hspak/go-ircevent"
)

// TODO: should probably create a struct with an interface for this data
var MaxKPM map[string]int
var KPM map[string]int
var TotalKappa map[string]int
var Minutes map[string]int
var DateKPM map[string]time.Time

func getPassword() string {
	pass, err := ioutil.ReadFile("password")
	if err != nil {
		log.Fatal("could not open password file")
	}
	return string(pass)
}

func launchBot(db *sql.DB, streamList chan *BotAction) {
	con := irc.IRC("kappakingbot", "kappakingbot")
	con.Password = getPassword()
	err := con.Connect("irc.twitch.tv:6667")
	if err != nil {
		log.Fatal(err)
	}

	joinedChannels := make(map[string]bool)
	go func() {
		for action := range streamList {
			stream := strings.ToLower(action.Channel)
			joined, exist := joinedChannels[stream]
			if (action.Join && !exist) || (action.Join && !joined) {
				joinedChannels[stream] = true

				// better slow than kicked
				time.Sleep(time.Second * 2)
				log.Println("Bot: joining", stream)
				con.Join("#" + stream)
			} else if !action.Join {
				joinedChannels[stream] = false
				log.Println("Bot: parting", stream)
				con.Part("#" + stream)
			}

			// required, either I crash or get kicked without it
			time.Sleep(time.Second)
		}
	}()

	// this buffer size is arbitrary
	kappaCounter := make(chan KappaData, 128)
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		count := strings.Count(e.Message(), "Kappa")
		if count == 0 {
			return
		}

		name := strings.ToLower(e.Arguments[0][1:])
		if _, ok := KPM[name]; !ok {
			KPM[name] = 0
		}

		kappaCounter <- KappaData{Name: name, Count: count}
		TotalKappa[name] += count

		// subtract counts after minute
		go func() {
			time.Sleep(time.Minute)
			kappaCounter <- KappaData{Name: name, Count: -count}

			// try to disconnect completely instead of trying to reconnect
			if !con.Connected() {
				log.Println("Bot: manually reconnecting")
				con.Disconnect()
				time.Sleep(time.Second)
				err := con.Connect("irc.twitch.tv:6667")
				if err != nil {
					log.Fatal(err)
				}
			}
		}()
	})

	// use channels to avoid data hazard?
	go func() {
		for data := range kappaCounter {
			KPM[data.Name] += data.Count
			if KPM[data.Name] > MaxKPM[data.Name] {
				kpmdate := time.Now().UTC()
				MaxKPM[data.Name] = KPM[data.Name]
				DateKPM[data.Name] = kpmdate
			}
		}
	}()

	con.Loop()
}
