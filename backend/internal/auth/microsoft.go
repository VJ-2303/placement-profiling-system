package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var (
	ErrFailedToGetUserInfo = errors.New("failed to get user info from Microsoft")
)

type MicrosoftOAuth struct {
	Config *oauth2.Config
}

type MicrosoftUserInfo struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	Mail              string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
}

func NewMicrosoftOAuth(clientID, clientSecret, redirectURL string) *MicrosoftOAuth {
	return &MicrosoftOAuth{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "profile", "email", "User.Read"},
			Endpoint:     microsoft.AzureADEndpoint("common"),
		},
	}
}

func (m *MicrosoftOAuth) GetAuthURL(state string) string {
	return m.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (m *MicrosoftOAuth) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return m.Config.Exchange(ctx, code)
}

func (m *MicrosoftOAuth) GetUserInfo(ctx context.Context, token *oauth2.Token) (*MicrosoftUserInfo, error) {
	client := m.Config.Client(ctx, token)

	resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrFailedToGetUserInfo
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo MicrosoftUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (m *MicrosoftOAuth) GetUserPhotoURL(ctx context.Context, token *oauth2.Token, userID string) string {
	// Use UI Avatars service to generate placeholder images based on name
	// This will be replaced with actual uploaded images later
	// Returns empty string - we'll generate the placeholder URL in the backend based on name
	return ""
}
