package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"math/rand/v2"
	"strconv"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
	"vincadrn.com/santuy/internal/model"
)

type OAuthURLResponse struct {
	Status   int    `json:"status"`
	OAuthURL string `json:"oauth_url"`
}

func RequestAuth() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		store := Session()
		session, err := store.Get(r, API_SESSION_NAME)
		if err != nil {
			log.Fatal(err)
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
		log.Println("---- Session in `login`:", session.Values, "\nURL:", url)

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
		if hostName != conf.Client.Host {
			http.Error(w, "invalid hostname", http.StatusForbidden)
			return
		}
		savedState := session.Values["oauth2_state"]
		if state != savedState {
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

		// e := json.NewEncoder(w)
		// e.SetIndent("", "  ")
		// e.Encode(*token)

		ctx := context.Background()
		peopleService, err := people.NewService(ctx, option.WithTokenSource(OauthConfig.TokenSource(ctx, token)))

		userInfo, err := peopleService.People.Get("people/me").PersonFields("emailAddresses").Do()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		emailAddress := userInfo.EmailAddresses[0].Value
		log.Println("User info: ", *userInfo)
		log.Println("Email addresses: ", userInfo.EmailAddresses)
		log.Println("Email address: ", emailAddress)
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
