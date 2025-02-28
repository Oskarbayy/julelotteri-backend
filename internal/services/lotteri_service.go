package services

import (
	"context"
	"fmt"

	julelotteri "github.com/oskarbayy/julelotteri-backend/generated"
	"google.golang.org/protobuf/types/known/emptypb"
)

// create a struct
type LotteriService struct {
	julelotteri.UnimplementedLotteriServiceServer
}

// define methods for the struct | has to line up with the proto file
func (s *LotteriService) GetWinner(ctx context.Context, req *emptypb.Empty) (*julelotteri.Player, error) {
	fmt.Println("Client requesting winner...")
	var player *julelotteri.Player = &julelotteri.Player{
		Name: "Test Navn",
		Id:   42,
	}

	fmt.Println("Responding to client...")
	return player, nil
}
