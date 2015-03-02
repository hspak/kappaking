package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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

var LiveStreams map[string]bool

// It is OK if this returns nil. It will only prevent the DB from updating
func getTopStreams(first bool) *Streams {
	res, err := http.Get("https://api.twitch.tv/kraken/streams")
	if err != nil {
		log.Println("info: no response from api")
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error: could not read api response")
		return nil
	}

	dat := new(Streams)
	if err := json.Unmarshal(body, dat); err != nil {
		log.Println("could not parse json")
		return nil
	}

	// TODO: think of a more efficient way
	if first {
		LiveStreams = make(map[string]bool)
		for _, d := range dat.Stream {
			LiveStreams[d.Channel.DisplayName] = true
		}
	} else {
		for k, _ := range LiveStreams {
			LiveStreams[k] = false
		}
		for _, d := range dat.Stream {
			LiveStreams[d.Channel.DisplayName] = true
		}
	}

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

func returnJSON(db *sql.DB) string {
	d, err := queryDB(db)
	if err != nil {
		// TODO: do more
		return "error"
	}
	wrapper := &Wrapper{Streams: d}
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
