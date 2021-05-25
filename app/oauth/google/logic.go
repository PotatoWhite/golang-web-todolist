package google

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

// create transaction id
func generateStateOauthCookie(writer http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(writer, cookie)

	return state
}

// get user Info from googleapis
func getGoogleUserInfo(code string) ([]byte, error) {
	token, err := GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	response, err := http.Get(OauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to GetUserInfo %s\n", err.Error())
	}

	return ioutil.ReadAll(response.Body)
}
