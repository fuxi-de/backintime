package services

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type GameService struct {
	DB *sql.DB
}

type GameData struct {
	UserId    string
	GameState string
}

type GameEntry struct {
	ID string
	GameData
}

func (gameService *GameService) CreateGame(token string, tokenService TokenService) *GameEntry {
	gameEntry := &GameEntry{
		ID: uuid.New().String(),
		GameData: GameData{
			UserId:    *tokenService.GetIdByToken(token),
			GameState: "",
		},
	}
	fmt.Println(gameEntry)
	sqlStmt := `
	insert into games(id, user, state) values(?, ?, ?)
	`
	_, err := gameService.DB.Exec(sqlStmt, gameEntry.ID, gameEntry.UserId, gameEntry.GameState)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil
	}

	return gameEntry
}
