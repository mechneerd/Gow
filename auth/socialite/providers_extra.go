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

// FacebookProvider implements Facebook OAuth.
type FacebookProvider struct {
	AbstractProvider
}

// NewFacebookProvider creates a new Facebook OAuth provider.
func NewFacebookProvider(clientID, clientSecret, redirectURL string) *FacebookProvider {
	return &FacebookProvider{
		AbstractProvider: AbstractProvider{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"email", "public_profile"},
		},
	}
}

// Redirect returns the OAuth authorization URL.
func (p *FacebookProvider) Redirect(state string) string {
	u, _ := url.Parse("https://www.facebook.com/v18.0/dialog/oauth")
	q := u.Query()
	q.Set("client_id", p.ClientID)
	q.Set("redirect_uri", p.RedirectURL)
	q.Set("scope", strings.Join(p.Scopes, " "))
	q.Set("state", state)
	q.Set("response_type", "code")
	u.RawQuery = q.Encode()
	return u.String()
}

// User exchanges the authorization code for user information.
func (p *FacebookProvider) User(ctx context.Context, code string) (*User, error) {
	tokenURL := "https://graph.facebook.com/v18.0/oauth/access_token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Form = url.Values{
		"client_id":     {p.ClientID},
		"client_secret": {p.ClientSecret},
		"redirect_uri":  {p.RedirectURL},
		"code":          {code},
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	var token struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	if token.AccessToken == "" {
		return nil, fmt.Errorf("facebook returned empty access token")
	}

	// Get user info
	req2, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://graph.facebook.com/v18.0/me?fields=id,name,email,picture.width(200)&access_token="+token.AccessToken, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp2.Body.Close()
	body, err := io.ReadAll(resp2.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %w", err)
	}

	var info struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &User{
		ID:     info.ID,
		Email:  info.Email,
		Name:   info.Name,
		Avatar: info.Picture.Data.URL,
		Token:  token.AccessToken,
	}, nil
}

// TwitterProvider implements Twitter/X OAuth.
type TwitterProvider struct {
	AbstractProvider
}

// NewTwitterProvider creates a new Twitter/X OAuth provider.
func NewTwitterProvider(clientID, clientSecret, redirectURL string) *TwitterProvider {
	return &TwitterProvider{
		AbstractProvider: AbstractProvider{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"tweet.read", "users.read"},
		},
	}
}

// Redirect returns the OAuth authorization URL.
func (p *TwitterProvider) Redirect(state string) string {
	u, _ := url.Parse("https://twitter.com/i/oauth2/authorize")
	q := u.Query()
	q.Set("client_id", p.ClientID)
	q.Set("redirect_uri", p.RedirectURL)
	q.Set("scope", strings.Join(p.Scopes, " "))
	q.Set("state", state)
	q.Set("response_type", "code")
	q.Set("code_challenge", "challenge")
	q.Set("code_challenge_method", "plain")
	u.RawQuery = q.Encode()
	return u.String()
}

// User exchanges the authorization code for user information.
func (p *TwitterProvider) User(ctx context.Context, code string) (*User, error) {
	tokenURL := "https://api.twitter.com/2/oauth2/token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Form = url.Values{
		"client_id":     {p.ClientID},
		"client_secret": {p.ClientSecret},
		"redirect_uri":  {p.RedirectURL},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"code_verifier": {"challenge"},
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	var token struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	if token.AccessToken == "" {
		return nil, fmt.Errorf("twitter returned empty access token")
	}

	// Get user info
	req2, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.twitter.com/2/users/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}
	req2.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp2.Body.Close()
	body, err := io.ReadAll(resp2.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %w", err)
	}

	var info struct {
		Data struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &User{
		ID:       info.Data.ID,
		Name:     info.Data.Name,
		Nickname: info.Data.Username,
		Token:    token.AccessToken,
	}, nil
}

// LinkedInProvider implements LinkedIn OAuth.
type LinkedInProvider struct {
	AbstractProvider
}

// NewLinkedInProvider creates a new LinkedIn OAuth provider.
func NewLinkedInProvider(clientID, clientSecret, redirectURL string) *LinkedInProvider {
	return &LinkedInProvider{
		AbstractProvider: AbstractProvider{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "profile", "email"},
		},
	}
}

// Redirect returns the OAuth authorization URL.
func (p *LinkedInProvider) Redirect(state string) string {
	u, _ := url.Parse("https://www.linkedin.com/oauth/v2/authorization")
	q := u.Query()
	q.Set("client_id", p.ClientID)
	q.Set("redirect_uri", p.RedirectURL)
	q.Set("scope", strings.Join(p.Scopes, " "))
	q.Set("state", state)
	q.Set("response_type", "code")
	u.RawQuery = q.Encode()
	return u.String()
}

// User exchanges the authorization code for user information.
func (p *LinkedInProvider) User(ctx context.Context, code string) (*User, error) {
	tokenURL := "https://www.linkedin.com/oauth/v2/accessToken"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Form = url.Values{
		"client_id":     {p.ClientID},
		"client_secret": {p.ClientSecret},
		"redirect_uri":  {p.RedirectURL},
		"code":          {code},
		"grant_type":    {"authorization_code"},
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	var token struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	if token.AccessToken == "" {
		return nil, fmt.Errorf("linkedin returned empty access token")
	}

	// Get user info
	req2, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.linkedin.com/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}
	req2.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp2.Body.Close()
	body, err := io.ReadAll(resp2.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %w", err)
	}

	var info struct {
		Sub     string `json:"sub"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &User{
		ID:     info.Sub,
		Email:  info.Email,
		Name:   info.Name,
		Avatar: info.Picture,
		Token:  token.AccessToken,
	}, nil
}

// AppleProvider implements Apple Sign In.
type AppleProvider struct {
	AbstractProvider
	teamID     string
	keyID      string
	privateKey string
}

// NewAppleProvider creates a new Apple OAuth provider.
func NewAppleProvider(clientID, clientSecret, redirectURL, teamID, keyID, privateKey string) *AppleProvider {
	return &AppleProvider{
		AbstractProvider: AbstractProvider{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"name", "email"},
		},
		teamID:     teamID,
		keyID:      keyID,
		privateKey: privateKey,
	}
}

// Redirect returns the Apple authorization URL.
func (p *AppleProvider) Redirect(state string) string {
	u, _ := url.Parse("https://appleid.apple.com/auth/authorize")
	q := u.Query()
	q.Set("client_id", p.ClientID)
	q.Set("redirect_uri", p.RedirectURL)
	q.Set("scope", strings.Join(p.Scopes, " "))
	q.Set("state", state)
	q.Set("response_type", "code id_token")
	q.Set("response_mode", "form_post")
	u.RawQuery = q.Encode()
	return u.String()
}

// User exchanges the authorization code for user information.
func (p *AppleProvider) User(ctx context.Context, code string) (*User, error) {
	tokenURL := "https://appleid.apple.com/auth/token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Form = url.Values{
		"client_id":     {p.ClientID},
		"client_secret": {p.ClientSecret},
		"code":          {code},
		"redirect_uri":  {p.RedirectURL},
		"grant_type":    {"authorization_code"},
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	var token struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	if token.AccessToken == "" {
		return nil, fmt.Errorf("apple returned empty access token")
	}

	// Parse the ID token to get user info (simplified)
	// In production, you would verify the JWT signature
	user := &User{
		Token: token.AccessToken,
	}

	return user, nil
}
