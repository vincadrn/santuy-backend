package auth

import "os"

const (
	NORTHFLANK = "northflank"
	ORACLE     = "oracle"
	LOCAL      = "local"
)

type OAuthObtainer interface {
	GetOAuthClientID() string
	GetOAuthClientSecret() string
}

type NorthflankOAuthCredentials struct {
	ClientID     string
	ClientSecret string
}

func (northflank NorthflankOAuthCredentials) GetOAuthClientID() string {
	panic("Not implemented")
}

func (northflank NorthflankOAuthCredentials) GetOAuthClientSecret() string {
	panic("Not implemented")
}

type OracleOAuthCredentials struct {
	ClientID     string
	ClientSecret string
}

func (oracle OracleOAuthCredentials) GetOAuthClientID() string {
	panic("Not implemented")
}

func (oracle OracleOAuthCredentials) GetOAuthClientSecret() string {
	panic("Not implemented")
}

type LocalOAuthCredentials struct {
	ClientID     string
	ClientSecret string
}

func (local LocalOAuthCredentials) GetOAuthClientID() string {
	return os.Getenv("OAUTH_CLIENT_ID")
}

func (local LocalOAuthCredentials) GetOAuthClientSecret() string {
	return os.Getenv("OAUTH_CLIENT_SECRET")
}
