package auth

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"

	"crypto/rand"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"

	config "vincadrn.com/santuy/configs"
	"vincadrn.com/santuy/internal/model"
	"vincadrn.com/santuy/internal/service"
	"vincadrn.com/santuy/internal/session"
)

type OAuthURLResponse struct {
	Status   int    `json:"status"`
	OAuthURL string `json:"oauth_url"`
}

func RequestAuth(svc *service.AccountService, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		sessionProvider := session.NewSessionProvider(w, r)

		OauthConfig.RedirectURL = fmt.Sprintf("%s/oauth2", r.Header.Get("Origin"))
		verifier := oauth2.GenerateVerifier()
		randomBytes := make([]byte, 8)
		_, err := rand.Read(randomBytes)
		if err != nil {
			log.Fatal(err)
		}
		randomState := hex.EncodeToString(randomBytes)
		url := OauthConfig.AuthCodeURL(randomState, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))

		err = sessionProvider.SetOAuth2Verifier(verifier)
		if err != nil {
			slog.Error("Cannot set oauth verifier")
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		err = sessionProvider.SetOAuth2State(randomState)
		if err != nil {
			slog.Error("Cannot set oauth state")
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

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
			http.Error(w, "OAuth error", http.StatusInternalServerError)
			return
		}

		w.Write(buf.Bytes())
	})
}

func RequestSession(svc *service.AccountService, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		sessionProvider := session.NewSessionProvider(w, r)

		var redirectURI OAuthRedirectURI
		bodyBytes, err := io.ReadAll(r.Body)
		defer r.Body.Close()
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

		savedState, err := sessionProvider.GetOAuth2State()
		if state != savedState {
			slog.Error("Cannot retrieve saved oauth state", "state", state, "saved state", savedState)

			if err != nil {
				slog.Error(err.Error())
			}

			http.Error(w, "Invalid state", http.StatusInternalServerError)

			return
		}

		savedCodeVerifier, err := sessionProvider.GetOAuth2Verifier()
		if err != nil {
			slog.Error("Cannot retrieve saved oauth verifier", "state", state, "saved state", savedState)
			slog.Error(err.Error())

			http.Error(w, "Invalid verifier", http.StatusInternalServerError)

			return
		}

		token, err := OauthConfig.Exchange(context.Background(), code, oauth2.VerifierOption(savedCodeVerifier))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Google OAuth email retrieval
		ctx := context.Background()
		peopleService, err := people.NewService(ctx, option.WithTokenSource(OauthConfig.TokenSource(ctx, token)))

		userInfo, err := peopleService.People.Get("people/me").PersonFields("names,emailAddresses").Do()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userName := userInfo.Names[0].DisplayName
		emailAddress := userInfo.EmailAddresses[0].Value

		err = sessionProvider.SetUserName(userName)
		if err != nil {
			slog.Error("Cannot set user name")
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		err = sessionProvider.SetUserEmail(emailAddress)
		if err != nil {
			slog.Error("Cannot set user email")
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		// Save user in db
		user := model.User{
			Name:  userName,
			Email: emailAddress,
		}
		svc.SaveUser(ctx, &user)

		// Clear oauth session
		sessionProvider.SetOAuth2State("")
		sessionProvider.SetOAuth2Verifier("")

		w.Write([]byte("OK"))
	})
}

func RequestLogout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		sessionProvider := session.NewSessionProvider(w, r)

		sessionProvider.SetUserEmail("")
		sessionProvider.SetUserName("")
		sessionProvider.SetCurrentGroupRole(session.GroupRole{})
		sessionProvider.SetOAuth2State("")
		sessionProvider.SetOAuth2Verifier("")

		w.Write([]byte("OK"))
	})
}
