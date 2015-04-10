package main

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/hspak/go-ircevent"
)

func getPassword() string {
	pass, err := ioutil.ReadFile("password")
	if err != nil {
		log.Fatal("could not open password file")
	}
	return string(pass)
}

func launchBot(db DB) {
	con := irc.IRC("kappakingbot", "kappakingbot")
	con.Password = getPassword()
	err := con.Connect("irc.twitch.tv:6667")
	if err != nil {
		log.Fatal(err)
	}

	// TODO: set all joinedChannels when disconnected
	joinedChannels := make(map[string]bool)
	go func() {
		for action := range streamList {
			stream := strings.ToLower(action.Channel)
			_, exist := joinedChannels[stream]
			// intentionally rejoining already joined channels
			// so that it can still join properly after disconnects
			if (action.Join && !exist) || (action.Join) {
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

	go func() {
		for {
			time.Sleep(time.Minute)
			for name, _ := range joinedChannels {
				Minutes[name] += 1
			}
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

		}()
	})

	go func() {
		for data := range kappaCounter {
			KPM[data.Name] += data.Count
			if KPM[data.Name] > MaxKPM[data.Name] {
				MaxKPM[data.Name] = KPM[data.Name]
				DateKPM[data.Name] = time.Now().UTC()
			}
		}
	}()
	con.Loop()
}
