package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/thoj/go-ircevent"
)

type KappaData struct {
	Name  string
	Count int
}

// TODO: should figure out a better way to share this data
var KPM map[string]int

func getPassword() string {
	pass, err := ioutil.ReadFile("password")
	if err != nil {
		log.Fatal("could not open password file")
	}
	return string(pass)
}

func launchBot(streamList chan string) {
	con := irc.IRC("kappakingbot", "kappakingbot")
	con.Password = getPassword()
	err := con.Connect("irc.twitch.tv:6667")
	if err != nil {
		log.Fatal(err)
	}

	con.Join("#ezjijy")
	// go func() {
	// for name := range streamList {
	// stream := strings.ToLower(name)
	// fmt.Println("joining", stream)
	// con.Join("#" + stream)

	// // required, either I crash or get kicked without it
	// time.Sleep(time.Second)
	// }
	// }()

	// this count is arbitrary
	kappaCounter := make(chan KappaData, 128)
	KPM = make(map[string]int)

	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		name := strings.Split(e.Host, ".")[0]
		count := strings.Count(e.Message(), "Kappa")

		if _, ok := KPM[name]; !ok {
			KPM[name] = 0
		}

		fmt.Println("Kappa from", name, "=", count)
		kappaCounter <- KappaData{Name: name, Count: count}

		// subtract counts after minute
		go func() {
			time.Sleep(time.Minute)
			kappaCounter <- KappaData{Name: name, Count: -count}
			fmt.Println("sub:", name, " =>", KPM[name])
		}()
	})

	// use channels to avoid data hazard?
	go func() {
		for data := range kappaCounter {
			time.Sleep(time.Second)
			KPM[data.Name] += data.Count
			fmt.Println("update:", data.Name, " =>", KPM[data.Name])
		}
	}()

	con.Loop()
}
