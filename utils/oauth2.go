package utils

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOauthConfig *oauth2.Config

func InitOAuth2() error {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		return err
	}

	// Initialize Google OAuth2 configuration
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	return nil
}

func GetGoogleUserInfo(ctx context.Context, code string) (*oauth2.Token, error) {
	return GoogleOauthConfig.Exchange(ctx, code)
}
