package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/thoj/go-ircevent"
)

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

	go func() {
		for name := range streamList {
			stream := strings.ToLower(name)
			fmt.Println("joining", stream)
			con.Join("#" + stream)
			time.Sleep(time.Second)
		}
	}()

	// on successful connect
	// con.AddCallback("001", func(e *irc.Event) {
	// con.Join("#ezjijy")
	// con.Join("#kappaking")
	// })
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		fmt.Println("msg", strings.Count(e.Message(), "Kappa"))
	})
	con.Loop()
}
