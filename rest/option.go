package rest

import "github.com/txsvc/stdlib/v2/settings"

type ClientOption interface {
	Apply(ds *settings.DialSettings)
}

// WithEndpoint returns a ClientOption that overrides the default endpoint to be used for a service.
func WithEndpoint(url string) ClientOption {
	return withEndpoint(url)
}

type withEndpoint string

func (w withEndpoint) Apply(ds *settings.DialSettings) {
	ds.Endpoint = string(w)
}

// WithCredentials returns a ClientOption that overrides the default credentials used for a service.
func WithCredentials(clientid, secret string) ClientOption {
	return withCredentials{
		clientID:     clientid,
		clientSecret: secret,
	}
}

type withCredentials struct {
	clientID     string
	clientSecret string
}

func (w withCredentials) Apply(ds *settings.DialSettings) {
	if ds.Credentials == nil {
		ds.Credentials = &settings.Credentials{}
	}
	ds.Credentials.ClientID = w.clientID
	ds.Credentials.ClientSecret = w.clientSecret
}

// WithCredentials returns a ClientOption that overrides the default token used for a service.
func WithToken(clientid, token string) ClientOption {
	return withToken{
		clientID: clientid,
		token:    token,
	}
}

type withToken struct {
	clientID string
	token    string
}

func (w withToken) Apply(ds *settings.DialSettings) {
	if ds.Credentials == nil {
		ds.Credentials = &settings.Credentials{}
	}
	ds.Credentials.ClientID = w.clientID
	ds.Credentials.Token = w.token
}
