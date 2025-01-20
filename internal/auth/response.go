package auth

type OAuthURLResponse struct {
	Status   int    `json:"status"`
	OAuthURL string `json:"oauth_url"`
}
