package auth

import (
	"github.com/google/wire"
	"poll.ant/internal/libs/oauth"
	"poll.ant/internal/services/auth/application"
	"poll.ant/internal/services/auth/infrastructure"
	"poll.ant/internal/services/auth/presentation"
)

var AuthSet = wire.NewSet(
	infrastructure.New,
	application.NewAuthService,
	presentation.NewAuthController,
	oauth.NewOAuthProvider,
)
