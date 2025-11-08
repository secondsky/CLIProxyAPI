// Package codex provides authentication and token management functionality
// for OpenAI's Codex AI services. It handles OAuth2 token storage, serialization,
// and retrieval for maintaining authenticated sessions with the Codex API.
package codex

import (
	"fmt"

	"github.com/router-for-me/CLIProxyAPI/v6/internal/accounts"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/misc"
)

// CodexTokenStorage stores OAuth2 token information for OpenAI Codex API authentication.
// It maintains compatibility with the existing auth system while adding Codex-specific fields
// for managing access tokens, refresh tokens, and user account information.
type CodexTokenStorage struct {
	// IDToken is the JWT ID token containing user claims and identity information.
	IDToken string `json:"id_token"`
	// AccessToken is the OAuth2 access token used for authenticating API requests.
	AccessToken string `json:"access_token"`
	// RefreshToken is used to obtain new access tokens when the current one expires.
	RefreshToken string `json:"refresh_token"`
	// AccountID is the OpenAI account identifier associated with this token.
	AccountID string `json:"account_id"`
	// LastRefresh is the timestamp of the last token refresh operation.
	LastRefresh string `json:"last_refresh"`
	// Email is the OpenAI account email address associated with this token.
	Email string `json:"email"`
	// Type indicates the authentication provider type, always "codex" for this storage.
	Type string `json:"type"`
	// Expire is the timestamp when the current access token expires.
	Expire string `json:"expired"`
}

// SaveTokenToFile serializes the Codex token storage to a JSON file.
// This method creates the necessary directory structure and writes the token
// data in JSON format to the specified file path for persistent storage.
//
// Parameters:
//   - authFilePath: The full path where the token file should be saved
//
// Returns:
//   - error: An error if the operation fails, nil otherwise
func (ts *CodexTokenStorage) SaveTokenToFile(authFilePath string) error {
	if ts.AccountID == "" {
		return fmt.Errorf("account_id is required for multi-account storage")
	}
	misc.LogSavingCredentials(authFilePath)
	return accounts.SaveProviderAccount("codex", ts.AccountID, func(existing map[string]any) map[string]any {
		for k, v := range map[string]any{
			"id_token":      ts.IDToken,
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
