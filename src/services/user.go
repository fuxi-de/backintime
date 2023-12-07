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

type UserService struct {
	DB *sql.DB
}

type UserResponse struct {
	Name string `json:"display_name"`
	Mail string `json:"email"`
}

type TokenResponse struct {
	Token string `json:"access_token"`
}

type UserData struct {
	Name  string
	Mail  string
	Token string
}

type UserEntry struct {
	UserData
	ID string
}

func (us *UserService) SaveUser(userData UserData) *UserEntry {
	user := &UserEntry{
		ID:       uuid.New().String(),
		UserData: userData,
	}

	sqlStmt := `
	insert or replace into users(id, token, name, mail) values(?, ?, ?, ?)
	`
	_, err := us.DB.Exec(sqlStmt, user.ID, user.Token, user.Name, user.Mail)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil
	}
	return user
}

func (us *UserService) GetUserInfo(token string) *UserData {
	userData := &UserData{
		Token: token,
		Name:  "",
		Mail:  "",
	}

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
	var test UserResponse
	json.Unmarshal(body, &test)

	userData.Mail = test.Mail
	userData.Name = test.Name

	return userData
}

func (us *UserService) ExtractToken(tokenResponse string) string {
	var test TokenResponse
	json.Unmarshal([]byte(tokenResponse), &test)
	return test.Token
}

func (us *UserService) CreateNewUser(tokenResponse string) *UserEntry {
	accessToken := us.ExtractToken(tokenResponse)
	userData := us.GetUserInfo(accessToken)
	userData.Token = tokenResponse

	return us.SaveUser(*userData)
}

func (us *UserService) GetAccessToken(email string) *TokenResponse {
	fmt.Println("getting token from db for: " + email)

	stmt, err := us.DB.Prepare("select token from users where mail = ?")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer stmt.Close()
	var token string
	err = stmt.QueryRow(email).Scan(&token)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var accessToken TokenResponse
	json.Unmarshal([]byte(token), &accessToken)
	return &accessToken
}
