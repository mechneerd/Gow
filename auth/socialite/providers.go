package socialite

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// AbstractProvider provides common OAuth logic.
type AbstractProvider struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// GoogleProvider implements Google OAuth.
type GoogleProvider struct {
	AbstractProvider
}

func NewGoogleProvider(clientID, clientSecret, redirectURL string) *GoogleProvider {
	return &GoogleProvider{
		AbstractProvider: AbstractProvider{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "profile", "email"},
		},
	}
}

func (p *GoogleProvider) Redirect(state string) string {
	u, _ := url.Parse("https://accounts.google.com/o/oauth2/v2/auth")
	q := u.Query()
	q.Set("client_id", p.ClientID)
	q.Set("redirect_uri", p.RedirectURL)
	q.Set("response_type", "code")
	q.Set("scope", strings.Join(p.Scopes, " "))
	q.Set("state", state)
	u.RawQuery = q.Encode()
	return u.String()
}

func (p *GoogleProvider) User(ctx context.Context, code string) (*User, error) {
	// Exchange code for token
	tokenURL := "https://oauth2.googleapis.com/token"
	resp, err := http.PostForm(tokenURL, url.Values{
		"client_id":     {p.ClientID},
		"client_secret": {p.ClientSecret},
		"redirect_uri":  {p.RedirectURL},
		"code":          {code},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token struct {
		AccessToken string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&token)

	// Get user info
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken
	resp2, err := http.Get(userInfoURL)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()
	body, _ := io.ReadAll(resp2.Body)

	var info struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	json.Unmarshal(body, &info)

	return &User{
		ID:     info.ID,
		Email:  info.Email,
		Name:   info.Name,
		Avatar: info.Picture,
		Token:  token.AccessToken,
	}, nil
}

// GitHubProvider implements GitHub OAuth.
type GitHubProvider struct {
	AbstractProvider
}

func NewGitHubProvider(clientID, clientSecret, redirectURL string) *GitHubProvider {
	return &GitHubProvider{
		AbstractProvider: AbstractProvider{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"user:email"},
		},
	}
}

func (p *GitHubProvider) Redirect(state string) string {
	u, _ := url.Parse("https://github.com/login/oauth/authorize")
	q := u.Query()
	q.Set("client_id", p.ClientID)
	q.Set("redirect_uri", p.RedirectURL)
	q.Set("scope", strings.Join(p.Scopes, " "))
	q.Set("state", state)
	u.RawQuery = q.Encode()
	return u.String()
}

func (p *GitHubProvider) User(ctx context.Context, code string) (*User, error) {
	tokenURL := "https://github.com/login/oauth/access_token"
	resp, err := http.PostForm(tokenURL, url.Values{
		"client_id":     {p.ClientID},
		"client_secret": {p.ClientSecret},
		"code":          {code},
		"redirect_uri":  {p.RedirectURL},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	vals, _ := url.ParseQuery(string(body))
	accessToken := vals.Get("access_token")

	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp2, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	var info struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
		Email string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	json.NewDecoder(resp2.Body).Decode(&info)

	return &User{
		ID:     fmt.Sprintf("%d", info.ID),
		Email:  info.Email,
		Name:   info.Name,
		Nickname: info.Login,
		Avatar: info.AvatarURL,
		Token:  accessToken,
	}, nil
}

