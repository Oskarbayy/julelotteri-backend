package services

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"

	julelotteri "github.com/oskarbayy/julelotteri-backend/generated"
	"github.com/xuri/excelize/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

// create a struct
type LotteriService struct {
	julelotteri.UnimplementedLotteriServiceServer
	DB *sql.DB
}

// define methods for the struct | has to line up with the proto file
func (s *LotteriService) GetWinner(ctx context.Context, req *emptypb.Empty) (*julelotteri.Player, error) {
	fmt.Println("Client requesting winner...")

	var id int32
	var name string
	var won bool
	// Select only players who haven't won (won = 0)
	err := s.DB.QueryRow("SELECT id, name, won FROM players WHERE won = 0 ORDER BY RANDOM() LIMIT 1").Scan(&id, &name, &won)
	if err != nil {
		if err == sql.ErrNoRows {
			// No players available with won=false, return dummy data.
			name = "_____"
			id = 0
			won = false
			fmt.Println("No players with won=false found, returning dummy data.")
		} else {
			return nil, fmt.Errorf("database query error: %v", err)
		}
	}

	// update the player so that they are marked as won
	_, updateErr := s.DB.Exec("UPDATE players SET won = 1 WHERE id = ?", id)
	if updateErr != nil {
		log.Printf("Failed to update player's won status: %v", updateErr)
	}

	player := &julelotteri.Player{
		Name: name,
		Id:   id,
		Won:  won,
	}

	fmt.Println("Responding to client with a valid winner...")
	return player, nil
}

func (s *LotteriService) ImportExcelFile(ctx context.Context, req *julelotteri.ImportExcelFileRequest) (*julelotteri.ImportExcelFileResponse, error) {
	// Clear the existing data in the players table.
	_, err := s.DB.Exec("DELETE FROM players")
	if err != nil {
		log.Printf("Failed to clear players table: %v", err)
		return &julelotteri.ImportExcelFileResponse{Success: false}, err
	}
	log.Println("Players table cleared successfully.")

	// Access the file data as a []byte slice.
	fileData := req.FileData
	reader := bytes.NewReader(fileData)

	// Open the Excel file using excelize.
	f, err := excelize.OpenReader(reader)
	if err != nil {
		log.Printf("Failed to open Excel file: %v", err)
		return &julelotteri.ImportExcelFileResponse{Success: false}, err
	}
	defer f.Close()

	// Get rows from "Ark1".
	rows, err := f.GetRows("Ark1")
	if err != nil {
		log.Printf("Failed to read rows: %v", err)
		return &julelotteri.ImportExcelFileResponse{Success: false}, err
	}

	for i, row := range rows {
		if i == 0 {
			// Skip header row.
			continue
		}
		if len(row) < 3 {
			continue
		}
		name := row[1]
		number := row[2]

		// Insert the data into the database.
		// Note: Make sure you use three placeholders for name, number, and won.
		_, err = s.DB.Exec("INSERT INTO players (name, number, won) VALUES (?, ?, ?)", name, number, false)
		if err != nil {
			log.Printf("Failed to insert player (%s, %s): %v", name, number, err)
			continue
		}
		log.Printf("Importing player: Name=%s, Number=%s", name, number)
	}

	return &julelotteri.ImportExcelFileResponse{Success: true}, nil
}

func (s *LotteriService) GetPlayers(ctx context.Context, req *emptypb.Empty) (*julelotteri.PlayerList, error) {
	// Query for players where won is false (0).
	rows, err := s.DB.Query("SELECT id, name, won FROM players WHERE won = 0")
	if err != nil {
		return nil, fmt.Errorf("failed to query players: %v", err)
	}
	defer rows.Close()

	players := []*julelotteri.Player{}
	for rows.Next() {
		var number int32
		var name string
		var won bool
		if err := rows.Scan(&number, &name, &won); err != nil {
			return nil, fmt.Errorf("failed to scan player: %v", err)
		}
		players = append(players, &julelotteri.Player{
			Id:   number,
			Name: name,
			Won:  won,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over players: %v", err)
	}
	print(&julelotteri.PlayerList{Players: players})
	return &julelotteri.PlayerList{Players: players}, nil
}
