package infrastructure

import (
	"clean-architecture/pkg/framework"
	"context"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// NewGmailService -> receive gmail service client
func NewGmailService(logger framework.Logger, env *framework.Env) *gmail.Service {
	ctx := context.Background()

	config := oauth2.Config{
		ClientID:     env.MailClientID,
		ClientSecret: env.MailClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.send"},
	}
	expiry, _ := time.Parse("2006-01-02", "2020-10-25")
	token := oauth2.Token{
		RefreshToken: env.MailTokenType,
		TokenType:    "Bearer",
		Expiry:       expiry,
	}
	var tokenSource = config.TokenSource(ctx, &token)
	srv, err := gmail.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		logger.Fatal("failed to receive gmail client", err.Error())
	}
	return srv

}
