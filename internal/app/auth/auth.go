package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"vincadrn.com/santuy/config"
	"vincadrn.com/santuy/model"
)

var (
	conf = config.Configuration()

	OauthConfig = oauth2.Config{}
)

func init() {
	var obtainer OAuthObtainer

	switch conf.SecretKeeper {
	case ORACLE:
		obtainer = OracleOAuthCredentials{}
	case NORTHFLANK:
		obtainer = NorthflankOAuthCredentials{}
	case LOCAL:
		obtainer = LocalOAuthCredentials{}
	}

	redirectHost := ""
	if os.Getenv("ENVIRONMENT") == "LOCAL" {
		redirectHost = conf.Client.Host + ":" + conf.Client.Port
	} else {
		redirectHost = conf.Client.Host
	}

	OauthConfig = oauth2.Config{
		ClientID:     obtainer.GetOAuthClientID(),
		ClientSecret: obtainer.GetOAuthClientSecret(),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
		RedirectURL:  fmt.Sprintf("http://%s/oauth2", redirectHost),
	}
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	log.Println("Authenticating ...")

	store := Session()
	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		log.Fatal(err)
	}

	fooVal := session.Values["foo"]
	if fooVal == "bar" {
		w.Write([]byte("autheddddd!"))
		return
	}

	verifier := oauth2.GenerateVerifier()
	randomizer := rand.ChaCha8{}
	randomState := strconv.FormatUint(randomizer.Uint64(), 16)
	url := OauthConfig.AuthCodeURL(randomState, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))

	session.Values["oauth2_verifier"] = verifier
	session.Values["oauth2_state"] = randomState
	err = session.Save(r, w)
	if err != nil {
		model.ResponseWithErrorDefault(w, err, http.StatusInternalServerError)
		return
	}
	log.Println("---- Session in `login`:", session.Values)

	oauthResponse := OAuthURLResponse{
		Status:   http.StatusOK,
		OAuthURL: url,
	}
	payload, err := json.Marshal(oauthResponse)
	if err != nil {
		model.ResponseWithErrorDefault(w, err, http.StatusInternalServerError)
	}
	w.Write(payload)
}

func OAuthCallback(w http.ResponseWriter, r *http.Request) {
	store := Session()
	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		log.Fatal(err)
	}

	r.ParseForm()

	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "invalid code", http.StatusBadRequest)
		return
	}

	retrievedState := r.Form.Get("state")
	savedState := session.Values["oauth2_state"]
	if retrievedState != savedState {
		http.Error(w, "invalid state", http.StatusInternalServerError)
		return
	}

	savedCodeVerifier, codeVerifierIsValid := session.Values["oauth2_verifier"].(string)
	if !codeVerifierIsValid {
		http.Error(w, "invalid verifier", http.StatusInternalServerError)
		return
	}
	token, err := OauthConfig.Exchange(context.Background(), code, oauth2.VerifierOption(savedCodeVerifier))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["foo"] = "bar"
	session.Values[42] = 43
	err = session.Save(r, w)
	log.Println("Session value after oauth callback:", session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(*token)
	log.Println("---- Session in `callback`:", session.Values)
}
