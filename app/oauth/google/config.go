package google

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"os"
)

var storeSessionId func(http.ResponseWriter, *http.Request, string) error

func SetSessionHandler(cookieStore func(writer http.ResponseWriter, request *http.Request, id string) error) {
	storeSessionId = cookieStore
}

// user info url
const OauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

// configuration
var GoogleOAuthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:3000/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

type GoogleUserId struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}
