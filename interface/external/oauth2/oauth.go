package oauth2

type OauthClient interface {
	GetAccessToken() (string, error)
}
