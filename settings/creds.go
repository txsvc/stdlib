// Most of the code is lifted from
// https://github.com/googleapis/google-api-go-client/blob/main/internal/settings.go
//
// For details and copyright etc. see above url.
package settings

import (
	"strings"

	"github.com/txsvc/stdlib/v2"
)

const (
	ProjectID    = "PROJECT_ID"
	ClientID     = "CLIENT_ID"
	ClientSecret = "CLIENT_SECRET"
	AccessToken  = "ACCESS_TOKEN"

	StateInit       State = iota - 2 // waiting to swap tokens
	StateInvalid                     // a config in this state should not be used
	StateUndefined                   // logged out
	StateAuthorized                  // logged in
)

type (
	State int

	Credentials struct {
		ProjectID    string `json:"project_id,omitempty"`
		ClientID     string `json:"client_id,omitempty"`
		ClientSecret string `json:"client_secret,omitempty"`
		Token        string `json:"token,omitempty"`
		Status       State  `json:"status,omitempty"`
		Expires      int64  `json:"expires,omitempty"` // 0 = never, > 0 = unix timestamp, < 0 = invalid
	}
)

func CredentialsFromEnv() *Credentials {
	c := &Credentials{
		ProjectID: stdlib.GetString(ProjectID, ""),
		ClientID:  stdlib.GetString(ClientID, ""),
		Token:     stdlib.GetString(AccessToken, ""),
		Status:    StateInit,
	}
	if c.Token == "" {
		c.ClientSecret = stdlib.GetString(ClientSecret, "")
	}
	return c
}

// Clone returns a deep copy of the credentials
func (c *Credentials) Clone() *Credentials {
	return &Credentials{
		ProjectID:    c.ProjectID,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Token:        c.Token,
		Status:       c.Status,
		Expires:      c.Expires,
	}
}

func (c *Credentials) Key() string {
	return strings.ToLower(c.ProjectID + "." + c.ClientID)
}

func (c *Credentials) HashedKey() string {
	return stdlib.Fingerprint(c.Key())
}

// IsValid test if Crendentials is valid
func (c *Credentials) IsValid() bool {
	if c.Expires < 0 || c.Status == StateInvalid || len(c.ClientID) == 0 || (len(c.Token) == 0 && len(c.ClientID) == 0) {
		return false
	}
	return !c.Expired()
}

// Expired only verifies just that, does not check all other attributes
func (c *Credentials) Expired() bool {
	if c.Expires == 0 {
		return false
	}
	return c.Expires < stdlib.Now()
}
