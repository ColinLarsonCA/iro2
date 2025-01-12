package greeting

import (
	"context"
	"database/sql"

	"github.com/ColinLarsonCA/iro2/backend/pb"
)

type service struct {
	pb.UnimplementedGreetingServiceServer

	db *sql.DB
}

func NewService(db *sql.DB) *service {
	return &service{db: db}
}

func (s *service) GetGreeting(ctx context.Context, in *pb.GreetingRequest) (*pb.GreetingResponse, error) {
	var message string
	err := s.db.QueryRow("SELECT message FROM greetings ORDER BY created_at DESC LIMIT 1").Scan(&message)
	if err != nil {
		return nil, err
	}
	return &pb.GreetingResponse{Message: message}, nil
}
