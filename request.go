package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Data struct {
	Display_name string
	Game         string
	Viewers      int
	Logo         string
	Status       string
	Url          string
}

type Wrapper struct {
	Streams []Data
}

type Streams struct {
	// {{{
	Stream []struct {
		// Id   int64  `json:"_id"`
		Game    string `json:"game"`
		Viewers int    `json:"viewers"`
		// CreatedAt string `json:"created_at"`
		// Links     struct {
		// Self string `json:"self"`
		// } `json:"_links"`
		// Preview struct {
		// Small    string `json:"small"`
		// Medium   string `json:"medium"`
		// Large    string `json:"large"`
		// Template string `json:"template"`
		// } `json:"preview"`
		Channel struct {
			// Links struct {
			// Self          string `json:"self"`
			// Follows       string `json:"follows"`
			// Commercial    string `json:"commercial"`
			// StreamKey     string `json:"stream_key"`
			// Chat          string `json:"chat"`
			// Features      string `json:"features"`
			// Subscriptions string `json:"subscriptions"`
			// Editors       string `json:"editors"`
			// Videos        string `json:"videos"`
			// Teams         string `json:"teams"`
			// } `json:"_links"`
			// Background                   interface{} `json:"background"`
			// Banner                       interface{} `json:"banner"`
			// BroadcasterLanguage          string      `json:"broadcaster_language"`
			DisplayName string `json:"display_name"`
			// Game        string `json:"game"`
			Logo string `json:"logo"`
			// Mature                       bool        `json:"mature"`
			Status string `json:"status"`
			// Partner                      bool        `json:"partner"`
			Url string `json:"url"`
			// VideoBanner                  string      `json:"video_banner"`
			// Id                           int         `json:"_id"`
			// Name                         string      `json:"name"`
			// CreatedAt                    string      `json:"created_at"`
			// UpdatedAt                    string      `json:"updated_at"`
			// Delay                        int         `json:"delay"`
			// Followers                    int         `json:"followers"`
			// ProfileBanner                string      `json:"profile_banner"`
			// ProfileBannerBackgroundColor interface{} `json:"profile_banner_background_color"`
			// Views                        int         `json:"views"`
			// Language                     string      `json:"language"`
		} `json:"channel"`
	} `json:"streams"`
	// Total int `json:"_total"`
	// Links struct {
	// Self     string `json:"self"`
	// Next     string `json:"next"`
	// Featured string `json:"featured"`
	// Summary  string `json:"summary"`
	// Followed string `json:"followed"`
	// } `json:"_links"`
	// }}}
}

func getTopStreams() *Streams {
	res, err := http.Get("https://api.twitch.tv/kraken/streams")
	if err != nil {
		log.Println("no response from api")
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("could not read response")
	}
	dat := new(Streams)
	if err := json.Unmarshal(body, dat); err != nil {
		log.Fatal("could not read json")
	}
	// prettyPrint(dat)
	fmt.Println(returnJSON(dat))
	return dat
}

func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.Marshal(v)

	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

func returnJSON(streams *Streams) string {
	d := make([]Data, 25)
	wrapper := &Wrapper{Stream: d}

	for i := 0; i < 25; i++ {
		d[i].Display_name = streams.Stream[i].Channel.DisplayName
		d[i].Game = streams.Stream[i].Game
		d[i].Viewers = streams.Stream[i].Viewers
		d[i].Logo = streams.Stream[i].Channel.Logo
		d[i].Status = streams.Stream[i].Channel.Status
		d[i].Url = streams.Stream[i].Channel.Url
	}

	out, err := JSONMarshal(&wrapper, true)
	if err != nil {
		// TODO: do more
		return "error"
	}
	return string(out)
}

func prettyPrint(streams *Streams) {
	for _, stream := range streams.Stream {
		fmt.Println("=== STREAM ===")
		fmt.Println(stream.Channel.DisplayName)
		fmt.Println(stream.Game)
		fmt.Println(stream.Viewers)
	}
}
