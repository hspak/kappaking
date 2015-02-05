package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/thoj/go-ircevent"
)

func getPassword() string {
	pass, err := ioutil.ReadFile("password")
	if err != nil {
		log.Fatal("could not open password file")
	}
	return string(pass)
}

func launchBot() {
	con := irc.IRC("kappakingbot", "kappakingbot")
	con.Password = getPassword()
	err := con.Connect("irc.twitch.tv:6667")
	if err != nil {
		log.Fatal(err)
	}

	// on successful connect
	con.AddCallback("001", func(e *irc.Event) {
		con.Join("#ezjijy")
	})
	con.AddCallback("JOIN", func(e *irc.Event) {
		con.Privmsg("#ezjijy", "HI!")
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		// process kappa's here
		fmt.Println("msg", e.Message())
	})
	con.Loop()
}
