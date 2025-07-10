package gapi

import (
	"fmt"

	db "github.com/terenjit/simplebank/db/sqlc"
	"github.com/terenjit/simplebank/pb"
	"github.com/terenjit/simplebank/token"
	"github.com/terenjit/simplebank/util"
	"github.com/terenjit/simplebank/worker"
)

// for gRPC server
type Server struct {
	pb.UnimplementedSimpleBankServer
	store           db.Store
	config          util.Config
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker:%v", err)
	}
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
