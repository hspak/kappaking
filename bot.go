package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"go-ircevent"
)

// TODO: should figure out a better way to share this data
var KPM map[string]int

func getPassword() string {
	pass, err := ioutil.ReadFile("password")
	if err != nil {
		log.Fatal("could not open password file")
	}
	return string(pass)
}

func launchBot(streamList chan *BotAction) {
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

	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		count := strings.Count(e.Message(), "Kappa")
		if count == 0 {
			return
		}

		name := strings.ToLower(e.Arguments[0][1:])
		if _, ok := KPM[name]; !ok {
			KPM[name] = 0
		}

		// fmt.Println("  Kappa add", name, "=>", count)
		kappaCounter <- KappaData{Name: name, Count: count}

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
			fmt.Println("  Kappa update:", data.Name, " =>", KPM[data.Name])
		}
	}()

	con.Loop()
}
