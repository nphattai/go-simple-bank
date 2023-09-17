package gapi

import (
	"fmt"

	db "github.com/nphattai/go-simple-bank/db/sqlc"
	"github.com/nphattai/go-simple-bank/pb"
	"github.com/nphattai/go-simple-bank/token"
	"github.com/nphattai/go-simple-bank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

func NewSever(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(util.RandomString(32))
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker, config: config}

	return server, nil
}
