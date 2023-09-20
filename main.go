package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/nphattai/go-simple-bank/api"
	db "github.com/nphattai/go-simple-bank/db/sqlc"
	_ "github.com/nphattai/go-simple-bank/doc/statik"
	"github.com/nphattai/go-simple-bank/gapi"
	"github.com/nphattai/go-simple-bank/pb"
	"github.com/nphattai/go-simple-bank/util"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Error().Err(err).Msg("can not load config")
	}

	if config.AppEnv == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Error().Err(err).Msg("can not connect to db")
	}

	store := db.NewStore(conn)

	runDBMigration(config.MigrationSourceURL, config.DBSource)

	go runGrpcGatewayServer(config, store)

	runGrpcServer(config, store)
}

func runDBMigration(sourceURL string, dbURL string) {
	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		log.Error().Err(err).Msg("can not create migrate")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Error().Err(err).Msg("can not run migration")
	}
}

func runGrpcGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewSever(config, store)
	if err != nil {
		log.Error().Err(err).Msg("can not create server")
	}

	grpcMux := runtime.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Error().Err(err).Msg("can not register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Error().Err(err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Error().Err(err).Msg("can not create listener")

	}

	log.Printf("start api gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Error().Err(err).Msg("can not start api gateway server")

	}
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewSever(config, store)
	if err != nil {
		log.Error().Err(err).Msg("can not create server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterSimpleBankServer(grpcServer, server)

	// TODO: ????
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Error().Err(err).Msg("can not create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Error().Err(err).Msg("can not start gRPC server")

	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewSever(config, store)

	if err != nil {
		log.Error().Err(err).Msg("can not create server")
	}
	server.Start(config.HTTPServerAddress)
}
