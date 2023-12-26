package main

import (
	"context"
	"database/sql"
	b64 "encoding/base64"
	"fmt"
	"fuxifuchs/backintime/src/services"
	"fuxifuchs/backintime/src/templates"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dotenvErr := godotenv.Load()
	if dotenvErr != nil {
		log.Fatalln("Error loading .env file")
	}

	db, err := sql.Open("sqlite3", "./db/backintime.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStmt := `
	create table if not exists users (id text not null primary key, name text not null, mail text not null, token text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	spotifyClientId := os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	userService := &services.UserService{
		DB: db,
	}

	e := echo.New()

	e.Static("/static", "assets")

	e.Static("/css", "dist")

	e.GET("/", func(c echo.Context) error {
		homepage := templates.HomePage("Florian")
		return homepage.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/auth/login", func(c echo.Context) error {
		scope := "streaming user-read-email user-read-private"

		spotifyAuthorizationUrl := "https://accounts.spotify.com/authorize/"

		authQueryParams := url.Values{}
		authQueryParams.Add("response_type", "code")
		authQueryParams.Add("scope", scope)
		authQueryParams.Add("redirect_uri", "http://localhost:1312/auth/callback")
		authQueryParams.Add("client_id", spotifyClientId)

		baseurl, err := url.Parse(spotifyAuthorizationUrl)
		if err != nil {
			fmt.Println("Malformed URL: ", err.Error())
			panic(err)
		}
		baseurl.RawQuery = authQueryParams.Encode()

		authUrl := baseurl.String()
		fmt.Println(authUrl)
		// return c.String(200, "ok")
		return c.Redirect(http.StatusFound, authUrl)
	})

	e.GET("/auth/callback", func(c echo.Context) error {
		code := c.QueryParam("code")

		spotifyTokenEndpoint := "https://accounts.spotify.com/api/token"
		formValues := url.Values{}
		formValues.Add("code", code)
		formValues.Add("redirect_uri", "http://localhost:1312/auth/callback")
		formValues.Add("grant_type", "authorization_code")

		encodedBasicAuth := b64.URLEncoding.EncodeToString([]byte(spotifyClientId + ":" + spotifyClientSecret))
		r, err := http.NewRequest("POST", spotifyTokenEndpoint, strings.NewReader(formValues.Encode()))

		r.Header.Add("Authorization", "Basic "+encodedBasicAuth)
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			log.Fatalln(err)
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

		user := userService.CreateNewUser(string(b))

		return c.Redirect(302, "/token/"+user.Mail)
	})

	e.GET("/token/:user", func(c echo.Context) error {
		token := userService.GetAccessToken(c.Param("user"))
		tokenPage := templates.UserPage("Flo", token.Token)
		return tokenPage.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/play/:user/:device", func(c echo.Context) error {
		token := userService.GetAccessToken(c.Param("user"))
		userService.PlaySong("something", token.Token, c.Param("device"))
		return c.String(200, "done")
	})
	e.Logger.Fatal(e.Start(":1312"))
}
