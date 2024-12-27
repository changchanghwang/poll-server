//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"poll.ant/internal/libs/db"
	"poll.ant/internal/server"
	"poll.ant/internal/services/users"
)

func InitializeServer() (*server.Server, error) {
	wire.Build(server.ProviderSet, db.InitDb, users.UserSet)
	return &server.Server{}, nil
}
