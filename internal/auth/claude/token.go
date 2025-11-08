// Package claude provides authentication and token management functionality
// for Anthropic's Claude AI services. It handles OAuth2 token storage, serialization,
// and retrieval for maintaining authenticated sessions with the Claude API.
package claude

import (
	"fmt"

	"github.com/router-for-me/CLIProxyAPI/v6/internal/accounts"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/misc"
)

// ClaudeTokenStorage stores OAuth2 token information for Anthropic Claude API authentication.
// It maintains compatibility with the existing auth system while adding Claude-specific fields
// for managing access tokens, refresh tokens, and user account information.
type ClaudeTokenStorage struct {
	// IDToken is the JWT ID token containing user claims and identity information.
	IDToken string `json:"id_token"`

	// AccessToken is the OAuth2 access token used for authenticating API requests.
	AccessToken string `json:"access_token"`

	// RefreshToken is used to obtain new access tokens when the current one expires.
	RefreshToken string `json:"refresh_token"`

	// LastRefresh is the timestamp of the last token refresh operation.
	LastRefresh string `json:"last_refresh"`

	// Email is the Anthropic account email address associated with this token.
	Email string `json:"email"`

	// Type indicates the authentication provider type, always "claude" for this storage.
	Type string `json:"type"`

	// Expire is the timestamp when the current access token expires.
	Expire string `json:"expired"`
}

// SaveTokenToFile serializes the Claude token storage to a JSON file.
// This method creates the necessary directory structure and writes the token
// data in JSON format to the specified file path for persistent storage.
//
// Parameters:
//   - authFilePath: The full path where the token file should be saved
//
// Returns:
//   - error: An error if the operation fails, nil otherwise
func (ts *ClaudeTokenStorage) SaveTokenToFile(authFilePath string) error {
	if ts.Email == "" {
		return fmt.Errorf("email is required for Claude account identification")
	}
	misc.LogSavingCredentials(authFilePath)
	accountID := ts.Email
	return accounts.SaveProviderAccount("claude", accountID, func(existing map[string]any) map[string]any {
		for k, v := range map[string]any{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
			"last_refresh":  ts.LastRefresh,
			"email":         ts.Email,
			"expired":       ts.Expire,
		} {
			if s, ok := v.(string); ok && s == "" {
				continue
			}
			existing[k] = v
		}
		return existing
	})
}
