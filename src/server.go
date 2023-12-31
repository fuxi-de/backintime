package main

import (
	"context"
	"database/sql"
	b64 "encoding/base64"
	"fmt"
	"fuxifuchs/backintime/src/middleware"
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
	create table if not exists tokens (id text not null primary key, token text not null, refresh_token text not null, expires integer not null);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	spotifyClientId := os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	tokenService := &services.TokenService{
		DB: db,
	}

	e := echo.New()

	e.Static("/static", "assets")

	e.Static("/css", "dist")

	e.GET("/", func(c echo.Context) error {
		homepage := templates.HomePage("Florian")
		return homepage.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/login", func(c echo.Context) error {
		loginPage := templates.LoginPage()
		return loginPage.Render(context.Background(), c.Response().Writer)
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
		tokenEntry := tokenService.CreateNewUser(string(b))
		callbackPage := templates.CallbackPage(tokenEntry.Token)
		return callbackPage.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/redirect", func(c echo.Context) error {
		target := c.QueryParam("target")
		fmt.Println(target)
		c.Response().Header().Set("HX-Redirect", "/")
		return c.String(200, "done")
	})
	r := e.Group("/user")
	{
		r.Use(middleware.Auth)
		r.GET("/", func(c echo.Context) error {
			fmt.Println("/user/ called")
			token := c.Get("token")
			if str, ok := token.(string); ok {
				userPage := templates.User("Flo", str)
				return userPage.Render(context.Background(), c.Response().Writer)
			} else {
				return c.Redirect(302, "/login")
			}
		})
		r.GET("/play/:device/:category", func(c echo.Context) error {
			token := c.Get("token")
			if str, ok := token.(string); ok {
				tokenService.PlaySong(c.Param("category"), str, c.Param("device"))
				return c.String(200, "done")
			} else {
				return c.Redirect(302, "/login")
			}
		})

	}

	e.Logger.Fatal(e.Start(":1312"))
}
