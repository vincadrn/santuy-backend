package auth

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"crypto/rand"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"

	config "vincadrn.com/santuy/configs"
	"vincadrn.com/santuy/internal/model"
)

type OAuthURLResponse struct {
	Status   int    `json:"status"`
	OAuthURL string `json:"oauth_url"`
}

func RequestAuth() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		store := Session()
		session, err := store.Get(r, API_SESSION_NAME)
		if err != nil {
			log.Fatal(err)
		}

		OauthConfig.RedirectURL = fmt.Sprintf("%s/oauth2", r.Header.Get("Origin"))
		verifier := oauth2.GenerateVerifier()
		randomBytes := make([]byte, 8)
		_, err = rand.Read(randomBytes)
		if err != nil {
			log.Fatal(err)
		}
		randomState := hex.EncodeToString(randomBytes)
		url := OauthConfig.AuthCodeURL(randomState, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))

		session.Values["oauth2_verifier"] = verifier
		session.Values["oauth2_state"] = randomState
		err = session.Save(r, w)
		if err != nil {
			model.ResponseWithErrorDefault(w, err, http.StatusInternalServerError)
			return
		}

		oauthResponse := OAuthURLResponse{
			Status:   http.StatusOK,
			OAuthURL: url,
		}

		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)

		err = enc.Encode(oauthResponse)
		if err != nil {
			model.ResponseWithErrorDefault(w, err, http.StatusInternalServerError)
		}
		w.Write(buf.Bytes())
	})
}

func RequestSession() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		store := Session()
		session, err := store.Get(r, API_SESSION_NAME)
		if err != nil {
			log.Fatal(err)
		}

		defer r.Body.Close()

		var redirectURI OAuthRedirectURI
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(bodyBytes, &redirectURI)
		if err != nil {
			log.Fatal(err)
		}

		parsedURI, err := url.Parse(redirectURI.RedirectURI)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(parsedURI)

		hostName := parsedURI.Host
		code := parsedURI.Query()["code"][0]
		state := parsedURI.Query()["state"][0]
		if hostName != config.GetAllowedClientHost() {
			http.Error(w, "invalid hostname", http.StatusForbidden)
			return
		}
		savedState := session.Values["oauth2_state"]
		if state != savedState {
			log.Println("--- state:", state, "---- savedState:", savedState)
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

		ctx := context.Background()
		peopleService, err := people.NewService(ctx, option.WithTokenSource(OauthConfig.TokenSource(ctx, token)))

		userInfo, err := peopleService.People.Get("people/me").PersonFields("emailAddresses").Do()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		emailAddress := userInfo.EmailAddresses[0].Value
		session.Values["email"] = emailAddress
		err = session.Save(r, w)
		if err != nil {
			model.ResponseWithErrorDefault(w, err, http.StatusInternalServerError)
			return
		}
		log.Println("Session values:", session)

		w.Write([]byte("OK"))
	})
}
