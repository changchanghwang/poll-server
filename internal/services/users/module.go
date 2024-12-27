package users

import (
	"github.com/google/wire"
	"poll.ant/internal/services/users/application"
	"poll.ant/internal/services/users/infrastructure"
	"poll.ant/internal/services/users/presentation"
)

var UserSet = wire.NewSet(
	infrastructure.NewUserRepository,
	application.NewUserService,
	presentation.NewUserController,
)
