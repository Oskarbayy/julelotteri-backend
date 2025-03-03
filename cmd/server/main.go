package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
	julelotteri "github.com/oskarbayy/julelotteri-backend/generated"
	"github.com/oskarbayy/julelotteri-backend/internal/services"
	"google.golang.org/grpc"
)

func main() {
	print("Starting server...\n")

	// Setup the database
	db, err := sql.Open("sqlite3", "lotteri.db")
	if err != nil {
		log.Fatalf("Failed to open the database: %v", err)
	}
	defer db.Close() // stop db when program exits

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS players (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		number TEXT NOT NULL,
		won BOOL NOT NULL
	);`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Setup the server
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen on port :8080, %v", err)
	}

	grpcServer := grpc.NewServer()
	lotteriService := &services.LotteriService{DB: db}

	julelotteri.RegisterLotteriServiceServer(grpcServer, lotteriService)

	print("gRPC server is running on port 8080")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port :8080 %v", err)
	}
}
