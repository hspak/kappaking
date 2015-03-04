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

	// I'm not sure if this has any effect
	for _, stream := range LiveStreams {
		con.Part("#" + stream)
	}

	go func() {
		for action := range streamList {
			stream := strings.ToLower(action.Channel)
			if action.Join {
				// better slow than kicked
				time.Sleep(time.Second * 2)
				log.Println("Bot: joining", stream)
				con.Join("#" + stream)
			} else {
				log.Println("Bot: parting", stream)
				con.Part("#" + stream)
			}

			// required, either I crash or get kicked without it
			time.Sleep(time.Second)
		}
	}()

	// this buffer size is arbitrary
	kappaCounter := make(chan KappaData, 128)
	KPM = make(map[string]int)
	MaxKPM = make(map[string]int)
	TotalKappa = make(map[string]int)

	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		count := strings.Count(e.Message(), "Kappa")
		if count == 0 {
			return
		}

		name := strings.ToLower(e.Arguments[0][1:])
		if _, ok := KPM[name]; !ok {
			KPM[name] = 0

			// TODO: these might need special logic to pull from db
			MaxKPM[name], TotalKappa[name], err = grabKappa(db, name)
			if err != nil {
				log.Println(err)
			}
		}

		// fmt.Println("  Kappa add", name, "=>", count)
		kappaCounter <- KappaData{Name: name, Count: count}
		TotalKappa[name] += 1

		// subtract counts after minute
		go func() {
			time.Sleep(time.Minute)
			kappaCounter <- KappaData{Name: name, Count: -count}
			// fmt.Println("  Kappa sub:", name, " =>", KPM[name])
		}()
	})

	// use channels to avoid data hazard?
	go func() {
		for data := range kappaCounter {
			// TODO: keep track of max KPM
			KPM[data.Name] += data.Count
			if KPM[data.Name] > MaxKPM[data.Name] {
				MaxKPM[data.Name] = KPM[data.Name]
			}
			// fmt.Println("  Kappa update:", data.Name, " =>", KPM[data.Name])
		}
	}()

	con.Loop()
}
