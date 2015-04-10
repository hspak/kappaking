package main

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/hspak/go-ircevent"
)

type KappaData struct {
	Name  string
	Count int
}

type Bot struct {
	db             *DB
	joinedChannels map[string]bool
	conn           *irc.Connection
	kappaCounter   chan KappaData
	stop           bool
}

func NewBot(database *DB) *Bot {
	return &Bot{
		db:             database,
		joinedChannels: make(map[string]bool),
		conn:           nil,
		kappaCounter:   make(chan KappaData, 128),
		stop:           false}
}

func (b *Bot) connect(nick, user, server string) {
	b.conn = irc.IRC(nick, user)
	b.setPassword()
	err := b.conn.Connect(server)
	if err != nil {
		log.Fatal(err)
	}
}

func (b *Bot) setPassword() {
	pass, err := ioutil.ReadFile("password")
	if err != nil {
		log.Fatal("could not open password file")
	}
	b.conn.Password = string(pass)
}

func (b *Bot) joinChannels() {
	for action := range b.db.streamList {
		if b.stop {
			return
		}
		stream := strings.ToLower(action.Channel)
		_, exist := b.joinedChannels[stream]
		// intentionally rejoining already joined channels
		// so that it can still join properly after disconnects
		if (action.Join && !exist) || (action.Join) {
			b.joinedChannels[stream] = true

			// better slow than kicked
			time.Sleep(time.Second * 2)
			log.Println("Bot: joining", stream)
			b.conn.Join("#" + stream)
		} else if !action.Join {
			b.joinedChannels[stream] = false
			log.Println("Bot: parting", stream)
			b.conn.Part("#" + stream)
		}

		// required, either I crash or get kicked without it
		time.Sleep(time.Second)
	}
}

func (b *Bot) countMinutes() {
	for {
		if b.stop {
			return
		}
		time.Sleep(time.Minute)
		for name, _ := range b.joinedChannels {
			b.db.cache.Store.Minutes[name] += 1
		}
	}
}

func (b *Bot) updateCounts() {
	for data := range b.kappaCounter {
		if b.stop {
			return
		}
		b.db.cache.KPM[data.Name] += data.Count
		if b.db.cache.KPM[data.Name] > b.db.cache.Store.MaxKPM[data.Name] {
			b.db.cache.Store.MaxKPM[data.Name] = b.db.cache.KPM[data.Name]
			b.db.cache.Store.DateKPM[data.Name] = time.Now().UTC()
		}
	}
}

func (b *Bot) trackKappas() {
	// this buffer size is arbitrary
	b.conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		count := strings.Count(e.Message(), "Kappa")
		if count == 0 {
			return
		}

		name := strings.ToLower(e.Arguments[0][1:])
		if _, ok := b.db.cache.KPM[name]; !ok {
			b.db.cache.KPM[name] = 0
		}

		b.kappaCounter <- KappaData{Name: name, Count: count}
		b.db.cache.Store.TotalKappa[name] += count

		// subtract counts after minute
		go func() {
			if b.stop {
				return
			}
			time.Sleep(time.Minute)
			b.kappaCounter <- KappaData{Name: name, Count: -count}
		}()
	})
}

func (b *Bot) ircKeepAlive() {
	errChan := b.conn.ErrorChan()
	for b.conn.Connected() {
		err := <-errChan
		if !b.conn.Connected() {
			break
		}
		b.conn.Log.Printf("Error, disconnected: %s\n", err)
		b.restartIrc()
	}
}

func (b *Bot) restartIrc() {
	log.Println("restarting...")
	b.conn.ClearCallback("PRIVMSG")
	b.conn.Disconnect()
	b.stop = true
	time.Sleep(time.Minute)
	log.Println("restarted...")
	b.stop = false
	b.start()
}

func (b *Bot) start() {
	b.connect("kappakingbot", "kappakingbot", "irc.twitch.tv:6667")
	go b.countMinutes()
	go b.joinChannels()
	go b.updateCounts()
	b.trackKappas()
}

func launchBot(db *DB) {
	bot := NewBot(db)
	bot.start()
	bot.ircKeepAlive()
}
