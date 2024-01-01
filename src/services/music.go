package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

type (
	MusicService struct{}
	PlaylistData struct {
		Tracks struct {
			Total int `json:"total"`
		}
	}
)

func (musicService *MusicService) PlaySong(playlistUri string, token string, deviceId string) {
	fmt.Println("playing uri: " + playlistUri)

	spotifyPlayerEndpoint := "https://api.spotify.com/v1/me/player/play"

	params := url.Values{}
	params.Add("device_id", deviceId)

	baseurl, err := url.Parse(spotifyPlayerEndpoint)
	if err != nil {
		fmt.Println("Malformed URL: ", err.Error())
		panic(err)
	}
	baseurl.RawQuery = params.Encode()

	playerEndpoint := baseurl.String()

	playlistIndex := getRandomPlaylistIndex(playlistUri, token)

	fmt.Println(playlistIndex)

	jsonBody := []byte(`{"context_uri": "spotify:playlist:` + playlistUri + `", "offset": { "position":` + strconv.Itoa(playlistIndex) + `}, "position_ms": 0}`)
	body := bytes.NewReader(jsonBody)
	fmt.Println(body)
	r, err := http.NewRequest("PUT", playerEndpoint, body)

	r.Header.Add("Authorization", "Bearer "+token)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	fmt.Println("res: " + res.Status)
}

func getRandomPlaylistIndex(playlistUri string, token string) int {
	playlistEndpoint := "https://api.spotify.com/v1/playlists/" + playlistUri

	fmt.Println("musicService" + playlistEndpoint)
	r, err := http.NewRequest("GET", playlistEndpoint, nil)

	r.Header.Add("Authorization", "Bearer "+token)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var playlistData PlaylistData
	json.Unmarshal([]byte(b), &playlistData)

	log.Printf("librariesInformation: %+v", playlistData)
	randomIndex := getRandomNumberInRange(playlistData.Tracks.Total)
	return randomIndex
}

func getRandomNumberInRange(max int) int {
	randomIndex := rand.Intn(max-0+1) + 0
	fmt.Println(randomIndex)
	return randomIndex
}
