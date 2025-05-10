package session

import (
	"encoding/gob"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

type SessionProvider interface {
	GetUserName() (string, error)
	GetUserEmail() (string, error)
	GetCurrentGroupRole() (string, error)
	GetOAuth2State() (string, error)
	GetOAuth2Verifier() (string, error)
	SetUserName(name string) error
	SetUserEmail(email string) error
	SetCurrentGroupRole(groupRole GroupRole) (string, error)
	SetOAuth2State(state string) error
	SetOAuth2Verifier(verifier string) error
	ClearSession() error
}

type GroupRole struct {
	GroupId string
	Role    string
}

const SESSION_NAME string = "santuysrv"

var cookieStore *sessions.CookieStore

func init() {
	// Register custom struct
	gob.Register(&GroupRole{})

	// Cookie setup
	cookieStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   1 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   true,
	}
}

type sessionProvider struct {
	store *sessions.CookieStore
	name  string
	r     *http.Request
	w     http.ResponseWriter
}

func NewSessionProvider(w http.ResponseWriter, r *http.Request) *sessionProvider {
	return &sessionProvider{
		store: cookieStore,
		name:  SESSION_NAME,
		w:     w,
		r:     r,
	}
}

func (sp *sessionProvider) GetUserName() (string, error) {
	session, err := sp.store.Get(sp.r, sp.name)
	if err != nil {
		slog.Error("Cannot get session")
		slog.Error(err.Error())

		return "", err
	}

	userName, ok := session.Values["name"].(string)
	if !ok {
		msg := "Cannot get user name"
		slog.Error(msg)

		return "", errors.New(msg)
	}

	return userName, nil
}

func (sp *sessionProvider) GetUserEmail() (string, error) {
	session, err := sp.store.Get(sp.r, sp.name)
	if err != nil {
		slog.Error("Cannot get session")
		slog.Error(err.Error())

		return "", err
	}

	email, ok := session.Values["email"].(string)
	if !ok {
		msg := "Cannot get user email"
		slog.Error(msg)

		return "", errors.New(msg)
	}

	return email, nil
}

func (sp *sessionProvider) GetCurrentGroupRole() (GroupRole, error) {
	session, err := sp.store.Get(sp.r, sp.name)
	if err != nil {
		slog.Error("Cannot get session")
		slog.Error(err.Error())

		return GroupRole{}, err
	}

	groupRole, ok := session.Values["group_role"].(*GroupRole)
	if !ok {
		msg := "Cannot get group ID and role"
		slog.Error(msg)

		return GroupRole{}, errors.New(msg)
	}

	return *groupRole, nil
}

func (sp *sessionProvider) GetOAuth2Verifier() (string, error) {
	session, err := sp.store.Get(sp.r, sp.name)
	if err != nil {
		slog.Error("Cannot get session")
		slog.Error(err.Error())

		return "", err
	}

	oAuth2Verifier, ok := session.Values["oauth2_verifier"].(string)
	if !ok {
		msg := "Cannot get oauth verifier"
		slog.Error(msg)

		return "", errors.New(msg)
	}

	return oAuth2Verifier, nil
}

func (sp *sessionProvider) GetOAuth2State() (string, error) {
	session, err := sp.store.Get(sp.r, sp.name)
	if err != nil {
		slog.Error("Cannot get session")
		slog.Error(err.Error())

		return "", err
	}

	oAuth2State, ok := session.Values["oauth2_state"].(string)
	if !ok {
		msg := "Cannot get oauth state"
		slog.Error(msg)

		return "", errors.New(msg)
	}

	return oAuth2State, nil
}

func (sp *sessionProvider) SetUserName(name string) error {
	session, err := sp.store.Get(sp.r, sp.name)
	if err != nil {
		slog.Error("Cannot get session")
		slog.Error(err.Error())

		return err
	}

	session.Values["name"] = name

	return session.Save(sp.r, sp.w)
}

func (sp *sessionProvider) SetUserEmail(email string) error {
	session, err := sp.store.Get(sp.r, sp.name)
	if err != nil {
		slog.Error("Cannot get session")
		slog.Error(err.Error())

		return err
	}

	session.Values["email"] = email

	return session.Save(sp.r, sp.w)
}

func (sp *sessionProvider) SetOAuth2Verifier(oAuth2Verifier string) error {
	session, err := sp.store.Get(sp.r, sp.name)
	if err != nil {
		slog.Error("Cannot get session")
		slog.Error(err.Error())

		return err
	}

	session.Values["oauth2_verifier"] = oAuth2Verifier

	return session.Save(sp.r, sp.w)
}

func (sp *sessionProvider) SetCurrentGroupRole(groupRole GroupRole) error {
	session, err := sp.store.Get(sp.r, sp.name)
	if err != nil {
		slog.Error("Cannot get session")
		slog.Error(err.Error())

		return err
	}

	session.Values["group_role"] = groupRole

	return session.Save(sp.r, sp.w)
}

func (sp *sessionProvider) SetOAuth2State(oAuth2State string) error {
	session, err := sp.store.Get(sp.r, sp.name)
	if err != nil {
		slog.Error("Cannot get session")
		slog.Error(err.Error())

		return err
	}

	session.Values["oauth2_state"] = oAuth2State

	return session.Save(sp.r, sp.w)
}

func (sp *sessionProvider) ClearSession() error {
	session, err := sp.store.Get(sp.r, sp.name)
	if err != nil {
		slog.Error("Cannot get session")
		slog.Error(err.Error())

		return err
	}

	session.Values["name"] = nil
	session.Values["email"] = nil
	session.Values["group_role"] = nil
	session.Values["oauth2_verifier"] = nil
	session.Values["oauth2_state"] = nil

	return session.Save(sp.r, sp.w)
}
