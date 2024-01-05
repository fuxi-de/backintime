package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type TokenService struct {
	DB *sql.DB
}

type UserResponse struct {
	Name string `json:"display_name"`
	Mail string `json:"email"`
}

type TokenData struct {
	Mail         string `json:"mail"`
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expires      int    `json:"expires_in"`
}

type TokenEntry struct {
	ID string
	TokenData
}

func (tokenService *TokenService) SaveUser(tokenData TokenData) *TokenEntry {
	tokenEntry := &TokenEntry{
		ID:        uuid.New().String(),
		TokenData: tokenData,
	}

	tokenEntry.Mail = tokenService.GetUserInfo(tokenEntry.Token).Mail

	sqlStmt := `
	insert into tokens(id, mail, token, refresh_token, expires) values(?, ?, ?, ?, ?)
    on conflict do update set 
      token=excluded.token,
      refresh_token=excluded.refresh_token,
      expires=excluded.expires
	`
	_, err := tokenService.DB.Exec(sqlStmt, tokenEntry.ID, tokenEntry.Mail, tokenEntry.Token, tokenEntry.RefreshToken, tokenEntry.Expires)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil
	}
	return tokenEntry
}

func (tokenService *TokenService) GetUserInfo(token string) *UserResponse {
	spotifyProfileEndpoint := "https://api.spotify.com/v1/me"

	r, err := http.NewRequest("GET", spotifyProfileEndpoint, nil)

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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var userResponse UserResponse
	json.Unmarshal(body, &userResponse)

	return &userResponse
}

func (tokenService *TokenService) ExtractToken(tokenResponse string) string {
	var token TokenData
	json.Unmarshal([]byte(tokenResponse), &token)
	fmt.Println(token)
	return token.Token
}

func (tokenService *TokenService) CreateNewUser(tokenResponse string) *TokenEntry {
	var token TokenData
	json.Unmarshal([]byte(tokenResponse), &token)
	fmt.Println(token)

	return tokenService.SaveUser(token)
}

func (tokenService *TokenService) GetIdByToken(token string) *string {
	stmt, err := tokenService.DB.Prepare("select id from tokens where token = ?")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer stmt.Close()
	var id string
	err = stmt.QueryRow(token).Scan(&id)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &id
}

// func (tokenService *TokenService) GetAccessToken(email string) *TokenResponse {
// 	fmt.Println("getting token from db for: " + email)
//
// 	stmt, err := tokenService.DB.Prepare("select token from users where mail = ?")
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil
// 	}
// 	defer stmt.Close()
// 	var token string
// 	err = stmt.QueryRow(email).Scan(&token)
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil
// 	}
// 	var accessToken TokenResponse
// 	json.Unmarshal([]byte(token), &accessToken)
// 	return &accessToken
// }
