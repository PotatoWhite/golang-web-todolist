package google

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// redirect to oauth login page
func RedirectToGoogleLoginPage(writer http.ResponseWriter, request *http.Request) {
	// generate transantionId
	state := generateStateOauthCookie(writer)
	url := GoogleOAuthConfig.AuthCodeURL(state)
	http.Redirect(writer, request, url, http.StatusTemporaryRedirect)
}

// redirected from oauth login page
func CallBackOAuthResultAndPrintUserInfo(writer http.ResponseWriter, request *http.Request) {
	oauthstate, _ := request.Cookie("oauthstate")

	// check fail and fast exit
	if request.FormValue("state") != oauthstate.Value {
		errMsg := fmt.Sprintf("invalid oauth oauth state cookie:%s state:%s\n", oauthstate.Value, request.FormValue("state"))
		log.Println(errMsg)
		http.Error(writer, errMsg, http.StatusInternalServerError)
		return
	}

	// get user info
	data, err := getGoogleUserInfo(request.FormValue("code"))
	// handle error
	if err != nil {
		log.Println(err.Error())
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// storeSessionId info into session cookie
	var userInfo GoogleUserId
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		errMsg := fmt.Sprintf("invalid oauth oauth state cookie:%s state:%s\n", oauthstate.Value, request.FormValue("state"))
		log.Println(err.Error())
		http.Error(writer, errMsg, http.StatusInternalServerError)
		return
	}

	err = storeSessionId(writer, request, userInfo.Id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
}
